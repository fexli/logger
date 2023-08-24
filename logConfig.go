package logger

import (
	"os"
	"time"
)

type singleTime struct {
	Hour   int
	Min    int
	Second int
	TZone  *time.Location
}
type SliceTime []singleTime

func NewSliceTime(hour, min, sec int, timeZone ...*time.Location) *SliceTime {
	if len(timeZone) == 0 {
		return &SliceTime{
			{
				Hour:   hour,
				Min:    min,
				Second: sec,
				TZone:  time.Local,
			},
		}
	}
	return &SliceTime{
		{
			Hour:   hour,
			Min:    min,
			Second: sec,
			TZone:  timeZone[0],
		},
	}
}

func (s *SliceTime) And(hour, min, sec int, timeZone ...*time.Location) *SliceTime {
	if len(timeZone) == 0 {
		*s = append(*s, singleTime{
			Hour:   hour,
			Min:    min,
			Second: sec,
			TZone:  time.Local,
		})
		return s
	}
	*s = append(*s, singleTime{
		Hour:   hour,
		Min:    min,
		Second: sec,
		TZone:  timeZone[0],
	})
	return s
}

type Config struct {
	// 日志名称
	Name string
	// 日志插入首行描述
	Desc string
	// 过时日志存储路径
	OldLogPath func(info os.FileInfo) string
	// 日志最大存储时间（超时切片，优先生效）
	MaxLogTime time.Duration
	// 日志定时切片
	SliceWhen SliceTime
}

func findNextByWhen(cfg Config) time.Duration {
	if len(cfg.SliceWhen) != 0 {
		now := time.Now()
		minD := time.Hour * 24
		for _, v := range cfg.SliceWhen {
			t := time.Date(now.Year(), now.Month(), now.Day(), v.Hour, v.Min, v.Second, 0, v.TZone)
			if t.After(now) {
				minD = min(minD, t.Sub(now))
			} else {
				minD = min(minD, t.Add(time.Hour*24).Sub(now))
			}
		}
		return minD
	}
	if cfg.MaxLogTime != 0 {
		return cfg.MaxLogTime
	}
	return time.Hour * 24
}

func min(d time.Duration, sub time.Duration) time.Duration {
	if d > sub {
		return sub
	}
	return d
}
