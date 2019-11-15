all:
	make buildserver
	make buildclient
	make dockerimage

dockerimage:
	docker-compose up

buildserver:
	echo "Building server"
	CGO_ENABLED=0 GOOS=linux go build -o ./server/server ./server/server.go

buildclient:
	echo "Building client"
	CGO_ENABLED=0 GOOS=linux go build -o ./client/client ./client/client.go