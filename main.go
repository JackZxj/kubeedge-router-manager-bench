package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/JackZxj/kubeedge-router-manager-bench/er"
	"github.com/JackZxj/kubeedge-router-manager-bench/re"
)

var banchType string
var rest string
var restMethod string
var mqttTocic string
var mqttUrl string
var msg string
var msgBinary string
var num int
var concurrency int

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	flag.StringVar(&banchType, "type", "re", "Type of banch. Accept: 're' or 'er', means 'rest -> eventbus' or 'eventbus -> rest'.")
	flag.StringVar(&banchType, "t", "re", "Type of banch. Accept: 're' or 'er', means 'rest -> eventbus' or 'eventbus -> rest'.")

	flag.StringVar(&rest, "rest", "", "Restful path to kubeedge cloudcore. Ex: \n\t're' mode: http://127.0.0.1:9443/centos78-edge-0/default/a\n\t'er' mode: 127.0.0.1:8080/test-1 ")
	flag.StringVar(&rest, "r", "", "Restful path to kubeedge cloudcore. Ex: \n\t're' mode: http://127.0.0.1:9443/centos78-edge-0/default/a\n\t'er' mode: 127.0.0.1:8080/test-1 ")

	flag.StringVar(&restMethod, "X", "POST", "Specify request command to use, accept: GET/POST/PUT/DELETE.")

	flag.StringVar(&mqttTocic, "eventbus", "", "MQTT topic to kubeedge eventbus.")
	flag.StringVar(&mqttTocic, "e", "", "MQTT topic to kubeedge eventbus.")

	flag.StringVar(&mqttUrl, "mqtt-broker", "tcp://127.0.0.1:1883", "MQTT broker to edge. Ex: tcp://127.0.0.1:1883")
	flag.StringVar(&mqttUrl, "b", "tcp://127.0.0.1:1883", "MQTT broker to edge. Ex: tcp://127.0.0.1:1883")

	flag.StringVar(&msg, "msg", "", "Message to send.")
	flag.StringVar(&msg, "m", "", "Message to send.")

	flag.StringVar(&msgBinary, "msg-binary", "", "Binary message to send.")
	flag.StringVar(&msgBinary, "mb", "", "Binary message to send.")

	flag.IntVar(&num, "num", 1, "Number of meseages to send.")
	flag.IntVar(&num, "n", 1, "Number of meseages to send.")

	flag.IntVar(&concurrency, "concurrency", 1, "Maximum concurrency for sending.")
	flag.IntVar(&concurrency, "c", 1, "Maximum concurrency for sending.")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if num < 1 {
		num = 1
	}
	if concurrency < 1 {
		concurrency = 1
	}

	switch banchType {
	case "re":
		if rest == "" && mqttTocic == "" {
			return errors.New("rest or eventbus is required")
		}
		if rest != "" {
			if msg == "" && msgBinary == "" {
				return errors.New("msg or msg-binary is required")
			}
			// log.Println(banchType, rest, restMethod, msg, msgBinary, num, concurrency)
			return re.RunRestPub(rest, restMethod, msg, msgBinary, num, concurrency)
		} else {
			// log.Println(banchType, mqttUrl, mqttTocic)
			return re.RunEventbusSub(mqttUrl, mqttTocic, num)
		}
	case "er":
		if rest == "" && mqttTocic == "" {
			return errors.New("rest or eventbus is required")
		}
		if mqttTocic != "" {
			if msg == "" && msgBinary == "" {
				return errors.New("msg or msg-binary is required")
			}
			// log.Println(banchType, mqttUrl, mqttTocic, msg, msgBinary, num, concurrency)
			return er.RunEventbusPub(mqttUrl, mqttTocic, msg, msgBinary, num, concurrency)
		} else {
			// log.Println(banchType, rest)
			return er.RunRestSub(rest, num)
		}
	default:
		return errors.New("unknown type")
	}
}
