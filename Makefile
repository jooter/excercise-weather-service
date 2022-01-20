utest:
	go test -count=1 -cover ./...

build:
	go build ./cmd/weather-server/. 