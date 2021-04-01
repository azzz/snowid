package sequence

import (
	"errors"
	"github.com/azzz/snowid/internal/id64"
	"sync"
)

// MaxMachineIDValue is the biggest value for Machine ID. Machine ID is limited with 5 bits.
// MaxSequenceID is the biggest value for Sequence ID. Sequence ID is limited with 5 bits.
const (
	MaxMachineIDValue = 0x1F
	MaxSequenceID     = 0x1F
)

var (
	MachineIDTooBigErr  = errors.New("machine id is too big")
	SequenceIDTooBigErr = errors.New("sequence id is too big")
)

// Timer returns milliseconds passed from the epoch.
type Timer func() uint64

// Seq64 implements a sequence of 64-bits long IDs based on timestamp, incremental value,
// and unique combination of machine and sequence ids.
// Seq64 is concurrency safe.
// Seq64 generates incremental IDs in scope of a timestamp.
// It means if we call Seq64.Next() several times within a very short period (1 millisecond) we get incremental ids:
// 0b111...00, 0b111...01, 0b111....10 etc. The counter reset at new millisecond, and starts with 0.
//
// Example:
// seq, _ := NewSeq64(15, 1, timer)
// id := seq.Next()
// numID, _ := id.Uint64()
type Seq64 struct {
	machineID  uint64
	sequenceID uint64

	// number stores incremental counter.
	number        uint64
	// seqTimestamp is the last sequence timestamp.
	seqTimestamp  uint64

	mu    sync.Mutex
	timer Timer
}

// NewSeq64 constructs a new Sequence.
// machineID is the unique machine ID or process ID the application runs on.
// sequenceID is the unique sequence ID to support having multiple sequences in the same application instance.
func NewSeq64(machineID, sequenceID uint64, timer Timer) (*Seq64, error) {
	if machineID > MaxMachineIDValue {
		return nil, MachineIDTooBigErr
	}

	if sequenceID > MaxSequenceID {
		return nil, SequenceIDTooBigErr
	}

	return &Seq64{
		machineID:     machineID,
		sequenceID:    sequenceID,
		timer:         timer,
	}, nil
}

// Next generates a next ID based on the current sequence state.
func (s *Seq64) Next() id64.ID64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	ts := s.timer()
	if ts == s.seqTimestamp {
		s.number++
	} else {
		s.seqTimestamp = ts
		s.number = 0
	}

	return id64.New(ts, seq64MIDtoID64MID(s.machineID, s.sequenceID), s.number)
}

// takes both machine ID and sequence ID and compiles
func seq64MIDtoID64MID(machineID, sequenceID uint64) uint64 {
	mid := machineID << 5
	sid := sequenceID

	return mid | sid
}
