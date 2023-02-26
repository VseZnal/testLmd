package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"testLmd/services/product-service/config"
	product_service "testLmd/services/product-service/proto"
	pb "testLmd/services/product-service/proto/product-service"
)

func main() {
	conf := config.ProductConfig()

	productServiceHost := conf.HostProduct
	productServicePort := conf.PortProduct

	productServiceAddress := fmt.Sprintf("%s:%s", productServiceHost, productServicePort)

	err := product_service.Init()
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	lis, err := net.Listen("tcp", productServiceAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	pb.RegisterProductServiceServer(
		server,
		&product_service.Server{},
	)

	log.Printf("server listening at %v", lis.Addr())

	err = server.Serve(lis)

	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
