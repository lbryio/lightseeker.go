package engine

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/lbryio/lbry.go/v2/extras/errors"
)

type SDKEngine struct {
	endpoint string
}

type SDKResponse struct {
	Result *struct {
		Items []struct {
			ClaimID string `json:"claim_id"`
			Name    string `json:"name"`
		} `json:"items"`
	} `json:"result"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

type SDKBody struct {
	Method string    `json:"method"`
	Params SDKParams `json:"params"`
}
type SDKParams struct {
	Text     string `json:"text"`
	PageSize int    `json:"page_size"`
}

func NewSDKEngine(endpoint string) *SDKEngine {
	return &SDKEngine{endpoint: endpoint}
}

func (s *SDKEngine) GetEndpoint() string {
	return s.endpoint
}

func (s *SDKEngine) Query(terms string) (SearchResponse, error) {
	searchBody := SDKBody{
		Method: "claim_search",
		Params: SDKParams{
			Text:     terms,
			PageSize: 20,
		},
	}
	sb, err := json.Marshal(searchBody)
	req, err := http.NewRequest(http.MethodPost, s.GetEndpoint(), bytes.NewReader(sb))
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
	var resp SDKResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, errors.Err(err)
	}
	if resp.Error != nil {
		return nil, errors.Err(resp.Error.Message)
	}
	searchResponse := make(SearchResponse, len(resp.Result.Items))
	for i, v := range resp.Result.Items {
		searchResponse[i].ClaimID = v.ClaimID
		searchResponse[i].ClaimName = v.Name
	}
	return searchResponse, nil
}

func (s *SDKEngine) Version() (*SearchVersion, error) {
	statusBody := struct {
		Method string `json:"method"`
	}{
		Method: "version",
	}
	sb, err := json.Marshal(statusBody)

	r, err := http.NewRequest(http.MethodPost, s.GetEndpoint(), bytes.NewReader(sb))
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
		Result *struct {
			LbrynetVersion string `json:"lbrynet_version"`
		}
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	sr := statusResponse{}
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, errors.Err(err)
	}
	if sr.Error != nil {
		return nil, errors.Err(sr.Error.Message)
	}
	return &SearchVersion{
		SemVer:     sr.Result.LbrynetVersion,
		CommitHash: "unknown",
	}, nil
}
