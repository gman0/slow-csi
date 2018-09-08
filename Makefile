NAME=slowplugin

all: slowplugin

slowplugin:
	if [ ! -d ./vendor ]; then dep ensure; fi
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o  _output/$(NAME) ./cmd/slow

image: slowplugin 
	cp _output/$(NAME) deploy/docker
	docker build -t $(NAME) deploy/docker

clean:
	go clean -r -x
	rm -f deploy/docker/slowplugin
