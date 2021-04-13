package benchmark

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func Benchmark_HttpGetID(b *testing.B) {
	port := os.Getenv("PORT")
	if port == "" {
		b.Fatal("PORT env variable is missing")
	}
	url := fmt.Sprintf("http://localhost:%s/id64", port)

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Do(req)

			if err != nil {
				b.Errorf("failed to fetch id: %s", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.Errorf("expected response status 200, but got: %d", resp.StatusCode)
				continue
			}
		}
	})

	b.StopTimer()
}
