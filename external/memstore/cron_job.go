package memstore

import (
	"github.com/robfig/cron/v3"
)

type CronJob interface {
	Start()
	End()
	AddTask(spec string, cmd func())
}

type cronjob struct {
	job *cron.Cron
}

func NewCronJob() CronJob {
	return &cronjob{
		job: cron.New(),
	}
}

func (c *cronjob) Start() {
	c.job.Start()
}

func (c *cronjob) End() {
	c.job.Stop()
}

func (c *cronjob) AddTask(spec string, cmd func()) {
	c.job.AddFunc(spec, cmd)
}
