package main

import (
	"fmt"
	"sync"

	"search-benchmark/claim"
	"search-benchmark/data"
	"search-benchmark/db"
	"search-benchmark/engine"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	testData := map[string]map[string]string{
		"resolve channels": data.ChannelsToResolve,
		"resolve streams":  data.StreamsToResolve,
		"resolve titles":   data.TitlesToResolve,
	}
	testedInstance := "sdk"
	endpoints := map[string]engine.SearchEngine{
		"prod":  engine.NewLightHouseEngine("https://lighthouse.lbry.com/"),
		"dev":   engine.NewLightHouseEngine("https://dev.lighthouse.lbry.com/"),
		"local": engine.NewLightHouseEngine("http://localhost:50005/"),
		"sdk":   engine.NewSDKEngine("http://localhost:5279"),
	}
	tolerance := 3
	for desc, t := range testData {
		wg := &sync.WaitGroup{}
		eMB := claim.New(wg, 8, t)
		eMB.SetTolerance(tolerance)
		eMB.SetEngine(endpoints[testedInstance])
		eMB.Start()
		wg.Wait()
		for _, e := range eMB.Errors() {
			fmt.Println(errors.FullTrace(e))
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
