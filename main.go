package main

import (
	"fmt"
	abciserver "github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"os/signal"
	"syscall"
	"webchatABCI/client"
)

func createABCIServer() error {

	port, hasPort := os.LookupEnv("ABCI_APP_PORT")
	if !hasPort {
		port = "8080"
	}

	addr := fmt.Sprintf("localhost:%s", port)
	transport := "socket"

	fmt.Printf("Start a new server at %v as a %v server\n", addr, transport)
	app := client.NewWebChatApplication()

	server, _ := abciserver.NewServer(addr, transport, app)

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	server.SetLogger(logger.With("module", "abci-server"))

	if err := server.Start(); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			logger.Info(fmt.Sprintf("captured %v, exiting...", sig))
			_ = server.Stop()
			os.Exit(0)
		}
	}()

	select {}
}

func main() {
	err := createABCIServer()
	if err != nil {
		panic(err)
	}
}
