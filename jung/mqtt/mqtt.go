package mqtt

import (
	"fmt"
	"time"

	"de.janmeckelholt.jung/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var Client mqtt.Client

func ServeMqtt(conf *config.Config) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("mosquitto:1883")
	//opts.AddBroker("192.168.178.61:1883")
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(conf.MqttUserName)
	opts.SetPassword(conf.MqttPassword)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	Client = mqtt.NewClient(opts)
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		log.Errorf("user: __%s__, password: __%s__", conf.MqttUserName, conf.MqttPassword)
		panic(token.Error())
	}
	log.Infof("Serving MQTT on mosquitto:1883")

}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func Publish(client mqtt.Client, topic string, message string) {
	if !client.IsConnected() {
		log.Errorf("MQTT client not connected, cannot publish to topic %s", topic)
		return
	}
	log.Debugf("Publishing to topic %s with message: %s", topic, message)
	token := client.Publish(topic, 0, false, message)
	token.Wait()
	log.Debugf("Publish token created, waiting for completion...")
	if !token.Wait() {
		log.Errorf("Publish timeout for topic %s", topic)
		return
	}
	log.Debugf("Publish token wait completed")
	if token.Error() != nil {
		log.Errorf("Failed to publish to topic %s: %v", topic, token.Error())
		return
	}
	log.Infof("Successfully published to topic %s: %s", topic, message)
	time.Sleep(time.Second)
}
