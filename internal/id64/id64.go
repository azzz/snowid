package id64

import (
	"errors"
)

// MaxTSValue is the maximum value for Unix Timestamp part of the identifier.
// MaxMachineIDValue is the maximem value for MachineID part of the identifier
// MaxNumberValue is the maximum value for Number part of the identifier
const (
	MaxTSValue        = 0x1FFFFFFFFFF
	MaxMachineIDValue = 0x3FF
	MaxNumberValue    = 0xFFF
)

var (
	TimestampTooBigErr = errors.New("timestamp is too big")
	MachineIDTooBigErr = errors.New("machine id is too big")
	NumberTooBigErr    = errors.New("number is too big")
)

// ID64 implements 64 bits long snowflake identifier which includes timestamp, machine ID and an incrementing value.
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
type ID64 struct {
	ts        uint64
	machineID uint64
	number    uint64
}

// NewID64 creates ID64.
// ts is the identifier timestamp in milliseconds
// machine is the machine id the identifier runs on
// number is an incremental number
func New(ts, machineID, number uint64) ID64 {
	return ID64{
		machineID: machineID,
		number:    number,
		ts:      ts,
	}
}

// Uint64 returns uint representation of the identifier.
func (id ID64) Uint64() (uint64, error) {
	if id.ts > MaxTSValue {
		return 0, TimestampTooBigErr
	}

	if id.machineID > MaxMachineIDValue {
		return 0, MachineIDTooBigErr
	}

	if id.number > MaxNumberValue {
		return 0, NumberTooBigErr
	}

	ts := id.ts << 23
	mid := id.machineID << 12
	num := id.number

	result := ts | mid | num

	return result, nil
}