package id64

import (
	"reflect"
	"testing"
)

func TestID64_Uint64(t *testing.T) {
	type fields struct {
		ts        uint64
		machineID uint64
		number    uint64
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint64
		wantErr bool
	}{
		{
			"return error if timestamp is too big",
			fields{ts: MaxTSValue+1}, 0, true,
		},

		{
			"return error if machine id is too big",
			fields{machineID: MaxMachineIDValue+1}, 0, true,
		},

		{
			"return error if number is too big",
			fields{machineID: MaxNumberValue+1}, 0, true,
		},

		{
			"return 0 for empty id",
			fields{}, 0, false,
		},

		{
			"return the maximum id",
			fields{
				MaxTSValue,
				MaxMachineIDValue,
				MaxNumberValue,
			},
			0b1111111111111111111111111111111111111111101111111111111111111111,
			false,
		},

		{
			"return the id",
			fields{0b111, 0b111, 0b111},
			0b000000000000000000000000000000000000011100000000111000000000111,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := ID64{
				ts:        tt.fields.ts,
				machineID: tt.fields.machineID,
				number:    tt.fields.number,
			}
			got, err := id.Uint64()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Uint64() got = %064b, want %064b", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		ts        uint64
		machineID uint64
		number    uint64
	}
	tests := []struct {
		name string
		args args
		want ID64
	}{
		{
			"build the id",
			args{ts: MaxTSValue, machineID: MaxMachineIDValue, number: MaxNumberValue},
			ID64{
				ts:        MaxTSValue,
				machineID: MaxMachineIDValue,
				number:    MaxNumberValue,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.ts, tt.args.machineID, tt.args.number); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
