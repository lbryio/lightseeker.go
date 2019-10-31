package main

import (
	"fmt"
	"sync"

	"search-benchmark/claim"
	"search-benchmark/data"
)

func main() {
	testData := []map[string]string{
		data.ChannelsToResolve,
		data.StreamsToResolve,
		data.TitlesToResolve,
	}
	const prod = "https://lighthouse.lbry.com/"
	const dev = "https://dev.lighthouse.lbry.com/"
	const local = "http://localhost:50005/"
	for _, t := range testData {
		wg := &sync.WaitGroup{}
		eMB := claim.New(wg, 8, t)
		eMB.SetTolerance(3)
		eMB.SetEndpoint(local)
		eMB.Start()
		wg.Wait()
		for _, e := range eMB.Errors() {
			fmt.Println(e.Error())
		}
		fmt.Println(eMB.Summary())
	}
}
