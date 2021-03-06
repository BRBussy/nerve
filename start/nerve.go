package main

import (
	messagingConsumerInstance "gitlab.com/iotTracker/messaging/consumer/instance"
	basicMessagingHub "gitlab.com/iotTracker/messaging/hub/basic"
	messagingMessageHandler "gitlab.com/iotTracker/messaging/message/handler"
	asyncMessagingProducer "gitlab.com/iotTracker/messaging/producer/sync"

	"flag"
	basicJsonRpcClient "gitlab.com/iotTracker/brain/communication/jsonRpc/client/basic"
	authJsonRpcAdaptor "gitlab.com/iotTracker/brain/security/authorization/service/adaptor/jsonRpc"
	zx303DeviceJsonRpcAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/administrator/jsonRpc"
	zx303DeviceJsonRpcAuthenticator "gitlab.com/iotTracker/brain/tracker/zx303/authenticator/jsonRpc"
	zx303TaskJsonRpcAdministrator "gitlab.com/iotTracker/brain/tracker/zx303/task/administrator/jsonRpc"
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

	zx303TaskSubmittedMessageHandler "gitlab.com/iotTracker/nerve/messaging/message/handler/zx303/task/submitted"
	zx303TaskTransitionedMessageHandler "gitlab.com/iotTracker/nerve/messaging/message/handler/zx303/task/transitioned"
)

func main() {
	kafkaBrokers := flag.String("kafkaBrokers", "localhost:9092", "ipAddress:port of each kafka broker node (, separated)")
	brainUrl := flag.String("brainUrl", "http://localhost:9011/api-2", "url of brain service")
	brainAPIUserUsername := flag.String("brainAPIUserUsername", "c877d101-79c1-4226-a309-14823c692d3d", "username of brain api user")
	brainAPIUserPassword := flag.String("brainAPIUserPassword", "12345", "password for brain api user")
	flag.Parse()

	jsonRpcClient := basicJsonRpcClient.New(*brainUrl)
	if err := jsonRpcClient.Login(authJsonRpcAdaptor.LoginRequest{
		UsernameOrEmailAddress: *brainAPIUserUsername,
		Password:               *brainAPIUserPassword,
	}); err != nil {
		log.Fatal("unable to log into brain: " + err.Error())
	}
	log.Info("successfully logged into brain")

	go func() {
		if err := jsonRpcClient.MaintainLogin(); err != nil {
			log.Fatal("error maintaining json rpc client login: ", err.Error())
		}
	}()

	zx303DeviceAdministrator := zx303DeviceJsonRpcAdministrator.New(
		jsonRpcClient,
	)
	zx303DeviceAuthenticator := zx303DeviceJsonRpcAuthenticator.New(
		jsonRpcClient,
	)
	zx303TaskAdministrator := zx303TaskJsonRpcAdministrator.New(
		jsonRpcClient,
	)

	kafkaBrokerNodes := strings.Split(*kafkaBrokers, ",")

	// create a messaging hub
	messagingHub := basicMessagingHub.New()

	// create and start brainQueue producer
	brainQueueProducer := asyncMessagingProducer.New(
		kafkaBrokerNodes,
		"brainQueue",
	)
	if err := brainQueueProducer.Start(); err != nil {
		log.Fatal(err.Error())
	}

	// create and start nerveBroadcast consumer
	nerveBroadcastConsumer := messagingConsumerInstance.New(
		kafkaBrokerNodes,
		"nerveBroadcast",
		[]messagingMessageHandler.Handler{
			zx303TaskSubmittedMessageHandler.New(
				messagingHub,
				zx303TaskAdministrator,
			),
			zx303TaskTransitionedMessageHandler.New(
				messagingHub,
				zx303TaskAdministrator,
			),
		},
	)
	go func() {
		if err := nerveBroadcastConsumer.Start(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	// set up  server
	Server := server.New(
		"7018",
		"0.0.0.0",
		messagingHub,
		zx303DeviceAuthenticator,
		zx303DeviceAdministrator,
	)
	Server.RegisterMessageHandler(ServerMessage.Login, ServerLoginMessageHandler.New(
		zx303DeviceAuthenticator,
	))
	Server.RegisterMessageHandler(ServerMessage.Heartbeat, ServerHeartbeatMessageHandler.New())
	Server.RegisterMessageHandler(ServerMessage.GPSPosition, ServerGPSPositionMessageHandler.New(
		brainQueueProducer,
	))
	Server.RegisterMessageHandler(ServerMessage.GPSPosition2, ServerGPSPositionMessageHandler.New(
		brainQueueProducer,
	))
	Server.RegisterMessageHandler(ServerMessage.Status, ServerStatusMessageHandler.New(
		brainQueueProducer,
	))
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
