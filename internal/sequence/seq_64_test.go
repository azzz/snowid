package sequence

import (
	"reflect"
	"sync"
	"testing"
)

func TestNewSeq64(t *testing.T) {
	type args struct {
		machineID  uint64
		sequenceID uint64
		timer      Timer
	}
	tests := []struct {
		name    string
		args    args
		want    *Seq64
		wantErr bool
	}{
		{
			"returns error when machineID is too big",
			args{machineID: MaxMachineIDValue + 1},
			nil,
			true,
		},

		{
			"returns Seq64",
			args{machineID: 0b101, sequenceID: 0b111, timer: nil},
			&Seq64{
				machineID:     0b101,
				number:        0,
				seqTimestamp:  0,
				timer:         nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSeq64(tt.args.machineID, tt.args.timer)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSeq64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSeq64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeq64_Next(t *testing.T) {
	overflowTimer := func() uint64 {
		return MaxTSValue + 1
	}

	timerMaximum := func() uint64 {
		return MaxTSValue
	}

	timer42 := func() uint64 {
		return 0b101010
	}

	type fields struct {
		machineID    uint64
		number       uint64
		seqTimestamp uint64
		mu           sync.Mutex
		timer        Timer
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint64
		wantErr bool
	}{
		{
			"return error if timestamp is too big",
			fields{timer: overflowTimer}, 0, true,
		},

		{
			"return error if number is too big",
			fields{number: MaxNumberValue+1, timer: timer42, seqTimestamp: timer42()}, 0, true,
		},

		{
			"increment number for the same timestamp",
			fields{number: 0b10, timer: timer42, seqTimestamp: timer42(), machineID: 0b111},
			0b0000000000000000000000000000000000010101000000000111000000000011,
			false,
		},

		{
			"reset number for new timestamp",
			fields{number: 0b10, timer: timer42, seqTimestamp: timer42()-1, machineID: 0b111},
			0b0000000000000000000000000000000000010101000000000111000000000000,
			false,
		},

		{
			"return the maximum id",
			fields{
				machineID: MaxMachineIDValue,
				number: MaxNumberValue-1,
				timer: timerMaximum,
				seqTimestamp: timerMaximum(),
			},
			0b1111111111111111111111111111111111111111101111111111111111111111,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Seq64{
				machineID:    tt.fields.machineID,
				number:       tt.fields.number,
				seqTimestamp: tt.fields.seqTimestamp,
				mu:           tt.fields.mu,
				timer:        tt.fields.timer,
			}
			got, err := s.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Next() got = %064b, want %064b", got, tt.want)
			}
		})
	}
}