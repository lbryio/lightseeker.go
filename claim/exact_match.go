package claim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/uber-go/atomic"
)

type taskData struct {
	searchTerm string
	claimID    string
	searchURL  string
}

type errorStack struct {
	errors []error
	lock   *sync.Mutex
}

func (e *errorStack) Append(err error) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.errors = append(e.errors, err)
}

type ExactMatchBenchmark struct {
	wg               *sync.WaitGroup
	startTime        time.Time
	instaMatches     *atomic.Int32
	thresholdMatches *atomic.Int32
	matches          *atomic.Int32
	tolerance        int
	workers          int
	work             chan taskData
	errors           *errorStack
	runTime          time.Duration
	data             map[string]string
	endPoint         string
}

func New(wg *sync.WaitGroup, workers int, data map[string]string) *ExactMatchBenchmark {
	return &ExactMatchBenchmark{
		wg:   wg,
		work: make(chan taskData),
		errors: &errorStack{
			errors: []error{},
			lock:   &sync.Mutex{},
		},
		instaMatches:     atomic.NewInt32(0),
		thresholdMatches: atomic.NewInt32(0),
		matches:          atomic.NewInt32(0),
		workers:          workers,
		data:             data,
		endPoint:         "https://dev.lighthouse.lbry.com/",
	}
}

func (e *ExactMatchBenchmark) Timing() time.Duration {
	return e.runTime
}

func (e *ExactMatchBenchmark) Summary() string {
	instaRate := e.Rate() * 100
	thresholdRate := float64(e.thresholdMatches.Load()) / float64(len(e.data)) * 100
	wholesomeRate := float64(e.matches.Load()) / float64(len(e.data)) * 100
	return fmt.Sprintf(`Instant match rate: %.2f
Threshold match rate: %.2f
Wholesome match rate: %.2f
Errors: %d
Timing: %s`, instaRate, thresholdRate, wholesomeRate, len(e.errors.errors), e.Timing().String())
}

func (e *ExactMatchBenchmark) Errors() []error {
	return e.errors.errors
}

func (e *ExactMatchBenchmark) Rate() float64 {
	return float64(e.instaMatches.Load()) / float64(len(e.data))
}

func (e *ExactMatchBenchmark) Start() {
	e.startTime = time.Now()

	for i := 0; i < e.workers; i++ {
		e.wg.Add(1)
		go e.consume(i)
	}
	e.produce()
	fmt.Println("done producing")
	go func() {
		e.wg.Wait()
		e.runTime = time.Since(e.startTime)
	}()
}

func (e *ExactMatchBenchmark) SetTolerance(t int) {
	e.tolerance = t
}

func (e *ExactMatchBenchmark) produce() {
	for searchTerm, claimID := range e.data {
		e.work <- taskData{
			searchTerm: searchTerm,
			claimID:    claimID,
			searchURL:  fmt.Sprintf("%ssearch?s=%s&size=20", e.endPoint, url.QueryEscape(searchTerm)),
		}
	}
	close(e.work)
}

func (e *ExactMatchBenchmark) consume(worker int) {
	fmt.Printf("starting worker %d\n", worker)
	defer e.wg.Done()
outer:
	for {
		s, more := <-e.work
		if !more {
			return
		}

		req, err := http.NewRequest("GET", s.searchURL, nil)
		if err != nil {
			e.errors.Append(errors.Err(err))
			continue
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			e.errors.Append(errors.Err(err))
			continue
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			e.errors.Append(errors.Err(err))
			continue
		}
		var searchResponse []struct {
			Name    string `json:"name"`
			ClaimID string `json:"claimId"`
		}
		err = json.Unmarshal(body, &searchResponse)
		if err != nil {
			e.errors.Append(errors.Err(err))
			continue
		}
		for i, r := range searchResponse {
			if r.ClaimID == s.claimID {
				if i == 0 {
					e.instaMatches.Add(1)
				}
				if i < e.tolerance {
					e.thresholdMatches.Add(1)
				}
				e.matches.Add(1)
				continue outer
			}
		}
		fmt.Printf("no results for %s - %s\n", s.searchTerm, s.claimID) //todo: export it to a value rather than printing it
	}
}

func (e *ExactMatchBenchmark) SetEndpoint(endpoint string) {
	e.endPoint = endpoint
}
