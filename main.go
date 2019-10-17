package main

import (
	"fmt"
	"sync"

	"search-benchmark/channel"
)

func main() {
	wg := &sync.WaitGroup{}
	//benchmarks := make([]benchmark, 1)
	eMB := channel.New(wg, 8)
	//benchmarks[0] = eMB
	eMB.SetTolerance(3)
	eMB.Start()
	wg.Wait()
	fmt.Println(eMB.Summary())
	//fmt.Printf("Channels found in search resuts: %d/%d (%.2f%%)\n", found.Load(), len(data.ChannelsToResolve), float64(found.Load())/float64(len(data.ChannelsToResolve))*100.0)
}
