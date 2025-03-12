package main

import (
	"de.janmeckelholt.jung/config"
	htpplistener "de.janmeckelholt.jung/htpp"
	"de.janmeckelholt.jung/mqtt"
	"fmt"
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Errorf("Could not load env: %s", err.Error())
		return
	}
	conf := config.Config{}
	err = env.Parse(&conf)
	if err != nil {
		log.Errorf("Could not load serviceConfig: %s", err.Error())
	}

	go mqtt.ServeMqtt(&conf)
	htpplistener.ServeTLS(Handler(), 443)

}

func Handler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch {
		case strings.HasPrefix(req.URL.Path, "/jung") && req.Method == http.MethodPost:
			{
				body, err := io.ReadAll(req.Body)
				if err != nil {
					http.Error(rw, "could not read body", http.StatusBadRequest)
					return
				}
				log.Infof("body: %s", body)
				mqtt.Publish(mqtt.Client, "jung", string(body))
				mqtt.Publish(mqtt.Client, "jung", req.URL.Path)

				rw.WriteHeader(http.StatusOK)
				res, err := rw.Write(body)
				if err != nil {
					log.Errorf("error sending Jung response back: %s", err.Error())
				}
				log.Infof("Sending jung response: %d", res)
			}
		default:
			{
				http.Error(rw, fmt.Sprintf("path: %s, method: %s not supported", req.URL.Path, req.Method), http.StatusBadRequest)
				return
			}

		}
	})
}
