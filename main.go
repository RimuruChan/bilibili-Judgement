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
	logger.Info("Version	", Version, "	Running...")
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
	logger.Info("SESSDATA:	", sessdata)
	logger.Info("CRSF:	", csrf)
	split()
	juryInfo, err := getJuryInfo()
	if err != nil {
		logger.Fatal("failed to fetch jury info:", err.Error())
	}
	logger.Info("Account Info")
	split()
	logger.Info("Account:	", juryInfo.Uname)
	logger.Info("CaseTotal:	", juryInfo.CaseTotal)
	logger.Info("Status:	", juryInfo.Status)
	logger.Info("TermEnd:	", time.Unix(int64(juryInfo.TermEnd), 0).Format("2006-01-02 15:04:05"))
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

		logger.Info("CaseID:	", caseID)
		logger.Info("CaseType:	", caseInfo.CaseType)
		split()
		logger.Info("fetching opinions...")
		split()

		opinion := getFullOpinion(caseID)

		retryTimes := 0
		for len(opinion.List) < 10 {
			if retryTimes >= 5 {
				logger.Info("Retried 5 times! fetch another case...")
				caseID, err = getNext()
				if err != nil {
					logger.Fatal("failed to fetch next case:", err.Error())
				}
				split()
				logger.Info("CaseID:	", caseID)
				logger.Info("CaseType:	", caseInfo.CaseType)
				split()
				opinion = getFullOpinion(caseID)
				retryTimes = 0
			} else {
				logger.Info("Votes amount less than 10, retry fetching after 10 seconds...")
				time.Sleep(10 * time.Second)
				retryTimes++
				opinion = getFullOpinion(caseID)
			}
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
		logger.Info("VotedKey:		", maxKey)
		logger.Info("VotedTimes:		", maxValue)
		logger.Info("TotalVotedTimes:	", opinion.Total)
		logger.Info("VoteText:		", voteItem[maxKey])
		split()
		logger.Info("claiming case...")
		_, err = postVote(caseID, 0)
		if err != nil {
			logger.Fatal("failed to claim case:", err.Error())
		}
		logger.Info("claim succeeded. sleeping for 15s...")
		time.Sleep(15 * time.Second)
		logger.Info("voting ...")
		_, err = postVote(caseID, maxKey)
		if err != nil {
			logger.Fatal("failed to vote:", err.Error())
		}
		logger.Info("vote succeeded. sleeping for 15s to next round...")
		split()
		time.Sleep(15 * time.Second)
	}
}

func split() {
	logger.Info(strings.Repeat("-", 50))
}
