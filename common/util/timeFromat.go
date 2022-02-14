package util

import (
	"time"
)

func MySQLTimeToGoTime(t string) (string, error) {
	t1, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return "", err
	}
	t = t1.Format("2006-01-02 15:04:05")
	return t, err
}

func NowTimeWithCTS() time.Time {
	t1 := time.Unix(time.Now().Unix(), 0).Local()
	return t1
}
