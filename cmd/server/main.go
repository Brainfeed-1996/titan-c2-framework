package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/titan-c2/framework/internal/db"
	"github.com/titan-c2/framework/internal/server"
	"github.com/titan-c2/framework/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	port     = flag.Int("port", 9090, "The server port")
	certFile = flag.String("cert", "certs/server.crt", "TLS cert file")
	keyFile  = flag.String("key", "certs/server.key", "TLS key file")
)

func main() {
	flag.Parse()

	log.Println("[*] Titan C2 Server - v1.0.0-beta")
	log.Println("[*] Initializing Database...")
	
	database, err := db.NewSQLiteDB("titan.db")
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	defer database.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if *certFile != "" && *keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Printf("[!] Warning: Failed to load TLS keys, running insecure: %v", err)
		} else {
			opts = append(opts, grpc.Creds(creds))
		}
	}

	grpcServer := grpc.NewServer(opts...)
	titanService := server.NewTitanServer(database)
	rpc.RegisterTitanC2Server(grpcServer, titanService)

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("[*] Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("[*] Listening on :%d", *port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
