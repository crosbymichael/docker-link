
all:
	CGO_ENABLED=0 go build -a engine.go
	docker build -t engine .
