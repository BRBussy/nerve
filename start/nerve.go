package main

import (
	"gitlab.com/iotTracker/nerve/log"
	nerveServer "gitlab.com/iotTracker/nerve/server"
	zx303Server "gitlab.com/iotTracker/nerve/server/zx303"
	"os"
	"os/signal"
)

func main() {
	// set up zx303 server
	ZX303Server := zx303Server.New()
	go func() {
		err := ZX303Server.Start(&nerveServer.StartRequest{
			Port:      "5021",
			IPAddress: "0.0.0.0",
		})
		log.Fatal("zx303 server stopped: ", err.Error())
	}()

	//Wait for interrupt signal
	systemSignalsChannel := make(chan os.Signal, 1)
	signal.Notify(systemSignalsChannel)
	for {
		select {
		case s := <-systemSignalsChannel:
			log.Info("Application is shutting down.. ( ", s, " )")
			return
		}
	}
}
