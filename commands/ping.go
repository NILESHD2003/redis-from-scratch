package commands

import (
	"github.com/NILESHD2003/redis-from-scratch/core"
)

func HandlePING(args []string) (core.RESPValue, error) {
	if len(args) == 0 {
		return core.SimpleString("PONG"), nil
	}
	return core.BulkString(args[0]), nil
}

var PingCommand = core.CommandDefinition{
	Name:    "PING",
	Arity:   -1,
	Handler: HandlePING,
	Docs: core.CommandDocs{
		Summary:    "Returns PONG if no argument is provided. Otherwise, returns the provided argument.",
		Since:      "1.0.0",
		Group:      core.GroupConnection,
		Complexity: "O(1)",
		Arguments: []core.ArgumentDefinition{
			{
				Name:        "message",
				DisplayText: "message",
				Type:        core.ArgTypeString,
				Token:       "message",
				Summary:     "The message to be echoed back. If not provided, defaults to 'PONG'.",
				Since:       "1.0.0",
			},
		},
	},
}
