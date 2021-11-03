package main

import (
	"github.com/withmandala/go-log"
	"os"
	"strings"
	"time"
)

// Version of Program
const Version = "1.0"

var sessdata string
var csrf string

var logger = log.New(os.Stderr)

func main() {
	logger.Info("Version\t", Version, "\tRunning...")
	split()
	logger.Info("Arguments")
	split()
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("no specific csrf")
		}
	}()
	csrf = os.Args[1]
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("no specific sessdata")
		}
	}()
	sessdata = os.Args[2]
	logger.Info("SESSDATA:\t", sessdata)
	logger.Info("CRSF:\t", csrf)
	split()
	juryInfo, err := getJuryInfo()
	if err != nil {
		logger.Fatal("failed to fetch jury info:", err.Error())
	}
	logger.Info("Account Info")
	split()
	logger.Info("Account:\t", juryInfo.Uname)
	logger.Info("CaseTotal:\t", juryInfo.CaseTotal)
	logger.Info("Status:\t", juryInfo.Status)
	split()
	if juryInfo.Status != 1 {
		logger.Fatal("you are not in jury.")
	}
	for {
		logger.Info("fetching next case...")
		split()
		caseID, err := getNext()
		if err != nil {
			logger.Fatal("failed to fetch next case:", err.Error())
		}
		caseInfo, err := getCaseInfo(caseID)
		if err != nil {
			logger.Fatal("failed to fetch case info:", err.Error())
		}
		voteItem := map[int]string{}
		for _, item := range caseInfo.VoteItems {
			voteItem[item.Vote] = item.VoteText
		}

		logger.Info("CaseID:\t", caseID)
		logger.Info("CaseType:\t", caseInfo.CaseType)
		split()
		logger.Info("fetching opinions...")
		split()
		opinion, err := getOpinion(caseID, 1, 20)
		if err != nil {
			logger.Fatal("failed to fetch opinions:", err.Error())
		}
		voteMap := map[int]int{}
		for _, s := range opinion.List {
			voteMap[s.Vote] = voteMap[s.Vote] + 1
		}
		maxKey := 0
		maxValue := 0
		for k, v := range voteMap {
			if v > maxValue {
				maxKey = k
				maxValue = v
			}
		}
		logger.Info("MaxKey:\t", maxKey)
		logger.Info("MaxValue:\t", maxValue)
		logger.Info("VoteText:\t", voteItem[maxKey])
		split()
		logger.Info("claiming case...")
		_, err = postVote(caseID, 0)
		if err != nil {
			logger.Fatal("failed to claim case:", err.Error())
		}
		logger.Info("claim succeeded. sleeping for 10s...")
		time.Sleep(10 * time.Second)
		logger.Info("voting ...")
		_, err = postVote(caseID, maxKey)
		if err != nil {
			logger.Fatal("failed to vote:", err.Error())
		}
		logger.Info("vote succeeded. sleeping for 10s to next round...")
		split()
		time.Sleep(10 * time.Second)
	}
}

func split() {
	logger.Info(strings.Repeat("-", 50))
}
