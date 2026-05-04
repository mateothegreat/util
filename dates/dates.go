package dates

import (
	"log"
	"time"
)

func GetTime(str string) time.Time {
	layout := "2006-01-02T15:04:05-07:00"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Printf("error parsing time from %s: %v", str, err)
	}

	return t
}
