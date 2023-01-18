package main

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func getJuryInfo() (JuryInfo, error) {
	resData, err := request("GET", JuryURL, make(map[string]string))
	if err != nil {
		return JuryInfo{}, err
	}
	var ji JuryInfo
	err = mapstructure.Decode(resData.Data, &ji)
	if err != nil {
		return JuryInfo{}, err
	}
	return ji, nil
}

func getNext() (string, error) {
	resData, err := request("GET", NextURL, make(map[string]string))
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
	resData, err := request("GET", InfoURL, map[string]string{
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

func getFullOpinion(caseId string) opinion {
	opinion, err := getOpinion(caseId, 1, 20)
	if err != nil {
		logger.Fatal("failed to fetch opinions:", err.Error())
	}
	pg := 2
	for len(opinion.List) < opinion.Total {
		nextPageOpinion, err := getOpinion(caseId, pg, 20)
		if err != nil {
			logger.Fatal("failed to fetch opinions:", err.Error())
		}
		opinion.List = append(opinion.List, nextPageOpinion.List...)
		pg++
	}
	return opinion
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

func postVote(caseId string, vote int) (ResData, error) {
	rand.Seed(time.Now().UnixNano())
	insiders := strconv.Itoa(rand.Intn(2))
	return request("POST", voteURL, map[string]string{
		"case_id":   caseId,
		"vote":      strconv.Itoa(vote),
		"insiders":  insiders,
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
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return ResData{}, err
	}
	logger.Debug(string(bytes))
	var rd ResData
	err = json.Unmarshal(bytes, &rd)
	if err != nil {
		return ResData{}, err
	}
	if rd.Code != 0 {
		return ResData{}, errors.New("return code is not 0: " + rd.Message)
	}
	return rd, nil
}

func addHeader(header *http.Header) {
	header.Add("cookie", "bili_jct="+csrf+"; SESSDATA="+sessdata+";")
	header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.146 Safari/537.36")
}
