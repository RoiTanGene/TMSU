/*
Copyright 2011-2015 Paul Ruane.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package cli

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"tmsu/common/log"
	"tmsu/common/terminal"
	"tmsu/common/terminal/ansi"
	"tmsu/storage"
)

var HelpCommand = Command{
	Name:        "help",
	Synopsis:    "List subcommands or show help for a particular subcommand",
	Usages:      []string{"tmsu help [OPTION]... [SUBCOMMAND]"},
	Description: `Shows help summary or, where SUBCOMMAND is specified, help for SUBCOMMAND.`,
	Options:     Options{{"--list", "-l", "list commands", false, ""}},
	Exec:        helpExec,
}

var helpCommands map[string]*Command

func helpExec(store *storage.Storage, options Options, args []string) error {
	var colour bool
	if options.HasOption("--color") {
		when := options.Get("--color").Argument
		switch when {
		case "auto":
			colour = terminal.Colour() && terminal.Width() > 0
		case "":
		case "always":
			colour = true
		case "never":
			colour = false
		default:
			return fmt.Errorf("invalid argument '%v' for '--color'", when)
		}
	} else {
		colour = terminal.Colour() && terminal.Width() > 0
	}

	if options.HasOption("--list") {
		listCommands()
	} else {
		switch len(args) {
		case 0:
			summary(colour)
		default:
			commandName := args[0]
			describeCommand(commandName, colour)
		}
	}

	return nil
}

func summary(colour bool) {
	text := "TMSU"
	if colour {
		text = ansi.Bold(text)
	}
	fmt.Println(text)
	fmt.Println()

	var maxWidth int
	commandNames := make([]string, 0, len(helpCommands))
	for _, command := range helpCommands {
		commandName := command.Name
		maxWidth = int(math.Max(float64(maxWidth), float64(len(commandName))))
		commandNames = append(commandNames, commandName)
	}

	sort.Strings(commandNames)

	for _, commandName := range commandNames {
		command, _ := helpCommands[commandName]

		if command.Hidden && log.Verbosity < 2 {
			continue
		}

		synopsis := command.Synopsis
		if !colour {
			synopsis = ansi.Strip(synopsis)
		}

		line := fmt.Sprintf("  %-*v  %v", maxWidth, command.Name, synopsis)

		if command.Hidden && colour {
			line = ansi.Yellow(line)
		}

		terminal.PrintWrapped(line)
	}

	fmt.Println()

	text = "Global options:"
	if colour {
		text = ansi.Bold(text)
	}
	fmt.Println(text)
	fmt.Println()

	printOptions(globalOptions)

	fmt.Println()
	terminal.PrintWrapped("Specify subcommand name for detailed help on a particular subcommand, e.g. tmsu help files")

	fmt.Println()
	terminal.PrintWrapped("To read subcommands from standard input specify - as an argument.")
}

func listCommands() {
	commandNames := make([]string, 0, len(helpCommands))

	for _, command := range helpCommands {
		if command.Hidden && log.Verbosity < 2 {
			continue
		}

		commandNames = append(commandNames, command.Name)
	}

	sort.Strings(commandNames)

	for _, commandName := range commandNames {
		fmt.Println(commandName)
	}
}

func describeCommand(commandName string, colour bool) {
	command := findCommand(helpCommands, commandName)
	if command == nil {
		fmt.Printf("No such command '%v'.\n", commandName)
		return
	}

	// usages
	for _, usage := range command.Usages {
		if colour {
			usage = ansi.Bold(usage)
		}

		terminal.PrintWrapped(usage)
	}

	// description
	fmt.Println()
	description := ansi.ParseMarkup(command.Description)

	if !colour {
		description = ansi.Strip(description)
	}

	terminal.PrintWrapped(description)

	// examples
	if command.Examples != nil && len(command.Examples) > 0 {
		fmt.Println()

		text := "Examples:"
		if colour {
			text = ansi.Bold(text)
		}
		fmt.Println(text)
		fmt.Println()

		for _, example := range command.Examples {
			example = "  " + strings.Replace(example, "\n", "\n  ", -1) // preserve indent
			terminal.PrintWrapped(example)
		}
	}

	// aliases
	if command.Aliases != nil && len(command.Aliases) > 0 {
		fmt.Println()

		if command.Aliases != nil {
			text := "Aliases:"
			if colour {
				text = ansi.Bold(text)
			}
			fmt.Print(text)

			for _, alias := range command.Aliases {
				fmt.Print(" " + alias)
			}
		}

		fmt.Println()
	}

	// options
	if command.Options != nil && len(command.Options) > 0 {
		fmt.Println()

		text := "Options:"
		if colour {
			text = ansi.Bold(text)
		}
		fmt.Println(text)
		fmt.Println()

		printOptions(command.Options)
	}
}

func printOptions(options []Option) {
	maxWidth := 0
	for _, option := range options {
		maxWidth = int(math.Max(float64(maxWidth), float64(len(option.LongName))))
	}

	for _, option := range options {
		line := fmt.Sprintf("  %-2v %-*v  %v", option.ShortName, maxWidth-len(option.LongName), option.LongName, option.Description)
		terminal.PrintWrapped(line)
	}
}
