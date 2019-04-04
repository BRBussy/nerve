package server

import (
	"encoding/hex"
	serverException "gitlab.com/iotTracker/nerve/server/exception"
	"gitlab.com/iotTracker/nerve/server/message"
	"io"
	"strings"
)

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, io.EOF
	}

	// convert input data to hex string
	hexDataString := hex.EncodeToString(data)

	// look for start and end of message markers
	startIdx := strings.Index(hexDataString, message.StartMarker)
	endIdx := strings.Index(hexDataString, message.EndMarker)

	// if start could not be found return an error
	if startIdx == -1 {
		return 0, nil, serverException.StartOfMessageNotFound{
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
		messageHexStringBytes = []byte(hexDataString[startIdx : endIdx+len(message.EndMarker)])
	}
	tokenString := hexDataString[startIdx+len(message.StartMarker) : endIdx]

	// convert message hex string bytes back to bytes
	messageBytes := make([]byte, hex.DecodedLen(len(messageHexStringBytes)))
	noMessageBytes, err := hex.Decode(messageBytes, messageHexStringBytes)
	if err != nil {
		return 0, nil, serverException.DecodingError{Reasons: []string{err.Error()}}
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
