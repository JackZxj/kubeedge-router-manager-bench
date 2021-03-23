package util

import "log"

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func GetMsgToSend(msg, msgBinary string) ([]byte, error) {
	if msg != "" && msgBinary != "" {
		log.Println("Warning: msg-binary will be prioritized.")
	}
	var msgToSend []byte
	if msg != "" {
		msgToSend = []byte(msg)
	}
	if msgBinary != "" {
		fileMsg, err := ReadFile(msgBinary)
		if err != nil && msg == "" {
			return nil, err
		}
		if err == nil {
			msgToSend = fileMsg
		}
	}
	return msgToSend, nil
}

