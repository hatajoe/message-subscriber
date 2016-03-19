package runner

//go:generate stringer -type=Status
type Status int

const (
	Running Status = iota
	Stopped
	Aborted
)
