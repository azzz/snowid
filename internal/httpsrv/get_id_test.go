package httpsrv

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"testing"
)

type DummyRW struct {
	mu         sync.Mutex
	data       []byte
	statusCode int
}

func (d *DummyRW) Header() http.Header {
	return nil
}

func (d *DummyRW) Write(bytes []byte) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.data = bytes
	return len(bytes), nil
}

func (d *DummyRW) WriteHeader(statusCode int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.statusCode = statusCode
}

type StaticSequence struct {
	err error
	num uint64
}

func (s StaticSequence) Next() (uint64, error) {
	return s.num, s.err
}

func TestGetID_ServeHTTP(t *testing.T) {
	t.Run("responses with a generated number", func(t *testing.T) {
		seq := StaticSequence{err: nil, num: 1000}
		handler := GetID{
			logger:   logrus.New(),
			sequence: seq,
		}

		rw := &DummyRW{}
		req, _ := http.NewRequest(http.MethodGet, "", nil)
		handler.ServeHTTP(rw, req)

		if rw.statusCode != http.StatusOK {
			t.Errorf("expected response with status %d, but got %d", http.StatusOK, rw.statusCode)
		}

		decoded := idResponse{}
		if err := json.Unmarshal(rw.data, &decoded); err != nil {
			panic(err)
		}

		if decoded.String != "1000" {
			t.Errorf("expected response includes field String=%s, but got String=%s", "1000", decoded.String)
		}

		if decoded.Numeric != 1000 {
			t.Errorf("expected response includes field Numeric=%d, but got Numeric=%d", 1000, decoded.Numeric)
		}
	})
}
