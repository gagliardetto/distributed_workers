package main

import (
	dis "github.com/gagliardetto/authenticated-RPC-in-Go"

	"crypto/tls"
	"fmt"
	"time"
)

var server dis.DistribServer

func init() {
	server = dis.NewServer()

	cert, err := tls.LoadX509KeyPair("keys-for-server/server.public.crt", "keys-for-server/server.private.key")
	if err != nil {
		fmt.Println(err)
		return
	}
	server.Config.Cert = cert

	server.Config.Client.EnableCertVerification = true
	server.Config.Client.ServerName = "distributed.client"
	server.Config.Client.ServerCertPath = "keys-for-server/client.public.crt"
}

func main() {

	server.On(dis.ConnectEvent, func(cc dis.Context) {
		requestPack := dis.Pack{
			Destination: "myNameIs",
			Payload:     "I'd like to know you.",
		}
		responsePack, err := cc.Receive(requestPack)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Awesome to meet you,", responsePack.Payload)
	})

	server.On("ping-pong", func(cc dis.Context) {

		fmt.Printf("\nNew ping-pong strike from client: %q\n", cc.Data)
		time.Sleep(time.Millisecond * 10)

		requestPack := dis.Pack{
			Destination: "ping-pong",
			Payload:     "I'll win",
		}
		err := cc.Trigger(requestPack)
		if err != nil {
			fmt.Println(err)
		}
	})

	server.On("myNameIs", func(cc dis.Context) (string, error) {
		return "Captain America", nil
	})

	go func(server dis.DistribServer) {
		for {
			time.Sleep(time.Second * 2)
			fmt.Println("Connected clients: ", server.CountClients())
		}
	}(server)

	for {
		fmt.Println("Starting server...")
		if err := server.Run("127.0.0.1:3333"); err != nil {
			fmt.Println("Server went down")
			time.Sleep(time.Second * 3)
			continue
		}
	}

}
