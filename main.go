package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/pitsanujiw/go-boilerplate/cmd"
)

func main() {
	app := &cli.App{
		Name:     "Member Activate",
		Usage:    "microservice for Member activate",
		Version:  "0.0.1-alpha",
		Commands: cmd.Commands(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
