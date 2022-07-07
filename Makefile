run:
	go run main.go

init:
	go install .

test:
	go test -v ./...

mod:
	go mod tidy -v

doc:
	swag init --parseDependency --parseInternal

format:
	go fmt && swag fmt --exclude="api/interfaces.go,model,database/pagination_scope.go"