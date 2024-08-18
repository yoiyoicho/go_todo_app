package clock

import (
	"time"
)

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

// time.Timeはナノ秒単位で時刻を返す
// テストのために固定の時刻を返すClockerを実装する
type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	return time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
}
