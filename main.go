package main

import (
	"fmt"
	"sync"

	"search-benchmark/claim"
	"search-benchmark/data"
	"search-benchmark/db"

	"github.com/sirupsen/logrus"
)

func main() {
	testData := map[string]map[string]string{
		"resolve channels": data.ChannelsToResolve,
		"resolve streams":  data.StreamsToResolve,
		"resolve titles":   data.TitlesToResolve,
	}
	testedInstance := "dev"
	endpoints := map[string]string{
		"prod":  "https://lighthouse.lbry.com/",
		"dev":   "https://dev.lighthouse.lbry.com/",
		"local": "http://localhost:50005/",
	}
	tolerance := 3
	for desc, t := range testData {
		wg := &sync.WaitGroup{}
		eMB := claim.New(wg, 8, t)
		eMB.SetTolerance(tolerance)
		eMB.SetEndpoint(endpoints[testedInstance])
		eMB.Start()
		wg.Wait()
		for _, e := range eMB.Errors() {
			fmt.Println(e.Error())
		}
		fmt.Println(eMB.Summary())
		err := db.StoreResults(testedInstance, endpoints[testedInstance], desc, db.Results{
			Tolerance:     tolerance,
			InstantRate:   eMB.InstantRate(),
			ThresholdRate: eMB.ThresholdRate(),
			WholesomeRate: eMB.WholesomeRate(),
			Errors:        len(eMB.Errors()),
			Timing:        eMB.Timing().Milliseconds(),
		})
		if err != nil {
			logrus.Errorln(err.Error())
		}
	}
}
