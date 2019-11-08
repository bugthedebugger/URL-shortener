package main

import (
	"context"
	"log"

	pb "github.com/bprayush/url_shortener/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Error dialing to server")
		return
	}

	client := pb.NewURLShortenerClient(conn)

	response, err := client.AddURL(context.Background(), &pb.AddURLRequest{
		Url:            "https://www.google.com",
		CustomEndpoint: "google",
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	value, err := client.GetURL(context.Background(), &pb.GetURLRequest{
		URL: "https://www.google.com",
	})

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(response)
	log.Println(value)
}
