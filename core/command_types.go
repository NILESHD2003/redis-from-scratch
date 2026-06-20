package core

type ArgumentType string

const (
	ArgTypeString    ArgumentType = "string"
	ArgTypeInteger   ArgumentType = "integer"
	ArgTypeKey       ArgumentType = "key"
	ArgTypeDouble    ArgumentType = "double"
	ArgTypePattern   ArgumentType = "pattern"
	ArgTypeUnixTime  ArgumentType = "unix-time"
	ArgTypePureToken ArgumentType = "pure-token"
)

type CommandGroup string

const (
	GroupConnection CommandGroup = "connection"
	GroupString     CommandGroup = "string"
	GroupGeneric    CommandGroup = "generic"
	GroupServer     CommandGroup = "server"
)
