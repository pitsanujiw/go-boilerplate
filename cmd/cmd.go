package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/pitsanujiw/go-boilerplate/cmd/apigateway"
	"github.com/pitsanujiw/go-boilerplate/cmd/service"
	"github.com/pitsanujiw/go-boilerplate/cmd/worker"
)

func Commands() []*cli.Command {
	return []*cli.Command{
		service.Command(),
		worker.Command(),
		apigateway.Command(),
	}
}
