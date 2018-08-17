package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript/logger/standard"
	"github.com/gen0cide/h3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	l             = standard.NewStandardLogger(nil, "horadric3", "cli", false, false)
	displayBefore = true
	debugOutput   = false
)

func main() {
	app := cli.NewApp()

	app.Writer = color.Output
	app.ErrWriter = color.Output

	app.Name = "h3"
	app.Usage = "assess information security skills"
	app.Description = "Opinionated framework to assess and measure information security skills"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "enables verbose debug output",
			Destination: &debugOutput,
		},
	}

	app.Version = h3.Version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Alex Levinson",
			Email: "gen0cide.threats@gmail.com",
		},
	}
	app.Copyright = "(c) 2018 Alex Levinson"
	app.Commands = []cli.Command{
		discovercmd,
	}

	app.Before = func(c *cli.Context) error {
		if debugOutput {
			l.Logger.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	// ignore error so we don't exit non-zero and break gfmrun README example tests
	err := app.Run(os.Args)
	if err != nil {
		l.Fatalf("Error Encountered: %v", err)
	}
}
