package timeutils

import "time"

func Millisecs() time.Duration {
	now := time.Now()
	return time.Duration(now.UnixNano() / 1000000)
}

func Seconds() time.Duration {
	now := time.Now()
	return time.Duration(now.Unix())
}
