package helper

import "time"

const layout = "2006-01-02"

func TimeToString(t time.Time) string {
	return t.Format(layout)
}

func StringToTime(s string) (time.Time, error) {
	return time.Parse(layout, s)
}
