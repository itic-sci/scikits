package scikits

import (
	"strings"
	"time"
)

func GetNowStr() string {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	return timeStr
}

func GetNowStamp() int64 {
	timeStamp := time.Now().Unix()
	return timeStamp
}

func TimestampToStr(timestamp int64) string {
	timeNow := time.Unix(timestamp, 0)                  //2017-08-30 16:19:19 +0800 CST
	timeString := timeNow.Format("2006-01-02 15:04:05") //2015-06-15 08:52:32
	return timeString
}

// 2006-01-02
func TimestampToStrDay(timestamp int64) string {
	timeNow := time.Unix(timestamp, 0)         //2017-08-30 16:19:19 +0800 CST
	timeString := timeNow.Format("2006-01-02") //2015-06-15 08:52:32
	return timeString
}

func StrTimeToTime(timeStr string) time.Time {
	if !strings.Contains(timeStr, ":") {
		timeStr = strings.TrimSpace(timeStr) + " 00:00:00"
	}
	formatTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		panic(err)
	}
	return formatTime
}

func TimeToStr(timeVal time.Time) string {
	return timeVal.Format("2006-01-02 15:04:05")
}

func GetTimeBeforeDay(day int) time.Time {
	nowTime := time.Now()
	oldTime := nowTime.AddDate(0, 0, -day)
	return oldTime
}
