utest:
	go test -v -count=1 -cover ./...

build:
	go build ./cmd/weather-server/. 

itest:
	./weather-server & 
	sleep 1
	curl http://localhost:8080/v1/weather?city=melbourne
	pkill -c -f ./weather-server


cover:
	go test -v -count=1 -cover -coverprofile=provider.cov ./pkg/provider/.
	go tool cover -html=provider.cov
