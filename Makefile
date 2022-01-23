
default: utest itest

# unit tests
utest:
	go test -count=1 -cover ./...

# build binary
build:
	go build ./cmd/weather-server/. 

# integration test
# From log message, we should see hitting external provider only twice
itest-in-ci: build
	echo
	./weather-server & 
	echo
	sleep 1 # wait for server startup
	curl -s http://localhost:8080/v1/weather?city=melbourne # hit external
	curl -s http://localhost:8080/v1/weather?city=melbourne # not hit external
	sleep 4 # wait for cache expire
	curl -s http://localhost:8080/v1/weather?city=melbourne # hit external
	curl -s http://localhost:8080/v1/weather?location=melbourne # error
	curl -s http://localhost:8080/v1/weather?city=sydney # error
	curl -s http://localhost:8080/v1/notExistAPI # error

itest: itest-in-ci
	pkill -c -f ./weather-server || true # shutdown server # pkill failed in github actions

# coverage report
cover:
	# go test -v -count=1 -cover -coverprofile=provider.cov ./pkg/provider/.
	go test -v -count=1 -cover -coverprofile=provider.cov ./internal/handler/.
	go tool cover -html=provider.cov
