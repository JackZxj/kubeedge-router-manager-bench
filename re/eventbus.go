package re

import (
	"log"
	"time"

	"github.com/JackZxj/kubeedge-router-manager-bench/util"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func RunEventbusSub(mqttUrl, mqttTocic string, num int) error {
	cli := &util.Client{
		Url:             mqttUrl,
		Topic:           mqttTocic,
		QoS:             1,
		TimeoutInSecond: 300,
	}
	choke := make(chan []byte, 1024)

	handler := func(client MQTT.Client, msg MQTT.Message) {
		choke <- msg.Payload()
	}

	client := MQTT.NewClient(cli.NewMQTTClientOptions(handler))
	receiveCount := 0
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	startTime := time.Now()
	if token := client.Subscribe(cli.Topic, cli.QoS, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	timeout := time.Second * time.Duration(cli.TimeoutInSecond)

LOOP:
	for {
		select {
		case <-time.After(timeout):
			endTime := time.Now()
			log.Printf("Test done. Sum: %d, Faild: %d, Time-Spanding: %v",
				num, num-receiveCount, endTime.Sub(startTime.Add(timeout)))
			break LOOP
		case <-choke:
			receiveCount++
			if receiveCount >= num {
				endTime := time.Now()
				log.Printf("Test done. Sum: %d, Faild: %d, Time-Spanding: %v",
					num, receiveCount, endTime.Sub(startTime))
				break LOOP
			}
		}
	}

	client.Disconnect(250)
	return nil
}
