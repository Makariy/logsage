package models

import "time"

type ModelID uint
type AuthToken string

type TimeIntervalStep string

var timeIntervalToTimeMap = map[TimeIntervalStep]time.Duration{
	"day":   24 * time.Hour,
	"week":  24 * 7 * time.Hour,
	"month": 24 * 31 * time.Hour,
	"year":  24 * 31 * 365 * time.Hour,
}

func IsTimeIntervalDefined(interval string) bool {
	_, isDefined := timeIntervalToTimeMap[TimeIntervalStep(interval)]
	return isDefined
}

func ConvertTimeStepToTime(step TimeIntervalStep) int64 {
	duration, exists := timeIntervalToTimeMap[step]
	if !exists {
		panic("time interval step is not defined")
	}
	return int64(duration / time.Second)
}
