bench:
	PORT=8080 go test benchmark/http_test.go -bench=. -benchtime=100x -cpu=10
