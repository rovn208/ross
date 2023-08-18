package util

import "testing"

func TestIsSupportedExtensions(t *testing.T) {
	type args struct {
		ex string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"mp4", args{ex: "mp4"}, true},
		{"mov", args{ex: "mov"}, true},
		{"avi", args{ex: "avi"}, false},
		{"mp3", args{ex: "mp3"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSupportedExtensions(tt.args.ex); got != tt.want {
				t.Errorf("IsSupportedExtensions() = %v, want %v", got, tt.want)
			}
		})
	}
}
