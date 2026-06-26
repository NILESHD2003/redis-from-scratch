package core

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrUnknownCommand         = errors.New("ERR unknown command")
	ErrWrongNumberOfArguments = errors.New("Wrong number of arguments for command.")
)

func NewUnknownCommandError(command string) error {
	return fmt.Errorf("%w '%s'", ErrUnknownCommand, strings.ToUpper(command))
}

func Execute(command string, args []string) (RESPValue, error) {
	definition, ok := Registry[strings.ToUpper(command)]

	if !ok {
		return nil, NewUnknownCommandError(command)
	}

	err := validateArity(definition, args)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Executing command: %s with args: %v\n", command, args)

	return definition.Handler(args)
}

func validateArity(definition CommandDefinition, args []string) error {
	// args len + 1 because command itself is also counted as an argument as per redis docs
	arity := len(args) + 1
	// Exact Args matched
	if definition.Arity > 0 {
		if arity != definition.Arity {
			return fmt.Errorf("%w Expected %d arguments but received %d arguments", ErrWrongNumberOfArguments, definition.Arity, arity)
		}
	}

	// Min Args required if -arity
	if definition.Arity < 0 {
		if arity < -definition.Arity {
			return fmt.Errorf("%w Expected at least %d arguments but received %d arguments", ErrWrongNumberOfArguments, -definition.Arity, arity)
		}
	}

	return nil
}
