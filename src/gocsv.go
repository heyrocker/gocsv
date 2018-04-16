package main

import (
	"fmt"
	"os"
	"strings"
)

type Subcommand interface {
	Name() string
	Aliases() []string
	Description() string
	Run([]string)
}

var subcommands []Subcommand

func RegisterSubcommand(sub Subcommand) {
	subcommands = append(subcommands, sub)
}

func init() {
	RegisterSubcommand(&DescribeSubcommand{})
	RegisterSubcommand(&DimensionsSubcommand{})
	RegisterSubcommand(&HeadersSubcommand{})
	RegisterSubcommand(&ViewSubcommand{})
	RegisterSubcommand(&StatsSubcommand{})
	RegisterSubcommand(&RenameSubcommand{})
	RegisterSubcommand(&CleanSubcommand{})
	RegisterSubcommand(&TsvSubcommand{})
	RegisterSubcommand(&DelimiterSubcommand{})
	RegisterSubcommand(&HeadSubcommand{})
	RegisterSubcommand(&TailSubcommand{})
	RegisterSubcommand(&BeheadSubcommand{})
	RegisterSubcommand(&AutoincrementSubcommand{})
	RegisterSubcommand(&StackSubcommand{})
	RegisterSubcommand(&SplitSubcommand{})
	RegisterSubcommand(&FilterSubcommand{})
	RegisterSubcommand(&ReplaceSubcommand{})
	RegisterSubcommand(&SelectSubcommand{})
	RegisterSubcommand(&SortSubcommand{})
	RegisterSubcommand(&SampleSubcommand{})
	RegisterSubcommand(&UniqueSubcommand{})
	RegisterSubcommand(&JoinSubcommand{})
	RegisterSubcommand(&XlsxSubcommand{})
	RegisterSubcommand(&SqlSubcommand{})
}

func usageForSubcommand(subcommand Subcommand) string {
	retval := "  - " + subcommand.Name()
	aliases := subcommand.Aliases()
	if len(aliases) == 1 {
		retval += fmt.Sprintf(" (alias: %s)", aliases[0])
	} else if len(aliases) > 1 {
		retval += fmt.Sprintf(" (aliases: %s)", strings.Join(aliases, ", "))
	}
	retval += fmt.Sprintf("\n      %s\n", subcommand.Description())
	return retval
}

// Keep this in sync with the README.
func usage() string {
	usage := "Usage:\n"
	usage += "  Valid subcommands are:\n"
	for _, subcommand := range subcommands {
		usage += usageForSubcommand(subcommand)
	}
	usage += "See https://github.com/DataFoxCo/gocsv for more documentation."
	return usage
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Fprintln(os.Stderr, "Must provide a valid subcommand.")
		fmt.Fprintf(os.Stderr, "%s\n", usage())
		os.Exit(1)
		return
	}
	subcommandName := args[1]
	if subcommandName == "help" {
		fmt.Fprintf(os.Stderr, "%s\n", usage())
		return
	}
	for _, subcommand := range subcommands {
		if MatchesSubcommand(subcommand, subcommandName) {
			subcommand.Run(args[2:])
			return
		}
	}
	fmt.Fprintf(os.Stderr, "Invalid subcommand \"%s\"\n", subcommandName)
	fmt.Fprintf(os.Stderr, "%s\n", usage())
	os.Exit(1)
}

func MatchesSubcommand(sub Subcommand, name string) bool {
	if name == sub.Name() {
		return true
	}
	for _, alias := range sub.Aliases() {
		if alias == name {
			return true
		}
	}
	return false
}