package commands

import "github.com/NILESHD2003/redis-from-scratch/core"

func RegisterCommands() {
	core.Register(PingCommand)
}
