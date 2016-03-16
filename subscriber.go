package subscriber

// Subscriber is message subscriber interface for Runner
type Subscriber interface {
	Subscribe() error
	Abort() error
	End() error
}
