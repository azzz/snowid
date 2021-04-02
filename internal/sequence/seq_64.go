package sequence

import (
	"errors"
	"sync"
)

// MaxMachineIDValue is the biggest value for Machine ID. Machine ID is limited with 5 bits.
// MaxSequenceID is the biggest value for Sequence ID. Sequence ID is limited with 5 bits.
const (
	MaxTSValue        = 0x1FFFFFFFFFF
	MaxMachineIDValue = 0x3FF
	MaxNumberValue    = 0xFFF
)

var (
	MachineIDTooBigErr  = errors.New("machine id is too big")
	TimestampTooBigErr = errors.New("timestamp is too big")
	NumberTooBigErr    = errors.New("number is too big")
)

// Timer returns milliseconds passed from the epoch.
type Timer func() uint64

// Seq64 implements 64 bits long snowflake identifier which includes timestamp, machine ID and an incrementing value.
// Format:
// [ 00000000000000000000000000000000000000000    0     0000000000 000000000000 ]
//   |--------------(1)----------------------| |-(2)-|  |---(3)--| |---(4)----|
// 1) 41 bits for Unix Timestamp with milliseconds
// 2) Reserved bit
// 3) 10 bits for Machine ID
// 4) 12 bits for incrementing number
//
// | field      | max value (dec) | max value (hex) |
// |------------|-----------------|-----------------|
// | Timestamp  | 2199023255551   | 0x1FFFFFFFFFF   |
// | Machine ID | 1023            | 0x3FF           |
// | Number     | 4095            | 0xFFF           |
// The minimum value is 0 for each field.
//
// Seq64 generates incremental IDs in scope of a timestamp.
// It means if we call Seq64.Next() several times within a very short period (1 millisecond) we get incremental ids:
// 0b111...00, 0b111...01, 0b111....10 etc. The counter reset at new millisecond, and starts with 0.
//
// Example:
// seq, err := NewSeq64(15, timer)
// if err != nil {
//   panic(err)
// }
//
// id, err := seq.Next()
// if err != nil {
//   fmt.Println(id)
// }
type Seq64 struct {
	machineID  uint64

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
func NewSeq64(machineID uint64, timer Timer) (*Seq64, error) {
	if machineID > MaxMachineIDValue {
		return nil, MachineIDTooBigErr
	}

	return &Seq64{
		machineID:     machineID,
		timer:         timer,
	}, nil
}

func (s *Seq64) MachineID() uint64 {
	return s.machineID
}

// Next generates a next ID based on the current sequence state.
func (s *Seq64) Next() (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.timer() == s.seqTimestamp {
		s.number++
	} else {
		s.seqTimestamp = s.timer()
		s.number = 0
	}

	if s.seqTimestamp > MaxTSValue {
		return 0, TimestampTooBigErr
	}

	if s.number > MaxNumberValue {
		return 0, NumberTooBigErr
	}

	ts := s.seqTimestamp << 23
	mid := s.machineID << 12
	num := s.number

	result := ts | mid | num

	return result, nil
}
