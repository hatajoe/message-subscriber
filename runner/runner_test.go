package runner_test

import (
	"testing"

	"github.com/hatajoe/message-subscriber/runner"
)

type TestSubscriber struct {
	message chan map[string]interface{}
	abort   bool
	end     bool
}

func NewTestSubscriber(message chan map[string]interface{}) *TestSubscriber {
	return &TestSubscriber{
		message: message,
		abort:   false,
		end:     false,
	}
}

func (m *TestSubscriber) Subscribe() error {
	m.message <- map[string]interface{}{
		"key1": "value",
		"key2": 1,
		"key3": true,
	}
	return nil
}

func (m *TestSubscriber) UnSubscribe() error {
	return nil
}

func (m *TestSubscriber) Abort() error {
	m.abort = true
	return nil
}

func (m *TestSubscriber) End() error {
	m.end = true
	return nil
}

func TestRun(t *testing.T) {
	message := make(chan map[string]interface{})
	r := runner.NewRunner(runner.Option{InitialState: runner.Running})
	go func() {
		for m := range message {
			if m, ok := m["key1"]; !ok {
				t.Error("err: map expect has a key `key1'")
			} else {
				if v, ok := m.(string); !ok {
					t.Error("err: `key1' expecting type string")
				} else {
					if v != "value" {
						t.Error("err: `key1' value expecting `value'")
					}
				}
			}
			if m, ok := m["key2"]; !ok {
				t.Error("err: map expect has a key `key2'")
			} else {
				if v, ok := m.(int); !ok {
					t.Error("err: `key2' expecting type int")
				} else {
					if v != 1 {
						t.Error("err: `key2' value expecting 1")
					}
				}

			}
			if m, ok := m["key3"]; !ok {
				t.Error("err: map expect has a key `key2'")
			} else {
				if v, ok := m.(bool); !ok {
					t.Error("err: `key3' expecting type bool")
				} else {
					if !v {
						t.Error("err: `key3' value expecting true")
					}
				}

			}
			r.SetState(runner.Aborted)
			break
		}
	}()
	sub := NewTestSubscriber(message)
	r.Run(sub)
}
