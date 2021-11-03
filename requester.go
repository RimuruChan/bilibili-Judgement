package main

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"strconv"
)

func getJuryInfo() (juryInfo, error) {
	resData, err := request("GET", juryURL, make(map[string]string))
	if err != nil {
		return juryInfo{}, err
	}
	var ji juryInfo
	err = mapstructure.Decode(resData.Data, &ji)
	if err != nil {
		return juryInfo{}, err
	}
	return ji, nil
}

func getNext() (string, error) {
	resData, err := request("GET", nextURL, make(map[string]string))
	if err != nil {
		return "", err
	}
	obj := resData.Data["case_id"]
	switch obj.(type) {
	case string:
		return obj.(string), nil
	default:
		return "", errors.New("case id must be string")
	}
}

func getCaseInfo(caseId string) (caseInfo, error) {
	resData, err := request("GET", infoURL, map[string]string{
		"case_id": caseId,
	})
	if err != nil {
		return caseInfo{}, err
	}
	var ci caseInfo
	err = mapstructure.Decode(resData.Data, &ci)
	if err != nil {
		return caseInfo{}, err
	}
	return ci, nil
}

func getOpinion(caseId string, pg int, ps int) (opinion, error) {
	resData, err := request("GET", opinionURL, map[string]string{
		"case_id": caseId,
		"pg":      strconv.Itoa(pg),
		"ps":      strconv.Itoa(ps),
	})
	if err != nil {
		return opinion{}, err
	}
	var o opinion
	err = mapstructure.Decode(resData.Data, &o)
	if err != nil {
		return opinion{}, err
	}
	return o, nil
}

func postVote(caseId string, vote int) (resData, error) {
	return request("POST", voteURL, map[string]string{
		"case_id":   caseId,
		"vote":      strconv.Itoa(vote),
		"content":   "",
		"anonymous": "1",
		"csrf":      csrf,
	})
}

func request(method string, api string, param map[string]string) (resData, error) {
	req, err := http.NewRequest(method, api, nil)
	if err != nil {
		return resData{}, err
	}
	query := req.URL.Query()
	for k, v := range param {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()
	addHeader(&req.Header)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return resData{}, err
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resData{}, err
	}
	var rd resData
	err = json.Unmarshal(bytes, &rd)
	if err != nil {
		return resData{}, err
	}
	if rd.Code != 0 {
		return resData{}, errors.New("return code is not 0: " + rd.Message)
	}
	return rd, nil
}

func addHeader(header *http.Header) {
	header.Add("cookie", "bili_jct="+csrf+"; SESSDATA="+sessdata+";")
	header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.146 Safari/537.36")
}
