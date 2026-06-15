.PHONY: seed worker-run server-run gofmt

seed:
	rm -f server/foo.sqlite3
	go -C server run ./cmd/seed/main.go

worker-run:
	go -C server run . -worker

server-run:
	go -C server run . -server

gofmt:
	gofmt -w $$(find server -name '*.go')
