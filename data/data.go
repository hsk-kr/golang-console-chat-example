package data

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

type Data struct {
	MessageType byte
	Message string
}

const (
	MAX_BUFFER = 1024
)

const (
	MESSAGE_TYPE_NICKNAME = 1
	MESSAGE_TYPE_MESSAGE = 2
)

func GenerateNicknameData(nickname string) []byte {
	return dataToBytes(MESSAGE_TYPE_NICKNAME, nickname)
}

func GenerateMessageData(message string) []byte {
	return dataToBytes(MESSAGE_TYPE_MESSAGE, message)
}

func ConvertBytesToData(b []byte) Data {
	data := Data{}
	decoder := gob.NewDecoder(bytes.NewBuffer(b))
	err := decoder.Decode(&data)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return data
}

func dataToBytes(messageType byte, message string) []byte {
	d := Data{}
	d.MessageType = messageType
	d.Message = message

	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(d)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	return buffer.Bytes()
}