package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/nochzato/hush/internal/hushcore"
	"github.com/nochzato/hush/internal/passutils"
	"github.com/urfave/cli/v2"
)

func getMasterPassword() (string, error) {
	fmt.Print("Enter your master password: ")
	masterPassword, err := passutils.ReadPassword(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read master password: %w", err)
	}
	fmt.Println()
	return masterPassword, nil
}

func main() {
	app := &cli.App{
		Name:  "hush",
		Usage: "A CLI tool for password management",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initialize hush and set the master password",
				Action: func(ctx *cli.Context) error {
					masterPassword, err := getMasterPassword()
					if err != nil {
						return err
					}
					return hushcore.InitHush(masterPassword)
				},
			},
			{
				Name:      "add",
				Aliases:   []string{"a"},
				Usage:     "Add a new password entry",
				ArgsUsage: "<name>",
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() < 1 {
						return fmt.Errorf("missing account name")
					}
					name := ctx.Args().First()

					fmt.Print("Enter the password: ")
					password, err := passutils.ReadPassword(os.Stdin)
					if err != nil {
						return fmt.Errorf("failed to read password: %w", err)
					}
					fmt.Println()

					masterPassword, err := getMasterPassword()
					if err != nil {
						return err
					}
					err = hushcore.AddPassword(name, password, masterPassword)
					if err != nil {
						if _, ok := err.(*passutils.PasswordStrengthError); ok {
							fmt.Println("Error: ", err)
						}
						return fmt.Errorf("failed to add password: %w", err)
					}
					fmt.Printf("Password for '%s' added successfully.\n", name)
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List all password names",
				Action: func(ctx *cli.Context) error {
					passwordNames, err := hushcore.ListPasswordNames()
					if err != nil {
						return fmt.Errorf("failed to list passwords: %w", err)
					}

					for _, name := range passwordNames {
						fmt.Println(name)
					}

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

					masterPassword, err := getMasterPassword()
					if err != nil {
						return err
					}

					password, err := hushcore.GetPassword(name, masterPassword)
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
				Name:      "remove",
				Aliases:   []string{"rm"},
				Usage:     "Remove a password",
				ArgsUsage: "<name>",
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() < 1 {
						return fmt.Errorf("missing password name")
					}
					name := ctx.Args().First()

					masterPassword, err := getMasterPassword()
					if err != nil {
						return err
					}

					err = hushcore.RemovePassword(name, masterPassword)
					if err != nil {
						return fmt.Errorf("failed to remove password: %w", err)
					}

					return nil
				},
			},
			{
				Name:    "generate",
				Aliases: []string{"gen"},
				Usage:   "Generate a new password",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "length",
						Aliases: []string{"l"},
						Value:   16,
						Usage:   "Length of the generated password",
					},
					&cli.BoolFlag{
						Name:    "display",
						Aliases: []string{"d"},
						Usage:   "Display password instead of copying to clipboard",
					},
				},
				Action: func(ctx *cli.Context) error {
					length := ctx.Int("length")
					display := ctx.Bool("display")

					password, err := passutils.GeneratePassword(length)
					if err != nil {
						return fmt.Errorf("failed to generate password: %w", err)
					}

					if display {
						fmt.Printf("Generated password: %s\n", password)
					} else {
						if err = clipboard.WriteAll(password); err != nil {
							return fmt.Errorf("failed to copy password to clipboard: %w", err)
						}
						fmt.Println("Password copied to clipboard.")
					}

					fmt.Print("Do you want to save this password? (y/N): ")
					var response string
					fmt.Scanln(&response)

					if response == "y" || response == "Y" {
						fmt.Print("Enter a name for this password: ")
						var name string
						fmt.Scanln(&name)

						masterPassword, err := getMasterPassword()
						if err != nil {
							return err
						}

						if err = hushcore.AddPassword(name, password, masterPassword); err != nil {
							fmt.Errorf("failed to save password: %w", err)
						}
						fmt.Printf("Password saved as %q.\n", name)
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

					masterPassword, err := getMasterPassword()
					if err != nil {
						return err
					}
					err = hushcore.ImplodeHush(masterPassword)
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
