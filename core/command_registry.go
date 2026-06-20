package core

type ArgumentDefinition struct {
	Name        string
	DisplayText string
	Type        ArgumentType
	Token       string
	Summary     string
	Since       string
	// maybe later
	// deprecated_since, flags for arr of args and value and also key_spec_index
}

type CommandDocs struct {
	Summary    string
	Since      string
	Group      CommandGroup
	Complexity string
	Arguments  []ArgumentDefinition
}

type CommandDefinition struct {
	Name    string
	Arity   int
	Handler func(args []string) (string, error)
	Docs    CommandDocs
}

var Registry = map[string]CommandDefinition{
	"PING": {
		Name:    "PING",
		Arity:   -1,
		Handler: HandlePING,
		Docs: CommandDocs{
			Summary:    "Returns PONG if no argument is provided. Otherwise, returns the provided argument.",
			Since:      "1.0.0",
			Group:      GroupConnection,
			Complexity: "O(1)",
			Arguments: []ArgumentDefinition{
				{
					Name:        "message",
					DisplayText: "message",
					Type:        ArgTypeString,
					Token:       "message",
					Summary:     "The message to be echoed back. If not provided, defaults to 'PONG'.",
					Since:       "1.0.0",
				},
			},
		},
	},
}
