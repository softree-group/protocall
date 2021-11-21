package stapler

import (
	"bufio"
	"strings"
	"time"
)

type Phrase struct {
	Time time.Time
	User string
	Text string
}

func parseString(raw string) []Phrase {
	scanner := bufio.NewScanner(strings.NewReader(raw))
	scanner.Split(bufio.ScanLines)
	res := []Phrase{}
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ":")
		if len(line) < 2 {
			return nil
		}

		time, err := time.Parse(time.RFC3339, line[0])
		if err != nil {
			return nil
		}

		res = append(res, Phrase{
			Time: time,
			User: line[1],
			Text: line[2],
		})
	}
	return res
}
