biname = teledoor
architectures = amd64 arm

all:
	go fmt *.go
	go test *.go
	GOOS=linux GOARCH=amd64 go build -o $(biname)_amd64 *.go ; \
	GOOS=linux GOARCH=arm go build -o $(biname)_arm *.go ; \
