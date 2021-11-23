package stapler

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

type Phrase struct {
	Time time.Time
	User string
	Text string
}

func parse(raw string) ([]Phrase, error) {
	scanner := bufio.NewScanner(strings.NewReader(raw))
	scanner.Split(bufio.ScanLines)
	res := []Phrase{}
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		line := strings.Split(scanner.Text(), ":")
		if len(line) < 2 {
			return nil, fmt.Errorf("invalid line format")
		}

		time, err := time.Parse(time.RFC3339, line[0])
		if err != nil {
			return nil, err
		}

		res = append(res, Phrase{
			Time: time,
			User: line[1],
			Text: line[2],
		})
	}
	return res, nil
}
