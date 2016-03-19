package runner

import "time"

type Option struct {
	InitialState  Status
	SleepDuration time.Duration
}
