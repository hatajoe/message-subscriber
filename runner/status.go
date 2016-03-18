package runner

//go:generate stringer -type=Status
type Status int

const (
	Stopped Status = iota
	Running
	Aborted
)
