.PHONY: play

play:
	CGO_CFLAGS="-DMA_NO_RUNTIME_LINKING" CGO_LDFLAGS="-framework CoreFoundation -framework CoreAudio -framework AudioToolbox" go run main.go

