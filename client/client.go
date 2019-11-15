package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/bprayush/url_shortener/proto"
	"google.golang.org/grpc"
)

var url string
var short string

func main() {

	flag.StringVar(&url, "url", "", "Url to shorten.")
	flag.StringVar(&short, "short", "", "Custom short url to redirect from.")
	flag.Parse()

	if url == "" || short == "" {
		log.Fatalln("Usage: client -url <url> -short <custom short endpoint>")
		return
	}

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Error dialing to server")
		return
	}

	client := pb.NewURLShortenerClient(conn)

	response, err := client.AddURL(context.Background(), &pb.AddURLRequest{
		Url:            url,
		CustomEndpoint: short,
	})
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Println(response)
}
