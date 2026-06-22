package main

import (
	"flag"

	"github.com/NILESHD2003/redis-from-scratch/commands"
	"github.com/NILESHD2003/redis-from-scratch/config"
	"github.com/NILESHD2003/redis-from-scratch/server"
)

var conf config.Config

func setup_flags() {
	flag.StringVar(&conf.Host, "host", "0.0.0.0", "Host to listen on")
	flag.IntVar(&conf.Port, "port", 6379, "Port to listen on")
	flag.Parse()
}

func main() {
	setup_flags()
	// Register commands
	commands.RegisterCommands()
	server.StartServer(conf)
}
