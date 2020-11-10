package slowtestsdk

import (
	"log"
	"net/http"
	"sync"
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

func (r *SlowTestRun) Exit() {
	log.Println("send exit wait")
	close(r.exit)
}

func (r *SlowTestRun) Wait() {

	log.Println("wait start")
	<-r.exit
	log.Println("wait ok")
}

func (r *SlowTestRun) RunCase(f func()) {
	log.Println("start run case")
	go f()
	r.Wait()
	log.Println("end run case")
}

func (r *SlowTestRun) StartServe() {
	log.Println("StartServe")

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/exit", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(200)
			_, _ = w.Write([]byte("exit ok \n"))
			r.Exit()
		})
		err := http.ListenAndServe("0.0.0.0:6200", mux)
		if err != nil {
			log.Fatal(err)
		}
	}()
	r.Wait()
	time.Sleep(time.Second)
	log.Println("stopServe")

}
