package srt

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

func toMillis(hour, min, sec, ms int) uint32 {
	var result uint32 = uint32(ms)
	result += uint32(sec) * 1000
	result += uint32(min) * 1000 * 60
	result += uint32(hour) * 1000 * 60 * 60
	return result
}

// MsToSrtFormat converts milliseconds to string for srt format.
func MsToSrtFormat(ms uint32) string {
	hour := ms / (1000 * 60 * 60)
	ms %= 1000 * 60 * 60
	min := ms / (1000 * 60)
	ms %= 1000 * 60
	sec := ms / 1000
	ms %= 1000
	return fmt.Sprintf("%v:%02v:%02v,%03v", hour, min, sec, ms)
}

// MillisFromSrtFormat converts srt formatted time string to milliseconds.
func MillisFromSrtFormat(strTime string) (uint32, error) {
	re := regexp.MustCompile(`(\d+):(\d+):(\d+),(\d+)`)
	const (
		_   = iota
		h   = iota
		min = iota
		sec = iota
		ms  = iota
	)

	match := re.FindStringSubmatch(strTime)
	if len(match) == 0 {
		return 0, errors.New("Wrong format: " + strTime)
	}

	intHour, err := strconv.Atoi(match[h])
	if err != nil {
		return 0, err
	}

	intMin, err := strconv.Atoi(match[min])
	if err != nil {
		return 0, err
	}

	intSec, err := strconv.Atoi(match[sec])
	if err != nil {
		return 0, err
	}

	intMs, err := strconv.Atoi(match[ms])
	if err != nil {
		return 0, err
	}

	return toMillis(intHour, intMin, intSec, intMs), nil
}
