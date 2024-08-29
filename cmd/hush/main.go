package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/nochzato/hush/internal/hushcore"
	"github.com/nochzato/hush/internal/whisper"
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
						if _, ok := err.(*whisper.PasswordStrengthError); ok {
							fmt.Println("Error: ", err)
						}
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
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "display",
						Aliases: []string{"d"},
						Usage:   "Display the password instead of copying to clipboard",
					},
				},
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() < 1 {
						return fmt.Errorf("missing password name")
					}
					name := ctx.Args().First()
					displayPassword := ctx.Bool("display")
					password, err := hushcore.GetPassword(name)
					if err != nil {
						return fmt.Errorf("failed to get password: %w", err)
					}
					if displayPassword {
						fmt.Println(password)
					} else {
						err = clipboard.WriteAll(password)
						if err != nil {
							return fmt.Errorf("failed to copy password to clipboard: %w", err)
						}
						fmt.Println("Password copied to clipboard.")
					}
					return nil
				},
			},
			{
				Name:  "implode",
				Usage: "Delete all data and remove the .hush directory",
				Action: func(ctx *cli.Context) error {
					fmt.Println("WARNING: This will delete all your stored passwords and remove the .hush directory.")
					fmt.Println("This action is non-reversible and all data will be lost.")
					fmt.Print("Are you sure you want to continue? (y/N): ")

					reader := bufio.NewReader(os.Stdin)
					response, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read user input: %w", err)
					}

					response = strings.TrimSpace(strings.ToLower(response))
					if response != "y" {
						fmt.Println("Operation cancelled.")
						return nil
					}

					err = hushcore.ImplodeHush()
					if err != nil {
						return fmt.Errorf("failed to implode hush: %w", err)
					}

					fmt.Println("Hush has been successfully imploded. All data has been deleted.")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
