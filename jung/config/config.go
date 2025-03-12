package config

type Config struct {
	MqttUserName string `env:"MQTT_USERNAME"`
	MqttPassword string `env:"MQTT_PASSWORD"`
}
