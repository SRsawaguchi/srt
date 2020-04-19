package srt

import "testing"

func TestToMillis(t *testing.T) {
	expect := uint32(17033733)
	result := toMillis(4, 43, 53, 733)
	if result != expect {
		t.Errorf("result must be %v but %v", expect, result)
	}
}

func TestMsToSrtFormat(t *testing.T) {
	expect := "2:08:10,089"
	ms, _ := MillisFromSrtFormat(expect)
	result := MsToSrtFormat(ms)
	if result != expect {
		t.Errorf("Expected value is %v but %v.", expect, result)
	}
}

func TestMillisFromSrtFormat(t *testing.T) {
	expect := uint32(7690689)
	result, err := MillisFromSrtFormat("02:08:10,689")
	if err != nil {
		t.Error(err)
	}
	if result != expect {
		t.Errorf("result must be %v but %v", expect, result)
	}
}
