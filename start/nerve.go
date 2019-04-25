package main

import (
	asyncMessagingProducer "gitlab.com/iotTracker/messaging/producer/sync"

	"flag"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/authorization/service/adaptor/jsonRpc"
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
	brainUrl := flag.String("brainUrl", "http://localhost:9011/api", "url of brain service")
	brainAPIUserUsername := flag.String("brainAPIUserUsername", "f03c27d6-2eb7-4156-a179-aec187a1baf1", "username of brain api user")
	brainAPIUserPassword := flag.String("brainAPIUserPassword", "lZXkc8YDQymXPpdgURVhcJ2JUsz2/J7aK8h7Cf9N8Gw=", "password for brain api user")
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

	jsonRpcClient := basicJsonRpcClient.New(*brainUrl)
	if err := jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
		UsernameOrEmailAddress: *brainAPIUserUsername,
		Password:               *brainAPIUserPassword,
	}); err != nil {
		log.Fatal("unable to log into brain: " + err.Error())
	}
	log.Info("successfully logged into brain")

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
	signal.Notify(systemSignalsChannel, os.Interrupt)
	for {
		select {
		case s := <-systemSignalsChannel:
			Server.Stop()
			log.Info("Application is shutting down.. ( ", s, " )")
			return
		}
	}
}
