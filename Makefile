# unit tests
utest:
	go test -count=1 -cover ./...

# build binary
build:
	go build ./cmd/weather-server/. 

# integration test
itest: build
	./weather-server & 
	sleep 1
	curl http://localhost:8080/v1/weather?city=melbourne
	curl http://localhost:8080/v1/weather?city=melbourne
	curl http://localhost:8080/v1/weather?city=melbourne
	curl http://localhost:8080/v1/weather?city=melbourne
	curl http://localhost:8080/v1/weather?city=melbourne
	pkill -c -f ./weather-server

# coverage report
cover:
	go test -v -count=1 -cover -coverprofile=provider.cov ./pkg/provider/.
	go tool cover -html=provider.cov
