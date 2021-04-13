bench:
	PORT=8080 go test benchmark/http_test.go -bench=. -benchtime=100x -cpu=10

run:
	LOG_LEVEL=info LISTEN=":8080" EPOCH="20210413001805" MACHINE_ID=333 go run cmd/ratatoskr/*.go