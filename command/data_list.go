package command

import (
	"strings"

	"tfschema/plugins"
)

// DataListCommand is a command which lists data sources.
type DataListCommand struct {
	Meta
}

// Run runs the procedure of this command.
func (c *DataListCommand) Run(args []string) int {
	if len(args) != 1 {
		c.UI.Error("The data list command expects PROVIDER")
		c.UI.Error(c.Help())
		return 1
	}

	providerName := args[0]

	client, err := tfschema.NewClient(providerName)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	defer client.Kill()

	res := client.DataSources()

	dataSources := []string{}
	for _, r := range res {
		dataSources = append(dataSources, r.Name)
	}

	c.UI.Output(strings.Join(dataSources, "\n"))

	return 0
}

// Help returns long-form help text.
func (c *DataListCommand) Help() string {
	helpText := `
Usage: tfschema data list PROVIDER
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns one-line help text.
func (c *DataListCommand) Synopsis() string {
	return "List data sources"
}
