package main

import (
	"bufio"
	"flag"
	"net"
	"fmt"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
	if err != nil {
		fmt.Println("Error!", err)
		return
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for{
		conn, err := ln.Accept()
		handleError(err)
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)
	msg, err := reader.ReadString('\n')
	handleError(err)
	fmt.Println("server received message:")
	fmt.Println(msg)
	temp := Message{clientid, msg}      //need to create a new temp type Message sending back to channel
	msgs <- temp
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	clientID := 0
	
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.
	ln, err := net.Listen("tcp", *portPtr)
	handleError(err)
	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	fmt.Println("server is ready ...")
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients channel
			// - start to asynchronously handle messages from this client
			clientID++
			clients[clientID] = conn
			go handleClient(conn,clientID,msgs)

		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			for i := range clients{
				if i != msg.sender{
					fmt.Println("ready to send back to client ", i, "with message: ", msg.message)
					clients[i].Write([]byte(msg.message))
				}
			}

		}
	}
}
