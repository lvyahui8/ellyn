package main

import (
	"github.com/lvyahui8/ellyn/ellyn_agent"
	"github.com/lvyahui8/ellyn/ellyn_ast"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	ellyn_agent.StartBackend = false
	conf := &ellyn_agent.Configuration{}
	app := &cli.App{
		Name:  "Ellyn",
		Usage: "",
		Commands: []*cli.Command{
			{
				Name:  "update",
				Usage: "update code",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "no-args",
						Destination: &conf.NoArgs,
					},
					&cli.BoolFlag{
						Name:        "no-demo",
						Destination: &conf.NoDemo,
					},
				},
				Action: func(ctx *cli.Context) error {
					dir, err := os.Getwd()
					if err != nil {
						return err
					}
					prog := ellyn_ast.NewProgram2(dir, *conf)
					defer prog.Destroy()
					prog.RollbackAll()
					prog.Visit()
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback code",
				Action: func(ctx *cli.Context) error {
					dir, err := os.Getwd()
					if err != nil {
						return err
					}
					prog := ellyn_ast.NewProgram(dir)
					defer prog.Destroy()
					prog.RollbackAll()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
