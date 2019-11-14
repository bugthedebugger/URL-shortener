package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

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
	val, err := s.rc.Set(req.CustomEndpoint, req.Url, 0).Result()

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
				OriginalURL:  val,
				ShortenedURL: req.URL,
			},
		},
	}, err
}

var wg = sync.WaitGroup{}

func main() {
	wg.Add(2)

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

	registerService := &Server{
		redisClient,
	}

	srv := grpc.NewServer()
	pb.RegisterURLShortenerServer(srv, registerService)

	m := mux.NewRouter()
	m.HandleFunc(
		"/{key}",
		func(rw http.ResponseWriter, r *http.Request) {
			parms := mux.Vars(r)
			redirectURL, err := registerService.GetURL(context.Background(), &pb.GetURLRequest{
				URL: parms["key"],
			})
			if err != nil {
				rw.Write([]byte(fmt.Sprintf("Could not find url %v, Error: %v", parms["key"], err.Error())))
				return
			}
			http.Redirect(rw, r, redirectURL.Url[0].OriginalURL, http.StatusSeeOther)
		},
	)

	go func() {
		log.Println("Starting server ...")
		if err := srv.Serve(listener); err != nil {
			log.Fatal("Error starting serverr on port 8080")
			return
		}
		wg.Done()
	}()

	go func() {
		log.Println("Starting http server")
		http.ListenAndServe(":8000", m)
		wg.Done()
	}()

	wg.Wait()
}
