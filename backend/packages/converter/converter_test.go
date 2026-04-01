package converter

import "testing"

func TestToInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    int64
		wantErr bool
	}{
		{name: "int", input: int(42), want: 42},
		{name: "int8", input: int8(12), want: 12},
		{name: "int16", input: int16(32000), want: 32000},
		{name: "int32", input: int32(123456), want: 123456},
		{name: "int64", input: int64(922337203685477000), want: 922337203685477000},
		{name: "numeric string", input: "123", want: 123},
		{name: "float string truncates", input: "45.9", want: 45},
		{name: "float64 value truncates", input: 99.8, want: 99},
		{name: "invalid string", input: "abc", wantErr: true},
		{name: "bool invalid", input: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("want=%d, got=%d", tt.want, got)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    float64
		wantErr bool
	}{
		{name: "float32", input: float32(1.5), want: 1.5},
		{name: "float64", input: 2.75, want: 2.75},
		{name: "int value", input: 10, want: 10},
		{name: "numeric string", input: "3.14", want: 3.14},
		{name: "scientific string", input: "1e3", want: 1000},
		{name: "invalid string", input: "xyz", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("want=%v, got=%v", tt.want, got)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    bool
		wantErr bool
	}{
		{name: "bool true", input: true, want: true},
		{name: "bool false", input: false, want: false},
		{name: "string true", input: "true", want: true},
		{name: "string false", input: "false", want: false},
		{name: "string one", input: "1", want: true},
		{name: "string zero", input: "0", want: false},
		{name: "invalid string", input: "not-bool", wantErr: true},
		{name: "numeric invalid", input: 123, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBool(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("want=%v, got=%v", tt.want, got)
			}
		})
	}
}
