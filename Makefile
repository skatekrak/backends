run:
	go run main.go

install:
	go install .

test:
	go test -v ./...

mod:
	go mod tidy -v
