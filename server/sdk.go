package server

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"
)

type SlowTestRun struct {
	exit chan struct{}
}

var r *SlowTestRun
var l sync.Once

func GetSlowTestRunner() *SlowTestRun {

	l.Do(func() {
		r = &SlowTestRun{
			exit: make(chan struct{}, 0),
		}
	})
	return r
}

func (r *SlowTestRun) sendExit() {
	log.Println("send exit waitExit")
	close(r.exit)
}

func (r *SlowTestRun) waitExit() {

	log.Println("waitExit start")
	<-r.exit
	log.Println("waitExit ok")
}

func (r *SlowTestRun) RunCase(t *testing.T, f func()) {
	if os.Getenv("SLOWTEST") != "true" {
		log.Println("skip slowtest,env SLOWTEST=" + os.Getenv("SLOWTEST"))
		t.Skip()
		return
	}
	log.Println("start run case")
	go f()
	r.waitExit()
	log.Println("end run case")
}

func (r *SlowTestRun) WaitSignal(m *testing.M) {

	if os.Getenv("SLOWTEST") != "true" {
		log.Println("skip slowtest,env SLOWTEST=" + os.Getenv("SLOWTEST"))
		return
	}
	log.Println("start testing.M")
	go func() {
		code := m.Run()
		log.Println("test exit code ", code)
	}()

	log.Println("start WaitSignal")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range c {
			log.Println("get sig", sig)
			r.sendExit()
		}
	}()

	r.waitExit()
	time.Sleep(time.Second)
	log.Println("end WaitSignal")

}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func RunCase(t *testing.T, f func()) {
	GetSlowTestRunner().RunCase(t, f)
}
func WaitSignal(m *testing.M) {
	GetSlowTestRunner().WaitSignal(m)
}
