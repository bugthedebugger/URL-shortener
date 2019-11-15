# URL SHORTENER

## Decription
A golang client CLI gRPC project on URL shortener. 

## TODO:
- Refactor the microservice (specifically the server)
- Make a client GUI for the web with sign in

## Installation
```
make
```

## Usage
### Shorten URL:
```
./client -url <URL to shorten> -shorten <Custom URL Tag>
```
### Visit Shorten URL:
```
host-url:port/custom-tag
```