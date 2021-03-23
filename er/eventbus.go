package er

import (
	"log"
	"time"

	"github.com/JackZxj/kubeedge-router-manager-bench/util"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func RunEventbusPub(mqttUrl, mqttTocic, msg, msgBinary string, num, concurrency int) error {
	msgToSend, err := util.GetMsgToSend(msg, msgBinary)
	if err != nil {
		return err
	}
	cli := &util.Client{
		Url:             mqttUrl,
		Topic:           mqttTocic,
		QoS:             1,
		TimeoutInSecond: 300,
	}

	client := MQTT.NewClient(cli.NewMQTTClientOptions(nil))
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	pool := util.NewPool(concurrency)
	resChan := util.NewResult(concurrency)
	startTime := time.Now()
	log.Println("Sending start:", startTime.UTC())

	go resChan.Loop()
	for i := 0; i < num; i++ {
		pool.Add(1)
		go func(){
			defer pool.Done()
			token := client.Publish(cli.Topic, cli.QoS, false, msgToSend)
			token.Wait()
		}()
	}

	pool.Wait()
	endTime := time.Now()
	log.Println("Sending end:", endTime.UTC())
	resChan.DoneChan <- true
	client.Disconnect(250)
	log.Printf("Test done.\nSum: %d, Faild: %d, Time-Spanding: %v",
		num, resChan.FailedCount, endTime.Sub(startTime))
	return nil
}
