package main

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"strconv"
)

func getJuryInfo() (JuryInfo, error) {
	resData, err := request("GET", JuryUrl, make(map[string]string))
	if err != nil {
		return JuryInfo{}, err
	}
	var juryInfo JuryInfo
	err = mapstructure.Decode(resData.Data, &juryInfo)
	if err != nil {
		return JuryInfo{}, err
	}
	return juryInfo, nil
}

func getNext() (string, error) {
	resData, err := request("GET", NextUrl, make(map[string]string))
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

func getCaseInfo(caseId string) (CaseInfo, error) {
	resData, err := request("GET", InfoUrl, map[string]string{
		"case_id": caseId,
	})
	if err != nil {
		return CaseInfo{}, err
	}
	var caseInfo CaseInfo
	err = mapstructure.Decode(resData.Data, &caseInfo)
	if err != nil {
		return CaseInfo{}, err
	}
	return caseInfo, nil
}

func getOpinion(caseId string, pg int, ps int) (Opinion, error) {
	resData, err := request("GET", OpinionUrl, map[string]string{
		"case_id": caseId,
		"pg":      strconv.Itoa(pg),
		"ps":      strconv.Itoa(ps),
	})
	if err != nil {
		return Opinion{}, err
	}
	var opinion Opinion
	err = mapstructure.Decode(resData.Data, &opinion)
	if err != nil {
		return Opinion{}, err
	}
	return opinion, nil
}

func postVote(caseId string, vote int) (ResData, error) {
	return request("POST", VoteUrl, map[string]string{
		"case_id":   caseId,
		"vote":      strconv.Itoa(vote),
		"content":   "",
		"anonymous": "1",
		"csrf":      csrf,
	})
}

func request(method string, api string, param map[string]string) (ResData, error) {
	req, err := http.NewRequest(method, api, nil)
	if err != nil {
		return ResData{}, err
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
		return ResData{}, err
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ResData{}, err
	}
	var resData ResData
	err = json.Unmarshal(bytes, &resData)
	if err != nil {
		return ResData{}, err
	}
	if resData.Code != 0 {
		return ResData{}, errors.New("return code is not 0: " + resData.Message)
	}
	return resData, nil
}

func addHeader(header *http.Header) {
	header.Add("cookie", "bili_jct="+csrf+"; SESSDATA="+sessdata+";")
	header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.146 Safari/537.36")
}
