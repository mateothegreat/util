package dates

import (
	"fmt"
	"time"
)

type DateLayout string

const (
	DateLayoutYYYYMMDD        DateLayout = "2006-01-02"
	DateLayoutYYYYMMDDTHHMMSS DateLayout = "2006-01-02T15:04:05"
)

func Parse(layout DateLayout, in string) (time.Time, error) {
	s, err := time.Parse(string(layout), in)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date: %w", err)
	}
	return s, nil
}
