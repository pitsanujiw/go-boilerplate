package health

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/pitsanujiw/go-boilerplate/config"
	"github.com/pitsanujiw/go-boilerplate/pkg/log"
	"github.com/pitsanujiw/go-boilerplate/pkg/server"
)

type Suite struct {
	suite.Suite
	logger  *log.Logger
	handler Handler
	server  server.HTTPServer
}

func Test_Suite_HTTPRequestHealth(t *testing.T) {
	s := new(Suite)
	suite.Run(t, s)
}

func (suite *Suite) SetupSuite() {
	logger, err := log.New(&config.App{Log: config.Log{Name: "testing"}})
	if err != nil {
		suite.Fail(err.Error())
		return
	}

	serv, _ := server.NewTest(logger)
	suite.server = serv
	suite.logger = logger
	suite.handler = New(serv, &config.App{
		Name: "__SERVICE_NAME__",
	})
}

// this function executes after all tests executed
func (suite *Suite) TearDownSuite() {
	// Do something here
}

func (suite *Suite) Test_Health() {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Content-Type", "application/json")
	// Perform the request plain with the app,
	// the second argument is a request latency
	// (set to -1 for no latency)
	resp, err := suite.server.Server().Test(req, -1)
	suite.Nil(err, "error should be nil")

	defer resp.Body.Close()
	//
	respByt, err := io.ReadAll(resp.Body)
	suite.Nil(err, "error should be nil")
	suite.NotNil(respByt, "body should not be nil")
	suite.Equal(http.StatusOK, resp.StatusCode, "http status code not BadRequest")
}
