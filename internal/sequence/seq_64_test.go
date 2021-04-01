package sequence

import (
	"github.com/azzz/snowid/internal/id64"
	"reflect"
	"sync"
	"testing"
)

func Test_seq64MIDtoID64MID(t *testing.T) {
	type args struct {
		machineID  uint64
		sequenceID uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			"return maximum value",
			args{MaxMachineIDValue, MaxSequenceID},
			0b1111111111,
		},

		{
			"return zero value",
			args{},
			0,
		},

		{
			"return value",
			args{0b101, 0b111},
			0b0010100111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seq64MIDtoID64MID(tt.args.machineID, tt.args.sequenceID); got != tt.want {
				t.Errorf("compileMachineID() = %010b, want %010b", got, tt.want)
			}
		})
	}
}

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
			"returns error when sequenceID is too big",
			args{sequenceID: MaxSequenceID + 1},
			nil,
			true,
		},

		{
			"returns Seq64",
			args{machineID: 0b101, sequenceID: 0b111, timer: nil},
			&Seq64{
				machineID:     0b101,
				sequenceID:    0b111,
				number:        0,
				seqTimestamp:  0,
				timer:         nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSeq64(tt.args.machineID, tt.args.sequenceID, tt.args.timer)
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
	t.Run("increments internal counter for the same timestamps", func(t *testing.T) {
		timer := func() uint64 {
			return 42
		}

		seq := Seq64{
			machineID:     0b101,
			sequenceID:    0b111,
			number:        0,
			seqTimestamp:  0,
			mu:            sync.Mutex{},
			timer:         timer,
		}

		got := seq.Next()
		want := id64.New(timer(), 0b0010100111, 0)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Next() got = %v, want %v", got, want)
		}

		got = seq.Next()
		want = id64.New(timer(), 0b0010100111, 1)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Next() got = %v, want %v", got, want)
		}

		got = seq.Next()
		want = id64.New(timer(), 0b0010100111, 2)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Next() got = %v, want %v", got, want)
		}
	})

	t.Run("resets internal counter on a new timestamp", func(t *testing.T) {
		timer := func() uint64 {
			return 42
		}

		seq := Seq64{
			machineID:     0b101,
			sequenceID:    0b111,
			number:        5,
			seqTimestamp:  timer()-1,
			mu:            sync.Mutex{},
			timer:         timer,
		}

		got := seq.Next()
		want := id64.New(timer(), 0b0010100111, 0)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Next() got = %v, want %v", got, want)
		}
	})
}