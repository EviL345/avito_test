package pvzv1

import (
	"fmt"
	"github.com/EviL345/avito_test/internal/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type PVZServer struct {
	UnimplementedPVZServiceServer
	pvzService handler.PvzService
}

func Start(grpcServerPort string, pvzService handler.PvzService) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcServerPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	RegisterPVZServiceServer(server, &PVZServer{
		pvzService: pvzService,
	})

	reflection.Register(server)

	if err = server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
