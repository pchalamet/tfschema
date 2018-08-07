package command

import (
	"flag"
	"fmt"
	"strings"

	"tfschema/plugins"

	"github.com/posener/complete"
)

// DataShowCommand is a command which shows a type definition of data source.
type DataShowCommand struct {
	Meta
	format string
}

// Run runs the procedure of this command.
func (c *DataShowCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("data show", flag.ContinueOnError)
	cmdFlags.StringVar(&c.format, "format", "table", "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if len(cmdFlags.Args()) != 1 {
		c.UI.Error("The data show command expects DATA_SOURCE")
		c.UI.Error(c.Help())
		return 1
	}

	dataSource := cmdFlags.Args()[0]
	providerName, err := detectProviderName(dataSource)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	client, err := tfschema.NewClient(providerName)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	defer client.Kill()

	block, err := client.GetDataSourceSchema(dataSource)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	var out string
	switch c.format {
	case "table":
		out, err = block.FormatTable()
	case "json":
		out, err = block.FormatJSON()
	default:
		c.UI.Error(fmt.Sprintf("Unknown output format: %s", c.format))
		c.UI.Error(c.Help())
		return 1
	}

	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(out)

	return 0
}

// AutocompleteArgs returns the argument predictor.
func (c *DataShowCommand) AutocompleteArgs() complete.Predictor {
	return c.completePredictDataSource()
}

// AutocompleteFlags returns a mapping of supported flags and options.
func (c *DataShowCommand) AutocompleteFlags() complete.Flags {
	return nil
}

// Help returns long-form help text.
func (c *DataShowCommand) Help() string {
	helpText := `
Usage: tfschema data show [options] DATA_SOURCE

Options:

  -format=type    Set output format to table or json (default: table)
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns one-line help text.
func (c *DataShowCommand) Synopsis() string {
	return "Show a type definition of data source"
}
