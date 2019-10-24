package main

import (
	"fmt"
	"sync"

	"search-benchmark/claim"
	"search-benchmark/data"
)

func main() {
	testData := []map[string]string{
		//data.ChannelsToResolve,
		//data.StreamsToResolve,
		data.TitlesToResolve,
	}
	for _, t := range testData {
		wg := &sync.WaitGroup{}
		eMB := claim.New(wg, 8, t)
		eMB.SetTolerance(3)
		eMB.Start()
		wg.Wait()
		for _, e := range eMB.Errors() {
			fmt.Println(e.Error())
		}
		fmt.Println(eMB.Summary())
	}
}
