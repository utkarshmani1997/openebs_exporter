package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func getVolStats(URL string, obj interface{}) error {

	resp, err := http.Get(URL)
	if err != nil {
		return err
	}

	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	if err2 != nil {
		return err2
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(obj); err != nil {
		return err
	}
	return nil
}
