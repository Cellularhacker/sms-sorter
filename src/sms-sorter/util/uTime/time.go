package uTime

import "time"

var loc *time.Location

func init() {
	loc, _ = time.LoadLocation("Asia/Seoul")
}

func GetLoc() *time.Location {
	return loc
}

func GetKST(t *time.Time) time.Time {
	if t == nil {
		return time.Now().In(loc)
	}

	return t.In(loc)
}
