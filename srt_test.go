package srt

import (
	"testing"
)

func TestMarkAsDeleted(t *testing.T) {
	sub := Subtitle{
		Num:     1,
		StartMs: 0,
		EndMs:   10,
		Text:    "",
	}
	sub.MarkAsDeleted()
	if sub.Num != -1 {
		t.Errorf("Expected sub.Num is %v, but received %v.", -1, sub.Num)
	}
}

func TestIsDeleted(t *testing.T) {
	sub := Subtitle{
		Num:     -1,
		StartMs: 0,
		EndMs:   10,
		Text:    "",
	}
	if !sub.IsDeleted() {
		t.Errorf("Expected true but the result is false.")
	}
}

func TestSubtitleString(t *testing.T) {
	expect := `5
1:23:02,042 --> 1:23:12,000
This is
text
`
	startMs, _ := MillisFromSrtFormat("1:23:02,042")
	endMs, _ := MillisFromSrtFormat("1:23:12,000")
	sub := Subtitle{
		Num:     5,
		StartMs: startMs,
		EndMs:   endMs,
		Text:    "This is\ntext",
	}
	if sub.String() != expect {
		t.Errorf("Wrong .srt format [%v].", sub.String())
	}
}

func TestNewSubtitleFromString(t *testing.T) {
	sub, err := NewSubtitleFromString("5", "01:23:04,523", "01:23:14,523", "example")
	if err != nil {
		t.Error(err)
	}
	if sub.Num != 5 {
		t.Errorf("subtitle.num is not correct. expected %v but receive %v.", 5, sub.Num)
	}
	if sub.StartMs != 4984523 {
		t.Errorf("subtitle.startMs is not correct. expected %v but receive %v.", 4984523, sub.StartMs)
	}
	if sub.EndMs != 4994523 {
		t.Errorf("subtitle.endMs is not correct. expected %v but receive %v.", 4994523, sub.EndMs)
	}
	if sub.Text != "example" {
		t.Errorf("subtitle.text is not correct. expected %v but receive %v.", "example", sub.Text)
	}
}

func TestNewSrtFromString(t *testing.T) {
	str := `35
0:02:21,860 --> 0:02:24,520
入力の組み合わせは4つだけです。
awef

36
0:02:24,520 --> 0:02:27,400
最初の3つは、
0+0 = 0

36
0:02:24,520 --> 0:02:27,400

37
0:02:24,520 --> 0:02:27,400
最初の3つは、
0+0 = 0`
	srt, err := NewSrtFromString(str)
	if err != nil {
		t.Error(err)
	}
	if len(srt.Subtitles) != 4 {
		t.Errorf("Expected length of subtitle is %v, but %v.", 4, len(srt.Subtitles))
	}
}

func TestSrtString(t *testing.T) {
	str := `1
0:02:21,860 --> 0:02:24,520

3
0:02:24,520 --> 0:02:27,400
最初の3つは、
0+0 = 0

-1
0:02:24,520 --> 0:02:27,400

7
0:02:24,520 --> 0:02:27,400`
	expected := `1
0:02:21,860 --> 0:02:24,520

3
0:02:24,520 --> 0:02:27,400
最初の3つは、
0+0 = 0

7
0:02:24,520 --> 0:02:27,400


`
	srt, _ := NewSrtFromString(str)
	result := srt.String()
	if result != expected {
		t.Errorf("Wrong .srt format. [%v]", srt.String())
	}
}

func TestRenumber(t *testing.T) {
	str := `1
0:02:21,860 --> 0:02:24,520

3
0:02:24,520 --> 0:02:27,400
最初の3つは、
0+0 = 0

-1
0:02:24,520 --> 0:02:27,400

7
0:02:24,520 --> 0:02:27,400`
	srt, err := NewSrtFromString(str)
	if err != nil {
		t.Error(err)
	}
	if srt.GetCount() != 3 {
		t.Errorf("Expected count of subtitle is %v, but %v.", 3, srt.GetCount())
	}

	idx := srt.Renumber()
	if idx != 3 {
		t.Errorf("Expected index is %v, but %v.", 3, idx)
	}
	idx = 0
	for i := 0; i < len(srt.Subtitles); i++ {
		if srt.Subtitles[i].IsDeleted() {
			continue
		}
		idx++
		if srt.Subtitles[i].Num != idx {
			t.Errorf("Expected num is %v, but %v.", idx, srt.Subtitles[i].Num)
		}
	}
}

