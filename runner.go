package runner

import (
	"log"
	"time"
)

type Runner struct {
	option Option
	status chan Status
}

func NewRunner(opts Option) *Runner {
	return &Runner{
		option: opts,
		status: make(chan Status),
	}
}

func (m *Runner) ChangeState(st Status) {
	m.status <- st
}

func (m *Runner) Run(sub Subscriber) {
	status := m.option.InitialState
	duration := m.option.InitialDuration
	for {
		select {
		case st := <-m.status:
			status = st
		default:
			switch status {
			case Running:
				m.run(sub)
			case Stopped:
				time.Sleep(duration)
			case Aborted:
				return
			}
		}
	}
}

func (m *Runner) run(sub Subscriber) {
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
