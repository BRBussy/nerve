package main

import (
	"gitlab.com/iotTracker/nerve/log"
	"gitlab.com/iotTracker/nerve/server"
	ServerMessage "gitlab.com/iotTracker/nerve/server/message"
	ServerDeviceSetupMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/deviceSetup"
	ServerFactorySettingsRestoredMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/factorySettingsRestored"
	ServerGPSPositionMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/gpsPosition"
	ServerHeartbeatMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/heartbeat"
	ServerHibernationMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/hibernation"
	ServerLoginMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/login"
	ServerOfflineWIFIDataMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/offlineWIFIData"
	ServerStatusMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/status"
	ServerTimeSynchronisationMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/timeSynchronisation"
	"os"
	"os/signal"
)

func main() {
	// set up  server
	Server := server.New(
		"5021",
		"0.0.0.0",
	)

	Server.RegisterMessageHandler(ServerMessage.Login, ServerLoginMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.Heartbeat, ServerHeartbeatMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.GPSPosition, ServerGPSPositionMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.GPSPosition2, ServerGPSPositionMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.Status, ServerStatusMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.Hibernation, ServerHibernationMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.FactorySettingsRestored, ServerFactorySettingsRestoredMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.OfflineWIFIData, ServerOfflineWIFIDataMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.TimeSynchronisation, ServerTimeSynchronisationMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.DeviceSetup, ServerDeviceSetupMessageHandler.New())

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
