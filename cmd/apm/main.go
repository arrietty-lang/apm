package main

import (
	"github.com/arrietty-lang/apm"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "repository url",
				Action: func(cCtx *cli.Context) error {
					return apm.Get(cCtx.Args().First())
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
