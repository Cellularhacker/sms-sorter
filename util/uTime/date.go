package uTime

import (
	"fmt"
	"strconv"
	"time"
)

func GetNowDate() int {
	loc, _ := time.LoadLocation("Asia/Seoul")
	year, month, day := time.Now().In(loc).Date()
	str := fmt.Sprintf("%4d%2d%2d", year, month, day)
	i, _ := strconv.Atoi(str)
	return i
}

func GetKSTDateStrBeautify(t *time.Time) string {
	if t == nil {
		tt := time.Now()
		t = &tt
	}
	loc, _ := time.LoadLocation("Asia/Seoul")
	year, month, day := t.In(loc).Date()
	hour, minute, second := t.In(loc).Clock()

	return fmt.Sprintf("%04d-%02d-%02d_%02d:%02d:%02d", year, int(month), day, hour, minute, second)
}