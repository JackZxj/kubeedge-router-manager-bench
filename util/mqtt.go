package util

import (
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	Url             string
	Topic           string
	QoS             byte
	TimeoutInSecond int
}

func (cli *Client) NewMQTTClientOptions(defaultHandler MQTT.MessageHandler) *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(cli.Url)
	if cli.TimeoutInSecond < 1 {
		cli.TimeoutInSecond = 1
	}
	opts.SetConnectTimeout(time.Second * time.Duration(cli.TimeoutInSecond))

	if defaultHandler != nil {
		opts.SetDefaultPublishHandler(defaultHandler)
	}

	return opts
}
