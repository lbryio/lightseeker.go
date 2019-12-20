package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lbryio/lbry.go/v2/extras/errors"
)

type LightHouseEngine struct {
	endpoint string
}

type lightHouseResponse []struct {
	Name    string `json:"name"`
	ClaimID string `json:"claimId"`
}

func NewLightHouseEngine(endpoint string) *LightHouseEngine {
	return &LightHouseEngine{endpoint: endpoint}
}

func (lh *LightHouseEngine) GetEndpoint() string {
	return lh.endpoint
}

func (lh *LightHouseEngine) Query(terms string) (SearchResponse, error) {
	searchURL := fmt.Sprintf("%ssearch?s=%s&size=20", lh.GetEndpoint(), url.QueryEscape(terms))
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, errors.Err(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Err(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, errors.Err(err)
	}
	var resp lightHouseResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, errors.Err(err)
	}
	searchResponse := make(SearchResponse, len(resp))
	for i, v := range resp {
		searchResponse[i].ClaimID = v.ClaimID
		searchResponse[i].ClaimName = v.Name
	}
	return searchResponse, nil
}

func (lh *LightHouseEngine) Version() (*SearchVersion, error) {
	statusAPI := lh.endpoint + "status"
	r, err := http.NewRequest(http.MethodGet, statusAPI, nil)
	if err != nil {
		return nil, errors.Err(err)
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, errors.Err(err)

	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, errors.Err(err)
	}
	type statusResponse struct {
		Version         string `json:"Version"`
		SemanticVersion string `json:"SemanticVersion"`
		VersionMsg      string `json:"VersionMsg"`
	}
	sr := statusResponse{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, errors.Err(err)
	}
	return &SearchVersion{
		SemVer:     sr.SemanticVersion,
		CommitHash: sr.Version,
	}, nil
}
