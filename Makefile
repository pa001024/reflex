SHELL = bash
export GOPATH=`pwd`

all: clean windows windows32 linux linux32 darwin32 darwin
windows:
	export GOOS=windows &&\
	export GOARCH=amd64 &&\
	export GOCHAR=6 &&\
	export GOEXE=.exe &&\
	go build -o bin/windows-amd64/t.exe
windows32:
	export GOOS=windows &&\
	export GOARCH=386 &&\
	export GOCHAR=8 &&\
	export GOEXE=.exe &&\
	go build -o bin/windows-386/t.exe
linux:
	export GOOS=linux &&\
	export GOARCH=amd64 &&\
	export GOCHAR=6 &&\
	export GOEXE= &&\
	go build -o bin/linux-amd64/t
linux32:
	export GOOS=linux &&\
	export GOARCH=386 &&\
	export GOCHAR=8 &&\
	export GOEXE= &&\
	go build -o bin/linux-386/t
darwin32:
	export GOOS=darwin &&\
	export GOARCH=386 &&\
	export GOCHAR=8 &&\
	export GOEXE= &&\
	go build -o bin/darwin-386/t
darwin:
	export GOOS=darwin &&\
	export GOARCH=amd64 &&\
	export GOCHAR=6 &&\
	export GOEXE= &&\
	go build -o bin/darwin-amd64/t

test:
	bin/windows-amd64/t logo.png

clean:
	rm -rf t.exe t bin