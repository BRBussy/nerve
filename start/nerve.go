package main

import (
	"gitlab.com/iotTracker/nerve/log"
	"gitlab.com/iotTracker/nerve/server"
	ServerMessage "gitlab.com/iotTracker/nerve/server/message"
	ServerMessageHeartbeatHandler "gitlab.com/iotTracker/nerve/server/message/handler/heartbeat"
	ServerMessageLoginHandler "gitlab.com/iotTracker/nerve/server/message/handler/login"
	"os"
	"os/signal"
)

func main() {
	// set up  server
	Server := server.New(
		"5021",
		"0.0.0.0",
	)

	Server.RegisterMessageHandler(ServerMessage.Login, ServerMessageLoginHandler.New())
	Server.RegisterMessageHandler(ServerMessage.Heartbeat, ServerMessageHeartbeatHandler.New())

	go func() {
		err := Server.Start()
		log.Fatal(" server stopped: ", err.Error())
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
