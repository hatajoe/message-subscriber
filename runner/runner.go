package runner

import (
	"log"
	"time"

	"sync"

	"github.com/hatajoe/message-subscriber"
)

type Runner struct {
	status        Status
	sleepDuration time.Duration
	mu            sync.Mutex
}

func NewRunner(option Option) *Runner {
	return &Runner{
		status:        option.InitialState,
		sleepDuration: option.SleepDuration,
	}
}

func (m *Runner) GetState() Status {
	m.mu.Lock()
	status := m.status
	m.mu.Unlock()
	return status
}

func (m *Runner) SetState(st Status) {
	m.mu.Lock()
	m.status = st
	m.mu.Unlock()
}

func (m *Runner) Run(sub subscriber.Subscriber) {
	for {
		switch m.GetState() {
		case Stopped:
			time.Sleep(m.sleepDuration)
		case Running:
			m.run(sub)
		case Aborted:
			return
		}
	}
}

func (m *Runner) run(sub subscriber.Subscriber) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if err := sub.Subscribe(); err != nil {
		if err := sub.Abort(); err != nil {
			panic(err)
		}
		return
	}
	if err := sub.End(); err != nil {
		panic(err)
	}
}
