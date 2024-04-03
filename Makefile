.PHONY: play linux_release mac_release windows_release release_all

play:
	CGO_CFLAGS="-DMA_NO_RUNTIME_LINKING" CGO_LDFLAGS="-framework CoreFoundation -framework CoreAudio -framework AudioToolbox" go run main.go

linux_release:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ottermatic
	zip -j ottermatic_linux_amd64.zip ottermatic
	rm ottermatic

mac_release:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ottermatic
	zip -j ottermatic_mac_amd64.zip ottermatic
	rm ottermatic

windows_release:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ottermatic.exe
	zip -j ottermatic_windows_amd64.zip ottermatic.exe
	rm ottermatic.exe

release_all: linux_release mac_release windows_release
