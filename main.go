package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Errorf("Could not load env: %s", err.Error())
		return
	}

	err = env.Parse(&config.Config)
	if err != nil {
		log.Errorf("Could not load serviceConfig: %s", err.Error())
	}

	mqtt.ServeMqtt(srv)
	serveTLS(apiHandler, dependencies.Configs["http_gateway-API"].Port)

}

func serveTLS(handler http.Handler, addr int) {
	tlsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", addr),
		Handler: handler,
	}
	lis, err := net.Listen("tcp", tlsServer.Addr)
	if err != nil {
		return
	}

	teardown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		shutdownErr := tlsServer.Shutdown(ctx)
		if shutdownErr != nil {
			log.Fatal("TLS-server: Shutdown Error!")
		}
	}

	log.Infof("Listening on :%d", addr)
	serveErr := tlsServer.ServeTLS(lis, "volumes-data/certs/http_gateway-server-cert.pem", "secret/certs/http_gateway-server-key.pem")
	defer func() {
		teardown()
	}()
	if serveErr != nil {
		log.Fatalf("HTTP-Gateway-Server: Serving Error: %s", serveErr.Error())
	}
}

func rHandler(rs *server.HttpGatewayServer, srv *service.Service) http.Handler {
	r := regexphandler.RegexpHandler{}

	r.HandleFunc(regexp.MustCompile(fmt.Sprintf("^%s$", service.LoginRoute)), mux.Handler(service.LoginRoute, rs, srv))
	r.HandleFunc(regexp.MustCompile(fmt.Sprintf("^%s.*$", service.JungRoute)), mux.Handler(service.JungRoute, rs, srv))
	r.HandleFunc(regexp.MustCompile(fmt.Sprintf("^%s$", config.ApiPrefix+config.RunPrefix+"/health")), mux.Handler("/health", rs, srv))
	r.HandleFunc(regexp.MustCompile(fmt.Sprintf("^%s$", config.ApiPrefix+config.RunPrefix+"/athlete")), mux.Handler("/athlete", rs, srv))
	r.HandleFunc(regexp.MustCompile(fmt.Sprintf("^%s$", config.ApiPrefix+config.RunPrefix+"/activities")), mux.Handler("/activities", rs, srv))
	r.HandleFunc(regexp.MustCompile(fmt.Sprintf("^%s$", config.ApiPrefix+config.RunPrefix+"/athlete/create")), mux.Handler("/athlete/create", rs, srv))
	r.HandleFunc(regexp.MustCompile(fmt.Sprintf("^%s$", config.ApiPrefix+config.RunPrefix+"/client/create")), mux.Handler("/client/create", rs, srv))
	r.HandleFunc(regexp.MustCompile(fmt.Sprintf("^%s$", config.ApiPrefix+config.RunPrefix+"/weeksummary")), mux.Handler("/weeksummary", rs, srv))
	r.HandleFunc(regexp.MustCompile(fmt.Sprintf("^%s$", config.ApiPrefix+config.RunPrefix+"/weeksummaries")), mux.Handler("/weeksummaries", rs, srv))
	r.HandleFunc(regexp.MustCompile(fmt.Sprintf("^%s$", config.ApiPrefix+config.RunPrefix+"/activitiesToDB")), mux.Handler("/activitiesToDB", rs, srv))

	rWithAuth := server.AuthMiddleware(r)

	return rWithAuth
}
