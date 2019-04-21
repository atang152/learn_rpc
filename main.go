package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// A is to host RPC server
// B is to have RPC client talk to A's RPC server
// Resource: https://golang.org/pkg/net/rpc/

// API is a receiver which we will use Register to publishes the receiver's methods in the DefaultServer.
type API struct{}

// Person is a struct that A will use to expose it's RPC method
type Person struct {
	Name string
}

// SayHello is a RPC method
// RPC methods must look schematically like: func (t *T) MethodName(argType T1, replyType *T2) error
func (API) SayHello(person Person, reply *Person) error {
	*reply = person
	return nil
}

func main() {

	// ============================
	// SERVER
	// ============================

	// Allocate memory, returns a pointer to API and register the new object that we can call our method
	api := new(API)

	// Use this function to simulate a connection.
	// Pipe creates a synchronous, in-memory, full duplex network connection.https://golang.org/pkg/net/#Pipe
	connA, connB := net.Pipe()
	defer connA.Close()
	defer connB.Close()

	// Run the RPC Server (serving on connA) in go routine
	go func() {

		// NewServer returns a pointer to a RPC server
		svr := rpc.NewServer()

		// Publishes the receiver's method in the DefaultServer
		svr.Register(api)

		// ServeConn runs the DefaultServer on a single connection.
		// ServeConn blocks, serving the connection until the client hangs up.
		// The caller typically invokes ServeConn in a go statement
		svr.ServeConn(connA)

	}()

	// ============================
	// CLIENT
	// ============================

	var reply Person

	a := Person{"Anto"}

	// Setup a RPC Client (via connB) and call the server.
	client := rpc.NewClient(connB)

	// Call RPC method through server
	err := client.Call("API.SayHello", a, &reply)

	if err != nil {
		log.Fatal("error", err)
	}

	fmt.Println(reply.Name)

}
