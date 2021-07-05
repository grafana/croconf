// TODO: move parts of to this to their own CLI framework library

package main

import (
	"fmt"
	"strings"

	"go.k6.io/croconf"
)

type SubCommand struct {
	Command          string
	AddConfigOptions func() error
	Run              func() error
	// TODO: aliases, hidden commands, error handlers and callbacks, sub-sub
	// commands, etc. - add most of the things and features that CLI frameworks
	// like cobra, kingpin and kong have...
}

// TODO: polish, remove nolint
//nolint:forbidigo
func GetSubcommandHandler(
	cm *croconf.Manager, cliSource *croconf.SourceCLI,
	subCommands []SubCommand, scBinders []croconf.StringValueBinder,
) (func() error, error) {
	subCommandIDs := make([]string, 0, len(subCommands))
	subCommandsByID := make(map[string]SubCommand, len(subCommands))

	for _, sc := range subCommands {
		if _, ok := subCommandsByID[sc.Command]; ok {
			return nil, fmt.Errorf("subcommand %s has more than one handler", sc.Command)
		}
		subCommandIDs = append(subCommandIDs, sc.Command)
		subCommandsByID[sc.Command] = sc
	}
	possibleValues := fmt.Sprintf("possible values: %v", subCommandIDs)

	var showHelp bool
	var subCommand string

	cm.AddField(
		croconf.NewBoolField(
			&showHelp,
			cliSource.FromNameAndShorthand("help", "h"),
		),
		croconf.WithDescription("show help information"),
	)

	cm.AddField(
		croconf.NewStringField(&subCommand, scBinders...),
		croconf.WithDescription(fmt.Sprintf("sub-command (%s)", possibleValues)),
	)

	return func() error {
		if subCommand == "" {
			if showHelp {
				fmt.Println(getHelp(cm))
				return nil
			} else {
				return fmt.Errorf("you have to specify a sub-command (%s), run with --help for help", possibleValues)
			}
		}

		subCmd, ok := subCommandsByID[subCommand]
		if !ok {
			return fmt.Errorf("invalid sub-command '%s', %s", subCommand, possibleValues)
		}

		if err := subCmd.AddConfigOptions(); err != nil {
			return err
		}

		if showHelp {
			fmt.Printf("Help for subcommand %s:\n\n", subCommand)
			fmt.Println(getHelp(cm))
			return nil
		}

		return subCmd.Run()
	}, nil
}

// TODO: make this customizable/templateable
func getHelp(cm *croconf.Manager) string {
	var sb strings.Builder
	for _, field := range cm.Fields() {
		fmt.Fprintf(&sb, "Field '%s' (%s):\n", field.Name, field.Description)
		fmt.Fprintf(&sb, "\tDefault value: %s\n", field.DefaultValue)
		for _, b := range field.Bindings() {
			if fromSource, ok := b.(croconf.BindingFromSource); ok && fromSource.Source() != nil {
				fmt.Fprintf(
					&sb, "\tFrom %s: %s\n",
					fromSource.Source().GetName(), fromSource.BoundName(),
				)
			}
		}
	}

	return sb.String()
}
