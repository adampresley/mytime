package main

import (
	"github.com/adampresley/mytime/cmd"
)

func main() {
	// var err error

	// cli.VersionPrinter = func(c *cli.Context) {
	// 	fmt.Printf("%s - Time tracking, invoicing, and reporting!\nVersion %s\n\n", Green("My Time"), BrightCyan(Version))
	// }

	// app := &cli.App{
	// 	Name:    "mytime",
	// 	Usage:   "Time tracking, invoicing, and reporting!",
	// 	Version: Version,
	// 	Flags: []cli.Flag{
	// 		&cli.StringFlag{
	// 			Name:  "category",
	// 			Usage: "Specify the category. E.g. mytime start \"projectcode\" --category \"categoryCode\"",
	// 			Value: "",
	// 		},
	// 	},
	// 	Commands: []*cli.Command{
	// 		{
	// 			Name:    "archive",
	// 			Aliases: []string{"a"},
	// 			Usage:   "Archives a client or project",
	// 			Subcommands: []*cli.Command{
	// 				projectActions.Archive(),
	// 			},
	// 		},
	// 		{
	// 			Name:    "create",
	// 			Aliases: []string{"c"},
	// 			Usage:   "Creates a new client, project, or category",
	// 			Subcommands: []*cli.Command{
	// 				clientActions.Create(),
	// 				categoryActions.Create(),
	// 				projectActions.Create(),
	// 			},
	// 		},
	// 		{
	// 			Name:    "list",
	// 			Aliases: []string{"l"},
	// 			Usage:   "Lists clients, projects, or categories",
	// 			Subcommands: []*cli.Command{
	// 				clientActions.List(),
	// 				categoryActions.List(),
	// 				projectActions.List(),
	// 			},
	// 		},
	// 		sessionActions.Start(),
	// 		{
	// 			Name:    "session",
	// 			Aliases: []string{"s", "sess"},
	// 			Usage:   "Actions working with sessions (timing)",
	// 			Subcommands: []*cli.Command{
	// 				sessionActions.Clean(),
	// 				sessionActions.Close(),
	// 				sessionActions.Invoice(),
	// 			},
	// 		},
	// 	},
	// }

	// app.EnableBashCompletion = true

	// if err := app.Run(os.Args); err != nil {
	// 	log.Fatal(err)
	// }

	cmd.Execute()
}
