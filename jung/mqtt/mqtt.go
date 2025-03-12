package mqtt

import (
	"de.janmeckelholt.jung/config"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var Client mqtt.Client

func ServeMqtt(conf *config.Config) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("mqtt:1883")
	//opts.AddBroker("192.168.178.61:1883")
	opts.SetClientID("go_mqtt_client")
	//opts.SetUsername(srv.Config.MqttUserName)
	//opts.SetPassword(srv.Config.MasterPassword)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	Client = mqtt.NewClient(opts)
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		log.Errorf("user: __%s__, password: __%s__", conf.MqttUserName, conf.MqttPassword)
		panic(token.Error())
	}

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
	token := client.Publish(topic, 0, false, message)
	token.Wait()
	time.Sleep(time.Second)
}
