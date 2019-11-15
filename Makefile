all:
	make dockerimage
	make buildserver
	make buildclient

dockerimage:
	docker run --rm -it -d -p 6379:6379 redis

buildserver:
	echo "Building server"
	CGO_ENABLED=0 GOOS=linux go build -o ./server/server ./server/server.go

buildclient:
	echo "Building client"
	CGO_ENABLED=0 GOOS=linux go build -o ./client/client ./client/client.go