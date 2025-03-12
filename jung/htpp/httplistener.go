package htpplistener

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

func ServeTLS(handler http.Handler, addr int) {
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
			log.Fatalf("TLS-server: Shuthdown Error %s", shutdownErr.Error())
		}
	}
	log.Infof("Listening on :%d", addr)
	serverErr := tlsServer.ServeTLS(lis, "./certs/jung-server-cert.pem", "./certs/jung-server-key.pem")
	defer func() {
		teardown()
	}()
	if serverErr != nil {
		log.Fatalf("Serving Error: %s", serverErr.Error())
	}

}
