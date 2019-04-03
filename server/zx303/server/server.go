package server

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"gitlab.com/iotTracker/nerve/log"
	nerveServer "gitlab.com/iotTracker/nerve/server"
	serverException "gitlab.com/iotTracker/nerve/server/exception"
	zx303ServerMessage "gitlab.com/iotTracker/nerve/server/zx303/message"
	zx303ServerMessageHandler "gitlab.com/iotTracker/nerve/server/zx303/message/handler"
	zx303ServerMessageHeartbeatHandler "gitlab.com/iotTracker/nerve/server/zx303/message/handler/heartbeat"
	zx303ServerMessageLoginHandler "gitlab.com/iotTracker/nerve/server/zx303/message/handler/login"
	zx303ServerMessageProxyHandler "gitlab.com/iotTracker/nerve/server/zx303/message/handler/proxy"
	zx303ServerException "gitlab.com/iotTracker/nerve/server/zx303/server/exception"
	"io"
	"net"
	"strings"
)

type server struct {
	Port           string
	IPAddress      string
	MessageHandler zx303ServerMessageHandler.Handler
}

func New() nerveServer.Server {
	handlerMap := make(map[zx303ServerMessage.Type]zx303ServerMessageHandler.Handler)

	// create and register handlers
	handlerMap[zx303ServerMessage.Login] = zx303ServerMessageLoginHandler.New()
	handlerMap[zx303ServerMessage.Heartbeat] = zx303ServerMessageHeartbeatHandler.New()

	return &server{
		MessageHandler: zx303ServerMessageProxyHandler.New(handlerMap),
	}
}

func (s *server) Start(request *nerveServer.StartRequest) error {
	s.Port = request.Port
	s.IPAddress = request.IPAddress
	log.Info(fmt.Sprintf("Starting ZX303 Server listening at %s:%s", s.IPAddress, s.Port))
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%s", s.IPAddress, s.Port))
	if err != nil {
		return serverException.Listen{Reasons: []string{"zx303", err.Error()}}
	}
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			return serverException.AcceptConnection{Reasons: []string{"zx303", err.Error()}}
		}
		go s.handleConnection(c)
	}
}

func (s *server) handleConnection(c net.Conn) {
	log.Info(fmt.Sprintf("ZX303 serving %s", c.RemoteAddr().String()))
	reader := bufio.NewReaderSize(c, 1024)
	scr := bufio.NewScanner(reader)
	scr.Split(splitFunc)
	for {
		// scan advances the scanner to the next token
		// which in this case is a complete message from the device
		// it returns false when the scan stops by reaching the end
		// of the input or an error
		for scr.Scan() {
			// create message from data token
			inMessage, err := zx303ServerMessage.New(string(scr.Bytes()))
			if err != nil {
				log.Warn(err.Error())
				continue
			}
			// handle the message
			outMessage, err := s.MessageHandler.Handle(inMessage)
			if err != nil {
				log.Warn(err.Error())
				continue
			}
			// if a message needs to be returned, return it
			if outMessage != nil {
				// send the message back
			}
		}
		// check to see if scanner stopped with an error
		if scr.Err() != nil {
			fmt.Println("scanning stopped with an error:", scr.Err().Error())
			break
		}
	}
}

const startMarker = "7878"
const endMarker = "0d0a"

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, io.EOF
	}

	// convert input data to hex string
	hexDataString := hex.EncodeToString(data)

	// look for start and end of message markers
	startIdx := strings.Index(hexDataString, startMarker)
	endIdx := strings.Index(hexDataString, endMarker)

	// if start could not be found return an error
	if startIdx == -1 {
		return 0, nil, zx303ServerException.StartOfMessageNotFound{
			Data:    hexDataString,
			Reasons: []string{"splitting failed"},
		}
	}

	// if end could not be found, but start could
	// signal the Scanner to read more data into the slice and try again
	// with a longer slice starting at the same point in the input
	if endIdx == -1 {
		return 0, nil, nil
	}

	// start and end could both be found
	// set the message and token strings
	var messageHexStringBytes []byte
	if endIdx == len(hexDataString)-1 {
		messageHexStringBytes = []byte(hexDataString[startIdx:])
	} else {
		messageHexStringBytes = []byte(hexDataString[startIdx : endIdx+len(endMarker)])
	}
	tokenString := hexDataString[startIdx+len(startMarker) : endIdx]

	// convert message hex string bytes back to bytes
	messageBytes := make([]byte, hex.DecodedLen(len(messageHexStringBytes)))
	noMessageBytes, err := hex.Decode(messageBytes, messageHexStringBytes)
	if err != nil {
		return 0, nil, zx303ServerException.DecodingError{Reasons: []string{err.Error()}}
	}

	// set the amount that the scanner should advance by
	advance = len(data) - (len(data) - noMessageBytes)

	// set the token
	token = []byte(tokenString)

	// set error to nil
	err = nil

	// return
	return advance, token, err
}
