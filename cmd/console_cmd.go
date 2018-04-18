package cmd

import (
	"github.com/urfave/cli"

	"github.com/ontio/ontology/cmd/console"
	"github.com/ontio/ontology/cmd/utils"
	"github.com/ontio/ontology/common/log"
)

var (
	ConsoleCommand = cli.Command{
		Action:   utils.MigrateFlags(ontConsole),
		Name:     "console",
		Usage:    "Start an interactive environment",
		Flags:    append(append(NodeFlags, InfoFlags...), ContractFlags...),
		Category: "CONSOLE COMMANDS",
		Description: `
			The Ontology terminal is an interactive shell environment
			which exposes a node admin interface as well as the ontology-go-sdk API.`,
	}
)

func ontConsole(ctx *cli.Context) error {
	cons, err := console.New(ontSdk)
	if err != nil {
		log.Fatalf("Failed to start the Ontology console: %v", err)
	}
	defer cons.Stop(false)

	cons.Welcome()
	cons.Interactive()

	return nil
}
