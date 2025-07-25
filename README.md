docker run -it -v .:/usr/src/myapp -w /usr/src/myapp golang:1.24 bash


GOOS=windows GOARCH=386 go build -o carnet_go.exe
