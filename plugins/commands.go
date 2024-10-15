package plugins

import (
	"errors"
	"fmt"
	"rory-pearson/pkg/log"
	"sync"
)

var (
	ErrCommandNotFound = errors.New("command not found")
)

type Command struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ArgTypes    []string                `json:"arg_types"`
	Function    func(args ...any) error `json:"-"`
}

type CommandsPlugin struct {
	Log      log.Log
	Commands map[string]Command

	mu sync.Mutex
}

func (p *CommandsPlugin) Initialize() (*CommandsPlugin, error) {
	return p, nil
}

func (p *CommandsPlugin) RegisterCommand(command Command) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Commands[command.ID] = command
	p.Log.Info().Str("command_id", command.ID).Str("command_name", command.Name).Msg("registered command")
}

func (p *CommandsPlugin) ExecuteCommand(command string, args ...any) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cmd, ok := p.Commands[command]
	if !ok {
		return ErrCommandNotFound
	}

	// Check if arguments are provided
	if len(cmd.ArgTypes) == 0 && len(args) == 0 {
		// No arguments expected and none provided — valid case
		return cmd.Function(args...)
	}

	// Check if there are missing arguments
	if len(args) < len(cmd.ArgTypes) {
		return fmt.Errorf("missing arguments: expected %d but got %d", len(cmd.ArgTypes), len(args))
	}

	// Validate argument types
	for i, arg := range args {
		expectedType := cmd.ArgTypes[i]
		switch expectedType {
		case "string":
			if _, ok := arg.(string); !ok {
				return fmt.Errorf("argument %d is not of type string", i+1)
			}
		case "int":
			if _, ok := arg.(int); !ok {
				return fmt.Errorf("argument %d is not of type int", i+1)
			}
		// Add other cases as needed
		default:
			return fmt.Errorf("unknown type %s for argument %d", expectedType, i+1)
		}
	}

	// If all checks pass, execute the command function
	return cmd.Function(args...)
}
