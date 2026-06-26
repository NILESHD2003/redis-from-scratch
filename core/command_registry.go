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
	Handler func(args []string) (RESPValue, error)
	Docs    CommandDocs
}

var Registry = make(map[string]CommandDefinition)

func Register(command CommandDefinition) {
	Registry[command.Name] = command
}
