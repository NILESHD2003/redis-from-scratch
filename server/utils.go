package server

import (
	"net"

	"github.com/NILESHD2003/redis-from-scratch/core"
)

func readIncomingCommand(c net.Conn) ([]string, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])

	if err != nil {
		return nil, err
	}

	tokens, err := core.DecodeRESPString(buf[:n])

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func respondToClient(c net.Conn, cmd string) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}
