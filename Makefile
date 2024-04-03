.PHONY: play linux_release mac_release windows_release release_all clean

CGO_CFLAGS := -DMA_NO_RUNTIME_LINKING
CGO_LDFLAGS := -framework CoreFoundation -framework CoreAudio -framework AudioToolbox

play:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go run main.go

linux_release:
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ottermatic
	zip -j ottermatic_linux_amd64.zip ottermatic

mac_release:
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ottermatic
	zip -j ottermatic_mac_amd64.zip ottermatic

windows_release:
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ottermatic.exe
	zip -j ottermatic_windows_amd64.zip ottermatic.exe

release_all: linux_release mac_release windows_release
	
clean:
	rm -f ottermatic
	rm -f ottermatic.exe
