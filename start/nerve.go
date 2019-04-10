package main

import (
	asyncMessagingProducer "gitlab.com/iotTracker/brain/messaging/producer/sync"

	"flag"
	"gitlab.com/iotTracker/nerve/log"
	"gitlab.com/iotTracker/nerve/server"
	ServerMessage "gitlab.com/iotTracker/nerve/server/message"
	ServerChargeCompleteMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/chargeComplete"
	ServerChargerConnectedMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/chargerConnected"
	ServerChargerDisconnectedMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/chargerDisconnected"
	ServerDeviceSetupMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/deviceSetup"
	ServerFactorySettingsRestoredMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/factorySettingsRestored"
	ServerGPSPositionMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/gpsPosition"
	ServerHeartbeatMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/heartbeat"
	ServerHibernationMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/hibernation"
	ServerLoginMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/login"
	ServerManualPositionMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/manualPosition"
	ServerOfflineWIFIDataMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/offlineWIFIData"
	ServerStatusMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/status"
	ServerTimeSynchronisationMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/timeSynchronisation"
	ServerWhiteListSynchronisationMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/whiteListSynchronisation"
	ServerWIFIPositionMessageHandler "gitlab.com/iotTracker/nerve/server/message/handler/wifiPosition"
	"os"
	"os/signal"
	"strings"
)

func main() {
	kafkaBrokers := flag.String("kafkaBrokers", "localhost:9092", "ipAddress:port of each kafka broker node (, separated)")
	flag.Parse()

	// set up kafka messaging
	kafkaBrokerNodes := strings.Split(*kafkaBrokers, ",")
	brainQueueProducer := asyncMessagingProducer.New(
		kafkaBrokerNodes,
		"brainQueue",
	)
	log.Info("Starting brainQueue producer")
	if err := brainQueueProducer.Start(); err != nil {
		log.Fatal(err.Error())
	}

	// set up  server
	Server := server.New(
		"7018",
		"0.0.0.0",
	)

	Server.RegisterMessageHandler(ServerMessage.Login, ServerLoginMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.Heartbeat, ServerHeartbeatMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.GPSPosition, ServerGPSPositionMessageHandler.New(
		brainQueueProducer,
	))
	Server.RegisterMessageHandler(ServerMessage.GPSPosition2, ServerGPSPositionMessageHandler.New(
		brainQueueProducer,
	))
	Server.RegisterMessageHandler(ServerMessage.Status, ServerStatusMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.Hibernation, ServerHibernationMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.FactorySettingsRestored, ServerFactorySettingsRestoredMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.OfflineWIFIData, ServerOfflineWIFIDataMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.TimeSynchronisation, ServerTimeSynchronisationMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.DeviceSetup, ServerDeviceSetupMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.WhiteListSynchronisation, ServerWhiteListSynchronisationMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.WIFIPosition, ServerWIFIPositionMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.ChargeComplete, ServerChargeCompleteMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.ChargerConnected, ServerChargerConnectedMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.ChargerDisconnected, ServerChargerDisconnectedMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.ManualPosition, ServerManualPositionMessageHandler.New())

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
