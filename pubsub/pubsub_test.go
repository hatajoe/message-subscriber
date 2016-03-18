package pubsub_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/hatajoe/message-subscriber/pubsub"
	"github.com/hatajoe/message-subscriber/runner"
)

var pool *redis.Pool

var (
	server      string = os.Getenv("MSGSUB_REDIS_SERVER")
	maxIdle     string = os.Getenv("MSGSUB_REDIS_MAX_IDLE")
	idleTimeout string = os.Getenv("MSGSUB_REDIS_IDLE_TIMEOUT")
	channel     string = os.Getenv("MSGSUB_REDIS_CHANNEL")
)

func TestMain(m *testing.M) {
	mi, err := strconv.Atoi(maxIdle)
	if err != nil {
		os.Exit(1)
	}
	it, err := strconv.Atoi(idleTimeout)
	if err != nil {
		os.Exit(1)
	}
	pool = &redis.Pool{
		MaxIdle:     mi,
		IdleTimeout: time.Duration(it),
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", server)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	os.Exit(m.Run())
}

func TestPubSub(t *testing.T) {
	expect := []byte("this is test!")
	testEnd := make(chan bool)
	buffer := make(chan []byte)
	sub, err := pubsub.NewPubSub(pool.Get(), channel, buffer)
	if err != nil {
		t.Fatal(err)
	}
	r := runner.NewRunner(runner.Option{InitialState: runner.Running})
	go func() {
		for buf := range buffer {
			if r.GetState() == runner.Aborted {
				break
			}
			if string(buf) != string(expect) {
				t.Errorf("err: '%s' expecting got '%s'\n", string(expect), string(buf))
			}
			r.SetState(runner.Aborted)
		}
	}()
	go func() {
		r.Run(sub)
		testEnd <- true
	}()

	if err := publish(t, r, sub, testEnd, expect); err != nil {
		t.Fatal(err)
	}
}

func publish(t *testing.T, r *runner.Runner, sub *pubsub.PubSub, testEnd chan bool, expect []byte) error {
	c := pool.Get()
	for {
		if err := c.Send("PUBLISH", channel, expect); err != nil {
			return err
		}
		select {
		case <-testEnd:
			if err := sub.UnSubscribe(); err != nil {
				return err
			}
			return nil
		default:
		}
	}

}
