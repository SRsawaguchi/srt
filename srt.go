package srt

import (
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Subtitle represents each subtitle in .srt file.
type Subtitle struct {
	Num     int
	StartMs uint32
	EndMs   uint32
	Text    string
}

// MarkAsDeleted marks its subtitles as deleted.
func (s *Subtitle) MarkAsDeleted() {
	s.Num = -1
}

// IsDeleted can determine whether its subtitles have been delete or not.
func (s *Subtitle) IsDeleted() bool {
	return s.Num == -1
}

func (s *Subtitle) String() string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(s.Num))
	sb.WriteString("\n")
	sb.WriteString(MsToSrtFormat(s.StartMs))
	sb.WriteString(" --> ")
	sb.WriteString(MsToSrtFormat(s.EndMs))
	sb.WriteString("\n")
	sb.WriteString(s.Text)
	sb.WriteString("\n")
	return sb.String()
}

// Srt is struct represents .str file.
type Srt struct {
	Subtitles []*Subtitle
	count     int
}

func (s *Srt) String() string {
	var sb strings.Builder
	s.ForEach(func(_ int, sub *Subtitle) bool {
		sb.WriteString(sub.String())
		return true
	})
	sb.WriteString("\n")
	return sb.String()
}

// GetCount returns count of subtitles that doesn't delete.
func (s *Srt) GetCount() int {
	return s.count
}

// ForEach applies specific function to each subtitle.
// Only not deleted subtitles are passed to the function.
// This function returns how many subtitles not deleted are included.
// Specific function can receive the number of the subtitle that has not been deleted.
// Returning false of the function, ForEach stops calling.
func (s *Srt) ForEach(fn func(int, *Subtitle) bool) int {
	idx := 0
	for _, sub := range s.Subtitles {
		if sub.IsDeleted() {
			continue
		}
		idx++
		if !fn(idx, sub) {
			break
		}
	}
	return idx
}

func (s *Srt) countSubtitle() int {
	var cnt int = 0
	for _, sub := range s.Subtitles {
		if !sub.IsDeleted() {
			cnt++
		}
	}
	return cnt
}

// NewSubtitleFromString makes subtitle object from string property.
func NewSubtitleFromString(strNum, strStart, strEnd, text string) (*Subtitle, error) {
	s := &Subtitle{}
	var err error
	s.Num, err = strconv.Atoi(strNum)
	if err != nil {
		return nil, err
	}

	s.StartMs, err = MillisFromSrtFormat(strStart)
	if err != nil {
		return nil, err
	}

	s.EndMs, err = MillisFromSrtFormat(strEnd)
	if err != nil {
		return nil, err
	}

	s.Text = text
	return s, nil
}

// NewSrtFromString makes Srt object from string.
func NewSrtFromString(str string) (*Srt, error) {
	const (
		_      = iota
		number = iota
		start  = iota
		end    = iota
		text   = iota
	)
	str = str + "\n\n"
	srt := &Srt{
		count: 0,
	}
	re := regexp.MustCompile(`(?m)(-??\d+)\n([0-9,:]+) --> ([0-9,:]+)\n((.+\n)*)^\n`)

	for _, match := range re.FindAllStringSubmatch(str, -1) {
		sub, err := NewSubtitleFromString(match[number], match[start], match[end], match[text])
		if err != nil {
			return nil, err
		}
		srt.Subtitles = append(srt.Subtitles, sub)
		if !sub.IsDeleted() {
			srt.count++
		}
	}
	return srt, nil
}

// NewSrtFromFile make Srt object from .srt file.
func NewSrtFromFile(filePath string) (*Srt, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return NewSrtFromString(string(bytes))
}

// Renumber method restores the number of each subtitle from 1.
// This method returns last number.
func (s *Srt) Renumber() int {
	return s.ForEach(func(i int, sub *Subtitle) bool {
		sub.Num = i
		return true
	})
}

// DeleteSubtitle Delete specified subtitle.
// This method returns the index of the deleted subtitle. If the subtitle has num doesn't exist, it returns -1.
func (s *Srt) DeleteSubtitle(num int) int {
	var idx = -1
	for i, sub := range s.Subtitles {
		if sub.Num == num {
			sub.MarkAsDeleted()
			s.count--
			idx = i
			break
		}
	}
	return idx
}

// TrimTo is delete subtitles until Subtitle.EndMs is smaller than `ms` and other subtitles are subtracted `ms`from that StartMs and EndMs.
// This function returns int value that indicates how many subtitles are deleted.
func (s *Srt) TrimTo(ms uint32) int {
	var deleted int = 0
	s.ForEach(func(_ int, sub *Subtitle) bool {
		if sub.StartMs < ms {
			sub.MarkAsDeleted()
			deleted++
			s.count--
			return true
		}
		sub.StartMs -= ms
		sub.EndMs -= ms
		return true
	})
	return deleted
}

// Cut is delete subtitles between start and end.
// This function returns int value that indicates how many subtitles are deleted.
func (s *Srt) Cut(start, end uint32) int {
	var deleted int = 0
	duration := end - start
	s.ForEach(func(_ int, sub *Subtitle) bool {
		if sub.StartMs >= start && sub.StartMs < end {
			sub.MarkAsDeleted()
			s.count--
			deleted++
		} else if sub.StartMs >= end {
			sub.StartMs -= duration
			sub.EndMs -= duration
		}
		return true
	})
	return deleted
}

// DeleteIf deletes subtitle that judge function return true.
// This returns how many subtitles were deleted.
func (s *Srt) DeleteIf(judge func(int, *Subtitle) bool) int {
	var deleted int = 0
	s.ForEach(func(i int, sub *Subtitle) bool {
		if judge(i, sub) {
			sub.MarkAsDeleted()
			deleted++
			s.count--
		}
		return true
	})
	return deleted
}

// DeleteEmpty deletes subtitles that's text is empty.
// This function returns int value that indicates how many subtitles are deleted.
func (s *Srt) DeleteEmpty() int {
	return s.DeleteIf(func(_ int, sub *Subtitle) bool {
		return strings.TrimSpace(sub.Text) == ""
	})
}

// DeleteByDuration deletes subtitles that's duration less than durationMs.
// This function returns int value that indicates how many subtitles are deleted.
func (s *Srt) DeleteByDuration(durationMs uint32) int {
	return s.DeleteIf(func(_ int, sub *Subtitle) bool {
		return sub.EndMs-sub.StartMs <= durationMs
	})
}

// Sort sort subtitles by start time.
// It uses stable sort algorithm.
func (s *Srt) Sort() {
	sort.SliceStable(s.Subtitles, func(i, j int) bool {
		return s.Subtitles[i].StartMs < s.Subtitles[j].StartMs
	})
}
