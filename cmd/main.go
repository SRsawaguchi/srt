package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/SRsawaguchi/srt"
)

var (
	filePath = flag.String("f", "", "input SRT file.")
	cmd      = flag.String("c", "dump", "dump SRT file.")
	time     = flag.String("t", "00:00:00,000", "Time to trim from head of file.")
	between  = flag.String("b", "00:00:00,000-00:00:00,000", "start-end")
)

func main() {
	var err error = nil
	flag.Parse()

	switch *cmd {
	case "dump":
		err = handleDump()
	case "renumber":
		err = handleRenumber()
	case "trim":
		err = handleTrim()
	case "cut":
		err = handleCut()
	case "delete_empty":
		err = handleDeleteEmpty()
	case "delete_by_duration":
		err = handleDeleteByDuration()
	case "sort":
		err = handleSort()
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		os.Exit(1)
	}

	// srt, err := srt.NewSrtFromString(string(bytes))
}

func readFile(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func handleDump() error {
	strSrt, err := readFile(*filePath)
	if err != nil {
		return err
	}

	fmt.Println(strSrt)
	return nil
}

func handleRenumber() error {
	s, err := srt.NewSrtFromFile(*filePath)
	if err != nil {
		return err
	}
	s.Renumber()
	fmt.Println(s.String())
	return nil
}

func handleTrim() error {
	s, err := srt.NewSrtFromFile(*filePath)
	if err != nil {
		return err
	}

	ms, err := srt.MillisFromSrtFormat(*time)
	if err != nil {
		return err
	}
	s.TrimTo(ms)
	s.Renumber()
	fmt.Println(s.String())
	return nil
}

func handleCut() error {
	s, err := srt.NewSrtFromFile(*filePath)
	if err != nil {
		return err
	}

	duration := strings.Split(*between, "-")
	start, err := srt.MillisFromSrtFormat(duration[0])
	if err != nil {
		return err
	}
	end, err := srt.MillisFromSrtFormat(duration[1])
	if err != nil {
		return err
	}
	s.Cut(start, end)
	s.Renumber()
	fmt.Println(s.String())
	return nil
}

func handleDeleteEmpty() error {
	s, err := srt.NewSrtFromFile(*filePath)
	if err != nil {
		return err
	}

	s.DeleteEmpty()
	s.Renumber()
	fmt.Println(s.String())
	return nil
}

func handleDeleteByDuration() error {
	s, err := srt.NewSrtFromFile(*filePath)
	if err != nil {
		return err
	}

	ms, err := srt.MillisFromSrtFormat(*time)
	if err != nil {
		return err
	}
	s.DeleteByDuration(ms)
	s.Renumber()
	fmt.Println(s.String())
	return nil
}

func handleSort() error {
	s, err := srt.NewSrtFromFile(*filePath)
	if err != nil {
		return err
	}

	s.Sort()
	s.Renumber()
	fmt.Println(s.String())
	return nil
}
