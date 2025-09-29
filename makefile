BIN      := ./genuuid
REVISION := `git rev-parse --short HEAD`
VERSION  := `git tag -l | sort -rV | head -n1`
FLAG     := -ldflags='-X main.version='$(VERSION)' -X main.revision='$(REVISION)' -s -w -extldflags="-static" -buildid=' -a -tags netgo -installsuffix -trimpath

# GOOSの値で動作を制御できる。
ifeq ($(GOOS),windows)
	BIN  := $(BIN).exe
endif
ifeq ($(GOOS),linux)
	BIN  := $(BIN)
endif

all:
	cat ./makefile
build:
	make clean
	go build -o $(BIN)
release:
	make clean
	go build $(FLAG) -o $(BIN)
	make upx 
	@echo Success!
upx:
	upx --lzma $(BIN)
clean:
	go generate
	goimports -w *.go
	gofmt -w *.go

release-all:
	GOOS=linux   make release
	GOOS=windows make release

