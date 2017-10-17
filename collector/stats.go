package collector

import (
	"encoding/json"
	"net/http"
)

type statsResponse struct {
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
	Data       stats  `json:"data"`
}

type stats struct {
	VolStats []*volume
}

type volume struct {
	RevisionCounter int64         `json:"RevisionCounter"`
	ReplicaCounter  int64         `json:"ReplicaCounter"`
	SCSIIOCount     map[int]int64 `json:"SCSIIOCount"`

	ReadIOPS            string `json:"ReadIOPS"`
	TotalReadTime       string `json:"TotalReadTime"`
	TotalReadBlockCount string `json:"TotalReadBlockCount"`

	WriteIOPS            string `json:"WriteIOPS"`
	TotalWriteTime       string `json:"TotalWriteTime"`
	TotalWriteBlockCount string `json:"TotatWriteBlockCount"`

	UsedLogicalBlocks string `json:"UsedLogicalBlocks"`
	UsedBlocks        string `json:"UsedBlocks"`
	SectorSize        string `json:"SectorSize"`
}

func getVolStats(URL string) (*stats, error) {

	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sr statsResponse
	if err = json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, err
	}
	return &sr.Data, nil
}
