package runner

import "time"

type Option struct {
	InitialState    Status
	InitialDuration time.Duration
}
