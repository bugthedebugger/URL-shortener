package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/go-redis/redis"

	pb "github.com/bprayush/url_shortener/proto"
	"google.golang.org/grpc"
)

// Server type struct{}
type Server struct {
	redisClient *redis.Client
}

// AddURL (ctx context.Context, req *pb.AddURLRequest) (*pb.AddURLResponse, error)
func (s *Server) AddURL(ctx context.Context, req *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	return nil, nil
}

// GetURL (ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error)
func (s *Server) GetURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	return nil, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening on port 8080")
		return
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal("Error listening to redis server!")
		return
	}
	fmt.Printf("Redis: %v\n", pong)

	srv := grpc.NewServer()
	pb.RegisterURLShortenerServer(srv, &Server{
		redisClient,
	})
	if err := srv.Serve(listener); err != nil {
		log.Fatal("Error starting serverr on port 8080")
		return
	}

}
