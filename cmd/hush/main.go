package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nochzato/hush/internal/hushcore"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "hush",
		Usage: "A CLI tool for password management",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initialize hush and set the master password",
				Action: func(ctx *cli.Context) error {
					return hushcore.InitHush()
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "Add a new password entry",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "Name of the password entry",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password",
						Aliases:  []string{"p"},
						Usage:    "Password for the entry",
						Required: true,
					}},
				Action: func(ctx *cli.Context) error {
					name := ctx.String("name")
					password := ctx.String("password")
					err := hushcore.AddPassword(name, password)
					if err != nil {
						return fmt.Errorf("failed to add password: %w", err)
					}
					fmt.Printf("Password for '%s' added successfully.\n", name)
					return nil
				},
			},
			{
				Name:      "get",
				Aliases:   []string{"g"},
				Usage:     "Retrieve a password",
				ArgsUsage: "<name>",
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() < 1 {
						return fmt.Errorf("missing password name")
					}
					name := ctx.Args().First()
					password, err := hushcore.GetPassword(name)
					if err != nil {
						return fmt.Errorf("failed to get password: %w", err)
					}
					fmt.Println(password)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
