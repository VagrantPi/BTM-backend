
di_gen:
	go install github.com/google/wire/cmd/wire@latest
	wire ./internal/di
test:
	go test -v ./...