package strategy

import "time"

type Strategy interface {
	Next() (time.Duration, bool)
	NextWithRetriedTimes(retriedTimes int32) (time.Duration, bool)
}
