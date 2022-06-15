run:
	go run main.go

init:
	go install .

test:
	go test -v ./...

mod:
	go mod tidy -v
