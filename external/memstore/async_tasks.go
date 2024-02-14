package memstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"

	"github.com/pitsanujiw/go-boilerplate/config"
	"github.com/pitsanujiw/go-boilerplate/pkg/log"
)

// TaskPublisher is a redis-backed queue for publishing tasks that will be asynchronously
// processed by the TaskServer.
// Wrapper for the package https://github.com/hibiken/asynq
type TaskPublisher interface {
	PublishAsyncTask(task AsyncTask, opt TaskOptions) error

	Close() error
}

type taskPublisher struct {
	client *asynq.Client
}

// TaskServer is a redis-backed queue which processes tasks published by TaskPublisher
// Wrapper for the package https://github.com/hibiken/asynq
type TaskServer interface {
	RegisterTaskWorker(taskType TaskType, tw ProcessTaskFunc)
	Run() error
	StopAndShutdown()
}

type taskServer struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

type TaskOptions struct {
	ProcessAt time.Time
	MaxRetry  int
}

type ProcessTaskFunc func(ctx context.Context, task AsyncTask) error

type errorHandler struct {
	l *log.Logger
}

type RedisConfig struct {
	// Redis server address in "host:port" format.
	Addr string

	// Username to authenticate the current connection when Redis ACLs are used.
	// See: https://redis.io/commands/auth.
	Username string

	// Password to authenticate the current connection.
	// See: https://redis.io/commands/auth.
	Password string

	// Redis DB to select after connecting to a server.
	// See: https://redis.io/commands/select.
	DB int
}

type TaskType string

const (
	TaskTypeConfirm TaskType = "confirmed:task"
)

type AsyncTask struct {
	Type    TaskType
	TaskID  string
	Payload []byte
}

type TaskIdentifier interface {
	// TaskID of the underlying task
	// The format {a}:{b}:{c} is encouraged but not required
	// Read the docs https://github.com/hibiken/asynq
	TaskID() string

	GetTaskType() TaskType

	GetPayload() []byte
}

func NewTaskPublisher(cfg *config.App) TaskPublisher {
	return &taskPublisher{
		client: asynq.NewClient(asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr,
			Username: cfg.Redis.Username,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}),
	}
}

func NewTaskOptions(processAt time.Time, maxRetry int) TaskOptions {
	if processAt.IsZero() {
		processAt = time.Now()
	}

	return TaskOptions{
		ProcessAt: processAt,
		MaxRetry:  maxRetry,
	}
}

func (p *taskPublisher) PublishAsyncTask(task AsyncTask, opt TaskOptions) error {
	payload, err := json.Marshal(task)
	if err != nil {
		return err
	}

	t := asynq.NewTask(
		string(task.Type),
		payload,
		asynq.TaskID(task.TaskID),
		asynq.ProcessAt(opt.ProcessAt),
		asynq.MaxRetry(opt.MaxRetry),
	)

	if _, err := p.client.Enqueue(t); err != nil {
		return err
	}

	return nil
}

func (p *taskPublisher) Close() error {
	return p.client.Close()
}

func NewTaskServer(log *log.Logger, cfg *config.App) TaskServer {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr,
			Username: cfg.Redis.Username,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
		asynq.Config{
			Concurrency: 10,
			RetryDelayFunc: func(n int, _ error, _ *asynq.Task) time.Duration {
				exp := NewExponentialBackOff()
				d := exp.NextBackOff()
				for i := 0; i < n; i++ {
					d = exp.NextBackOff()
				}
				return d
			},
			DelayedTaskCheckInterval: time.Second,
			Logger:                   log.Sugar(),
			ErrorHandler:             &errorHandler{l: log},
		},
	)
	return &taskServer{
		srv: srv,
		mux: asynq.NewServeMux(),
	}
}

func (s *taskServer) RegisterTaskWorker(taskType TaskType, tw ProcessTaskFunc) {
	fn := func(ctx context.Context, task *asynq.Task) error {
		var t AsyncTask
		if err := json.Unmarshal(task.Payload(), &t); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}
		return tw(ctx, t)
	}
	s.mux.HandleFunc(string(taskType), fn)
}

func (s *taskServer) Run() error {
	return s.srv.Run(s.mux)
}

func (s *taskServer) StopAndShutdown() {
	s.srv.Stop()
	s.srv.Shutdown()
}

func (h *errorHandler) HandleError(ctx context.Context, task *asynq.Task, err error) {
	taskID, _ := asynq.GetTaskID(ctx)
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	h.l.Logger.Sugar().Infoln(
		fmt.Sprintf("taskID %v (type %v) encountered err [%v] on retry %v/%v", taskID, task.Type(), err, retried, maxRetry),
	)
}

func NewAsyncTask(t TaskIdentifier) AsyncTask {
	return AsyncTask{
		Type:    t.GetTaskType(),
		TaskID:  t.TaskID(),
		Payload: t.GetPayload(),
	}
}
