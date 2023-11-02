GO ?= go

.PHONY: all
all: ??-*

??-*: %: .bin %/main.go
	$(GO) build -trimpath -o ".bin/$@" "$@/main.go"

.bin:
	mkdir "$@"

.PHONY: clean
clean:
	$(RM) -r .bin
