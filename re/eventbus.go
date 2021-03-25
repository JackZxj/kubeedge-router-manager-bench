package re

import (
	"log"
	"sync"
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
	var once sync.Once
	var startTime time.Time
	handler := func(client MQTT.Client, msg MQTT.Message) {
		choke <- msg.Payload()
		once.Do(func() {
			startTime = time.Now()
		})
	}

	client := MQTT.NewClient(cli.NewMQTTClientOptions(handler))
	receiveCount := 0
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := client.Subscribe(cli.Topic, cli.QoS, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	timeout := time.Second * time.Duration(cli.TimeoutInSecond)

LOOP:
	for {
		select {
		case <-time.After(timeout):
			endTime := time.Now()
			log.Printf("It timed out %d seconds.\nSum: %d, Faild: %d, Time-Spending: %v",
				cli.TimeoutInSecond, num, num-receiveCount, endTime.Sub(startTime.Add(timeout)))
			break LOOP
		case <-choke:
			receiveCount++
			if receiveCount >= num {
				endTime := time.Now()
				log.Printf("Test done.\nSum: %d, Faild: %d, Time-Spending: %v",
					num, num-receiveCount, endTime.Sub(startTime))
				break LOOP
			}
		}
	}

	client.Disconnect(250)
	return nil
}
