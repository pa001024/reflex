SHELL = bash

NAME=t

all: clean init windows windows32 linux linux32 darwin32 darwin

init:
	go get
windows:
	export GOOS=windows &&\
	export GOARCH=amd64 &&\
	export GOCHAR=6 &&\
	export GOEXE=.exe &&\
	go build -o bin/windows-amd64/${NAME}.exe
windows32:
	export GOOS=windows &&\
	export GOARCH=386 &&\
	export GOCHAR=8 &&\
	export GOEXE=.exe &&\
	go build -o bin/windows-386/${NAME}.exe
linux:
	export GOOS=linux &&\
	export GOARCH=amd64 &&\
	export GOCHAR=6 &&\
	export GOEXE=.exe &&\
	go build -o bin/linux-amd64/${NAME}
linux32:
	export GOOS=linux &&\
	export GOARCH=386 &&\
	export GOCHAR=8 &&\
	export GOEXE= &&\
	go build -o bin/linux-386/${NAME}
darwin32:
	export GOOS=darwin &&\
	export GOARCH=386 &&\
	export GOCHAR=8 &&\
	export GOEXE= &&\
	go build -o bin/darwin-386/${NAME}
darwin:
	export GOOS=darwin &&\
	export GOARCH=amd64 &&\
	export GOCHAR=6 &&\
	export GOEXE= &&\
	go build -o bin/darwin-amd64/${NAME}

test:
	go build -o ../${NAME}.exe
	./${NAME}

clean:
	rm -rf ${NAME}.exe ${NAME} bin