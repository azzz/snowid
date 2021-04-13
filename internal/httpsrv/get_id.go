package httpsrv

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/azzz/snowid/internal/sequence"
	"github.com/sirupsen/logrus"
)

// Seq64Generator is any generator of 64-bits long sequences.
type Seq64Generator interface {
	Next() (uint64, error)
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type idResponse struct {
	Numeric uint64 `json:"numeric"`
	String  string `json:"string"`
}

var (
	numberTooBigResponse       = errorResponse{Code: "ID_OVERFLOW", Message: "the incremental number value is overflown, try later"}
	timestampOverflownResponse = errorResponse{Code: "TIMESTAMP_OVERFLOW", Message: "the timestamp value is overflow, the service is unavailablee"}
)

type GetID struct {
	logger   logrus.FieldLogger
	sequence Seq64Generator
}

func NewGetID(logger logrus.FieldLogger, seq Seq64Generator) *GetID {
	return &GetID{
		logger:   logger,
		sequence: seq,
	}
}

func (h *GetID) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id, err := h.sequence.Next()

	if errors.Is(err, sequence.NumberTooBigErr) {
		h.logger.Error(err)
		h.writeJson(w, numberTooBigResponse, http.StatusTooManyRequests)
		return
	} else if errors.Is(err, sequence.TimestampTooBigErr) {
		h.logger.Error(err)
		h.writeJson(w, timestampOverflownResponse, http.StatusServiceUnavailable)
		return
	} else if err != nil {
		h.logger.Error(err)
		h.writeJson(w, errorResponse{Code: "ERROR", Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("binary id: %064b", id)

	response := idResponse{
		Numeric: id,
		String:  strconv.FormatUint(id, 10),
	}

	h.writeJson(w, response, http.StatusOK)
}

func (h *GetID) writeJson(w http.ResponseWriter, v interface{}, status int) {
	data, err := json.Marshal(v)
	if err != nil {
		h.logger.Errorf("failed to marshal json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{}"))
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(data)
}
