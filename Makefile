all: build

build: linux darwin windows

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/ydcli
	@echo "To run the binary file on linux, use the command: ./bin/linux/ydcli --urls=<list of urls separated by comma>"

darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/ydcli
	@echo "To run the binary file on darwin, use the command: ./bin/darwin/ydcli --urls=<list of urls separated by comma>"

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/windows/ydcli.exe
	@echo "To run the binary file on windows, use the command: .\bin\windows\ydcli.exe --urls=<list of urls separated by comma>"

clean:
	rm -rf bin
