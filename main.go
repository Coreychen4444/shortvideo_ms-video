package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Coreychen4444/shortvideo"
	"github.com/Coreychen4444/shortvideo_ms-video/repo"
	"github.com/Coreychen4444/shortvideo_ms-video/service"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	mysql_db := repo.InitMysql()
	s := grpc.NewServer()
	pb.RegisterVideoServiceServer(s, service.NewVideoServer(repo.NewDbRepository(mysql_db)))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
