package htpplistener

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func Serve(handler http.Handler, addr int) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", addr),
		Handler: handler,
	}
	startServer(server, addr)
}

func ServeTLS(handler http.Handler, addr int) {
	tlsserver := &http.Server{
		Addr:    fmt.Sprintf(":%d", addr),
		Handler: handler,
	}
	startServerTLS(tlsserver, addr, "./certs/jung-server-cert.pem", "./certs/jung-server-key.pem")
}

func startServer(server *http.Server, addr int) {
	lis, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("listen error: %s", err.Error())
	}

	teardown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		shutdownErr := server.Shutdown(ctx)
		if shutdownErr != nil {
			log.Fatalf("server shutdown error: %s", shutdownErr.Error())
		}
	}

	log.Infof("Listening on :%d", addr)
	defer teardown()
	serverErr := server.Serve(lis)
	if serverErr != nil && serverErr != http.ErrServerClosed {
		log.Fatalf("Serving Error: %s", serverErr.Error())
	}
}

func startServerTLS(server *http.Server, addr int, certFile, keyFile string) {
	lis, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("listen error: %s", err.Error())
	}

	teardown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		shutdownErr := server.Shutdown(ctx)
		if shutdownErr != nil {
			log.Fatalf("TLS-server shutdown error: %s", shutdownErr.Error())
		}
	}

	log.Infof("Listening on :%d (TLS)", addr)
	defer teardown()
	serverErr := server.ServeTLS(lis, certFile, keyFile)
	if serverErr != nil && serverErr != http.ErrServerClosed {
		log.Fatalf("Serving Error: %s", serverErr.Error())
	}
}
