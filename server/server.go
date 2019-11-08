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
	rc *redis.Client
}

// AddURL (ctx context.Context, req *pb.AddURLRequest) (*pb.AddURLResponse, error)
func (s *Server) AddURL(ctx context.Context, req *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	val, err := s.rc.Set(req.Url, req.CustomEndpoint, 0).Result()

	if err != nil {
		return &pb.AddURLResponse{
			Message: err.Error(),
			Status:  "Error",
		}, err
	}

	return &pb.AddURLResponse{
		Message: "URL shortened",
		Status:  "Ok",
		Url: &pb.ShortenedURL{
			OriginalURL:  req.Url,
			ShortenedURL: val,
		},
	}, err
}

// GetURL (ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error)
func (s *Server) GetURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	val, err := s.rc.Get(req.URL).Result()
	fmt.Println(val)
	return &pb.GetURLResponse{
		Url: []*pb.ShortenedURL{
			&pb.ShortenedURL{
				OriginalURL:  req.URL,
				ShortenedURL: val,
			},
		},
	}, err
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening on port 8080")
		return
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Error listening to redis server! Error: %v", err.Error())
		return
	}
	fmt.Printf("Redis: %v\n", pong)

	srv := grpc.NewServer()
	pb.RegisterURLShortenerServer(srv, &Server{
		redisClient,
	})

	log.Println("Starting server ...")
	if err := srv.Serve(listener); err != nil {
		log.Fatal("Error starting serverr on port 8080")
		return
	}

}