func TestDeleteSubtitle(t *testing.T) {
	str := `1
0:02:21,860 --> 0:02:24,520

3
0:02:24,520 --> 0:02:27,400

5
0:02:24,520 --> 0:02:27,400

7
0:02:24,520 --> 0:02:27,400`
	srt, err := NewSrtFromString(str)
	if err != nil {
		t.Error(err)
	}

	idx := srt.DeleteSubtitle(5)
	if idx != 2 {
		t.Errorf("Expected index is %v, but %v.", 2, idx)
	}
	if srt.GetCount() != 3 {
		t.Errorf("Expected count of subtitle is %v, but %v.", 3, srt.GetCount())
	}
}

func TestTrimTo(t *testing.T) {
	str := `1
0:00:00,000 --> 0:00:10,000

2
0:00:10,000 --> 0:00:20,000

3
0:00:30,000 --> 0:00:40,000

4
0:00:50,000 --> 0:01:00,000

5
0:01:00,000 --> 0:01:10,000`
	srt, err := NewSrtFromString(str)
	if err != nil {
		t.Error(err)
	}
	ms, err := MillisFromSrtFormat("0:00:40,000")
	if err != nil {
		t.Error(err)
	}
	deleted := srt.TrimTo(ms)
	if deleted != 3 {
		t.Errorf("Expect %v subtitles were deleted, but %v", 3, deleted)
	}
	if srt.GetCount() != 2 {
		t.Errorf("Expect %v subtitles were remaining, but %v", 2, srt.GetCount())
	}
	expectStartMs, _ := MillisFromSrtFormat("0:00:10,000")
	if srt.Subtitles[3].StartMs != expectStartMs {
		t.Errorf("Expect StartMs is %v , but %v", expectStartMs, srt.Subtitles[3].StartMs)
	}
	expectEndMs, _ := MillisFromSrtFormat("0:00:20,000")
	if srt.Subtitles[3].EndMs != expectEndMs {
		t.Errorf("Expect EndMs is %v , but %v", expectEndMs, srt.Subtitles[3].EndMs)
	}
}

func TestCut(t *testing.T) {
	str := `1
0:00:00,000 --> 0:00:10,000

2
0:00:10,000 --> 0:00:20,000

3
0:00:30,000 --> 0:00:40,000

4
0:00:50,000 --> 0:01:00,000

5
0:01:00,000 --> 0:01:10,000

6
0:01:10,000 --> 0:01:20,000`
	srt, err := NewSrtFromString(str)
	if err != nil {
		t.Error(err)
	}
	start, _ := MillisFromSrtFormat("0:00:30,000")
	end, _ := MillisFromSrtFormat("0:01:00,000")
	deleted := srt.Cut(start, end)
	if deleted != 2 {
		t.Errorf("Expect %v subtitles were deleted, but %v", 2, deleted)
	}
}
func TestDeleteEmpty(t *testing.T) {
	str := `1
0:00:00,000 --> 0:00:10,000
aaa

2
0:00:10,000 --> 0:00:20,000

3
0:00:30,000 --> 0:00:40,000


4
0:00:50,000 --> 0:01:00,000
bbb

5
0:01:00,000 --> 0:01:10,000
   

6
0:01:10,000 --> 0:01:20,000`
	srt, err := NewSrtFromString(str)
	if err != nil {
		t.Error(err)
	}
	deleted := srt.DeleteEmpty()
	if deleted != 4 {
		t.Errorf("Expect %v subtitles were deleted, but %v", 4, deleted)
	}
}

func TestDeleteByDuration(t *testing.T) {
	str := `1
0:00:00,000 --> 0:00:10,000
aaa

2
0:00:10,000 --> 0:00:15,000

3
0:00:15,000 --> 0:00:25,000


4
0:00:25,000 --> 0:00:35,000
bbb

5
0:00:35,000 --> 0:00:40,000
   

6
0:00:40,000 --> 0:00:45,000`
	srt, err := NewSrtFromString(str)
	if err != nil {
		t.Error(err)
	}
	deleted := srt.DeleteByDuration(5 * 1000) // 5sec
	if deleted != 3 {
		t.Errorf("Expect %v subtitles were deleted, but %v", 3, deleted)
	}
}

func TestSort(t *testing.T) {
	str := `1
0:00:25,000 --> 0:00:35,000
aaa

2
0:00:15,000 --> 0:00:25,000

3
0:00:10,000 --> 0:00:15,000

4
0:00:00,000 --> 0:00:10,000

5
0:00:35,000 --> 0:00:40,000
   
6
0:00:40,000 --> 0:00:45,000`
	srt, err := NewSrtFromString(str)
	if err != nil {
		t.Error(err)
	}
	srt.Sort()
	idx := srt.Subtitles[0].Num
	if idx != 4 {
		t.Errorf("Expected first subtitle's num is %v but %v.", 4, idx)
	}
}
