package zx303

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"gitlab.com/iotTracker/nerve/log"
	nerveServer "gitlab.com/iotTracker/nerve/server"
	serverException "gitlab.com/iotTracker/nerve/server/exception"
	zx303ServerException "gitlab.com/iotTracker/nerve/server/zx303/exception"
	"net"
	"strings"
)

type server struct {
	Port      string
	IPAddress string
}

func New() nerveServer.Server {
	return &server{}
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

	reader := bufio.NewReaderSize(c, 4096)
	scr := bufio.NewScanner(reader)
	scr.Split(splitFunc)
	for {
		// scan advances the scanner to the next token
		// which in this case is a complete message from the device
		// it returns false when the scan stops by reaching the end
		// of the input or an error
		for scr.Scan() {
			fmt.Println(string(scr.Bytes()))
		}
		// check to see if scanner stopped with an error
		if scr.Err() != nil {
			fmt.Println("scanning stopped with an error:", scr.Err().Error())
			break
		}
	}
}

// data structure: 787807101304011333270D0A
// 7878 01 08 0d0a
const startMarker = "7878"
const endMarker = "0d0a"

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
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
