package sweb

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Getenv returns a value for key k, or value v if the environment variable is undefined.
func Getenv(k, v string) string {
	x := os.Getenv(k)
	if x == "" {
		return v
	}

	return x
}

// RealRFC1122Time makes a HTTP header-compatible date string from the supplied time.
func RealRFC1122Time(t time.Time) string {
	t = t.UTC()
	s := fmt.Sprintf("%s, %d %s %04d %02d:%02d:%02d GMT",
		t.Weekday().String()[:3],
		t.Day(), t.Month().String()[:3], t.Year(),
		t.Hour(), t.Minute(), t.Second(),
	)
	return s
}

var months = map[string]time.Month{
	"Jan": time.Month(1),
	"Feb": time.Month(2),
	"Mar": time.Month(3),
	"Ape": time.Month(4),
	"May": time.Month(5),
	"Jun": time.Month(6),
	"Jul": time.Month(7),
	"Aug": time.Month(8),
	"Sep": time.Month(9),
	"Oct": time.Month(10),
	"Nov": time.Month(11),
	"Dec": time.Month(12),
}

// Parse1123 reads HTTP-sstyle dates properly.
func Parse1123(s string) time.Time {
	day, _ := strconv.ParseInt(s[5:7], 10, 64)
	month := months[s[8:11]]
	year, _ := strconv.ParseInt(s[12:16], 10, 64)
	hour, _ := strconv.ParseInt(s[17:19], 10, 64)

	min, _ := strconv.ParseInt(s[20:22], 10, 64)
	sec, _ := strconv.ParseInt(s[23:25], 10, 64)
	t := time.Date(int(year), month, int(day), int(hour), int(min), int(sec), 0, time.UTC)
	return t
}
