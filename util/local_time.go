package util

import "time"

func LocalTime_getTimeClient(interval int) *time.Time {
	t := time.Now().UTC()
	t.Add(time.Duration(time.Duration.Hours(time.Duration(interval))))

	return &t
}
