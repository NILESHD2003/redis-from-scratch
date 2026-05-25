package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/NILESHD2003/redis-from-scratch/config"
)

func handleClientConnection(c net.Conn, concurrent_clients *int64) {
	defer c.Close()

	// Infinite loop to handle client connection and echo back messages i.e. accept messages from client and send back the same message to client
	for {
		cmd, err := readIncomingCommand(c)

		if err != nil {
			atomic.AddInt64(concurrent_clients, -1)

			if err != io.EOF {
				log.Fatal(
					"[Asynchronous]Error Reading Command from Client:",
					err,
				)
			}

			return
		}

		if err := respondToClient(c, cmd); err != nil {
			log.Fatal("[Asynchronous]Error Responding to Client: ", err)
			atomic.AddInt64(concurrent_clients, -1)
			return
		}
	}
}

func StartServer(config config.Config) {
	log.Println("Starting Redis server on ", config.Host, ":", config.Port)

	var concurrent_clients int64 = 0

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))

	if err != nil {
		log.Fatal("[Asynchronous]Error Starting TCP Server: ", err)
		panic(err)
	}

	// infinite loop to accept incoming connections
	for {
		c, err := lsnr.Accept()

		if err != nil {
			log.Fatal("[Asynchronous]Error Accepting Connection: ", err)
			continue
		}

		atomic.AddInt64(&concurrent_clients, 1)

		// spawing a new goroutine to handle the client connection and echo back messages i.e. accept messages from client and send back the same message to client
		go handleClientConnection(c, &concurrent_clients)
	}
}
