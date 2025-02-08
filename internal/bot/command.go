package bot

import (
	"github.com/bwmarrin/discordgo"
)

// Command represents an executable bot command.
type Command interface {
	Execute(b *BotController, msg *discordgo.MessageCreate, options []string)
	Help() string
}

// CommandRegistry holds available commands.
type CommandRegistry struct {
	commands map[string]Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
	}
}

func (cr *CommandRegistry) Register(name string, cmd Command) {
	cr.commands[name] = cmd
}

func (cr *CommandRegistry) Get(name string) (Command, bool) {
	cmd, ok := cr.commands[name]
	return cmd, ok
}
