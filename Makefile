all:
	go build -o numnum cmd/numnum/numnum.go

test:
	go run cmd/numnum/numnum.go
