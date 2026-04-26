package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"de.janmeckelholt.jung/config"
	htpplistener "de.janmeckelholt.jung/htpp"
	"de.janmeckelholt.jung/mqtt"
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type DatapointAttributes struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type Datapoint struct {
	Attributes DatapointAttributes `json:"attributes"`
}

type JungData struct {
	Data []Datapoint `json:"data"`
}

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
	htpplistener.Serve(Handler(), 80)

}

func Handler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if proto := req.Header.Get("X-Forwarded-Proto"); proto != "" {
			req.URL.Scheme = proto
		}
		if host := req.Header.Get("X-Forwarded-Host"); host != "" {
			req.Host = host
		}

		switch {
		case strings.HasPrefix(req.URL.Path, "/jung") && req.Method == http.MethodPost:
			{
				body, err := io.ReadAll(req.Body)
				if err != nil {
					http.Error(rw, "could not read body", http.StatusBadRequest)
					return
				}
				log.Infof("body: %s", body)

				// Parse the JSON data
				var jungData JungData
				err = json.Unmarshal(body, &jungData)
				if err != nil {
					log.Errorf("Failed to parse JSON: %v", err)
					http.Error(rw, "invalid JSON format", http.StatusBadRequest)
					return
				} 
				if len(jungData.Data) <= 0 {
					log.Errorf("No datapoints found in JSON")
					http.Error(rw, "no datapoints found in JSON", http.StatusBadRequest)
					return
				}
				title := jungData.Data[0].Attributes.Title
				value := jungData.Data[0].Attributes.Value
				log.Infof("Parsed datapoint - Title: %s, Value: %s", title, value)
				

				mqtt.Publish(mqtt.Client, req.URL.Path, value)

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
