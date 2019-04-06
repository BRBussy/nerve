package main

import (
	"gitlab.com/iotTracker/nerve/log"
	"gitlab.com/iotTracker/nerve/server"
	ServerMessage "gitlab.com/iotTracker/nerve/server/message"
	ServerMessageGPSPositionHandler "gitlab.com/iotTracker/nerve/server/message/handler/gpsPosition"
	ServerMessageHeartbeatHandler "gitlab.com/iotTracker/nerve/server/message/handler/heartbeat"
	ServerMessageLoginHandler "gitlab.com/iotTracker/nerve/server/message/handler/login"
	ServerMessageStatusHandler "gitlab.com/iotTracker/nerve/server/message/handler/status"
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
	Server.RegisterMessageHandler(ServerMessage.GPSPosition, ServerMessageGPSPositionHandler.New())
	Server.RegisterMessageHandler(ServerMessage.GPSPosition2, ServerMessageGPSPositionHandler.New())
	Server.RegisterMessageHandler(ServerMessage.Status, ServerMessageStatusHandler.New())

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
