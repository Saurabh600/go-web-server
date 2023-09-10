all: build run

build: cmd/main.go
	go build -o bin/app cmd/main.go

run: ./bin/app
	./bin/app

clean: ./bin/app
	@rm -fv ./bin/app
