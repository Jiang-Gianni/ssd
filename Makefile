fw:
	go run cmd/fw/fw.go

build:
	go build -ldflags="-w -s"