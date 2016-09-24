package utils_test

import (
	"errors"
	utils "github.com/joaodias/hugito-app/utils"
	"reflect"
	"testing"
)

// Mock a random reader that returns an error
func randomReaderError(n []byte) (int, error) {
	return 0, errors.New("Some error")
}

// Mock a random reader that returns no error
func randomReaderSuccess(n []byte) (int, error) {
	return len(n), nil
}

func TestContainsSubArray(t *testing.T) {
	type args struct {
		sub       []string
		reference []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"one string match", args{[]string{"test"}, []string{"test"}}, true},
		{"one string no match", args{[]string{"test"}, []string{"bla"}}, false},
		{"one string partial match", args{[]string{"test"}, []string{"tes"}}, false},
		{"multiple string match", args{[]string{"one", "two", "three", "four", "five"}, []string{"one", "four"}}, true},
		{"multiple string match", args{[]string{"one", "two", "three", "four", "five"}, []string{"six", "seven"}}, false},
	}
	for _, tt := range tests {
		if got := utils.ContainsSubArray(tt.args.sub, tt.args.reference); got != tt.want {
			t.Errorf("%q. ContainsSubArray() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAreStringsEqual(t *testing.T) {
	type args struct {
		x string
		y string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"equal strings", args{"abc", "abc"}, true},
		{"different strings", args{"abc", "xyz"}, false},
	}
	for _, tt := range tests {
		if got := utils.AreStringsEqual(tt.args.x, tt.args.y); got != tt.want {
			t.Errorf("%q. AreStringsEqual() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	type args struct {
		n      int
		reader utils.RandomReader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"read error", args{16, randomReaderError}, nil, true},
		{"read ok", args{16, randomReaderSuccess}, make([]byte, 16), false},
	}
	for _, tt := range tests {
		got, err := utils.GenerateRandomBytes(tt.args.n, tt.args.reader)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. GenerateRandomBytes() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GenerateRandomBytes() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGenerateRandomString(t *testing.T) {
	type args struct {
		s            int
		randomReader utils.RandomReader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"read error", args{16, randomReaderError}, "", true},
		{"read success", args{16, randomReaderSuccess}, "AAAAAAAAAAAAAAAAAAAAAA==", false},
	}
	for _, tt := range tests {
		got, err := utils.GenerateRandomString(tt.args.s, tt.args.randomReader)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. GenerateRandomString() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. GenerateRandomString() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
