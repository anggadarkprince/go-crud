package utilities

import "time"

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func FormatDateHuman(t time.Time) string {
	return t.Format("02 January 2006")
}

func FormatDateTimeHuman(t time.Time) string {
	return t.Format("02 January 2006 15:04:05")
}

func NowUnix() int64 {
	return time.Now().Unix()
}