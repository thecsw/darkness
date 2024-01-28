package narumi

import (
	"reflect"
	"testing"
	"time"
)

func TestConvertHoloscene(t *testing.T) {
	type args struct {
		HEtime string
	}
	tests := []struct {
		name  string
		args  args
		want  time.Time
		want1 bool
	}{
		{"Test 1", args{"127; 12022 H.E. 0000"}, time.Date(2022, time.January, 127, 0, 0, 0, 0, time.Local), true},
		{"Test 2", args{"127; 12022 H.E."}, time.Date(2022, time.January, 127, 0, 0, 0, 0, time.Local), true},
		{"Test 3", args{"127; 12024 H.E. 1234"}, time.Date(2024, time.January, 127, 12, 34, 0, 0, time.Local), true},
		{"Test 4", args{"127; 12024 H.E. 111"}, time.Date(2024, time.January, 127, 0, 0, 0, 0, time.Local), true},
		{"Test 5", args{""}, time.Time{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ConvertHoloscene(tt.args.HEtime)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertHoloscene() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ConvertHoloscene() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_extractHoloscene(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 string
		want3 string
	}{
		{"Test 1", args{"127; 12022 H.E. 0000"}, "127", "12022", "00", "00"},
		{"Test 2", args{"127; 12022 H.E."}, "127", "12022", "00", "00"},
		{"Test 3", args{"127; 12024 H.E. 1234"}, "127", "12024", "12", "34"},
		{"Test 4", args{"127; 12024 H.E. 111"}, "127", "12024", "00", "00"},
		{"Test 5", args{""}, "", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := extractHoloscene(tt.args.data)
			if got != tt.want {
				t.Errorf("extractHoloscene() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractHoloscene() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("extractHoloscene() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("extractHoloscene() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func Test_getHoloscene(t *testing.T) {
	type args struct {
		dayS    string
		yearS   string
		hourS   string
		minuteS string
	}
	tests := []struct {
		name  string
		args  args
		want  time.Time
		want1 bool
	}{
		{"Test 1", args{"127", "12022", "00", "00"}, time.Date(2022, time.January, 127, 0, 0, 0, 0, time.Local), true},
		{"Test 2", args{"127", "12022", "", ""}, time.Date(2022, time.January, 127, 0, 0, 0, 0, time.Local), true},
		{"Test 3", args{"127", "12024", "12", "34"}, time.Date(2024, time.January, 127, 12, 34, 0, 0, time.Local), true},
		{"Test 4", args{"127", "12024", "00", "00"}, time.Date(2024, time.January, 127, 0, 0, 0, 0, time.Local), true},
		{"Test 5", args{"", "", "", ""}, time.Time{}, false},
		{"Test 6", args{"127", "12022", "", ""}, time.Date(2022, time.January, 127, 0, 0, 0, 0, time.Local), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getHoloscene(tt.args.dayS, tt.args.yearS, tt.args.hourS, tt.args.minuteS)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getHoloscene() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getHoloscene() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
