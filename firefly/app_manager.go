package firefly

import (
	"time"

	"github.com/imroc/req"
)

type Account struct {
	Id         string `json:"id"`
	Attributes struct {
		Name           string  `json:"name"`
		Type           string  `json:"type"`
		IsActive       bool    `json:"active"`
		CurrentBalance float32 `json:"current_balance"`
	} `json:"attributes"`
}

type Transaction struct {
	Type             string `json:"type"`
	Date 		     time.Time `json:"date"`
	Currency 	     string `json:"currency_code"`
	Amount 		     float32 `json:"amount,string"`
	SourceType       string `json:"source_type"`
	DestinationType  string `json:"destination_type"`
}

type ParentTransaction struct {
	Id         string `json:"id"`
	Attributes struct {
		CreatedOn    string  `json:"name"`
		UpdateOn     string  `json:"type"`
		Transactions []Transaction `json:"transactions"`
	} `json:"attributes"`
}

type ListAccount struct {
	Accounts []Account `json:"data"`
}
type ListTransaction struct {
	ParentTransactions []ParentTransaction `json:"data"`
	Meta Meta `json:"meta"`
}

type Meta struct {
	Pagination struct {
		CurrentPage int `json:"current_page"`
		TotalPages int `json:"total_pages"`
		PerPage int `json:"per_page"`
	} `json:"pagination"`
}

var (
	ACCESS_TOKEN   string
	API_BASE_URL string
)

func FetchAccountList() []Account {
	authHeader := req.Header{
		"Accept":        "application/json",
		"Authorization": "Bearer " + ACCESS_TOKEN,
	}
	param := req.Param{
		"type": "asset",
	}
	reqUrl := API_BASE_URL + "accounts"

	rawResp, _ := req.Get(reqUrl, authHeader, param)

	var lstAccounts ListAccount
	rawResp.ToJSON(&lstAccounts)

	return lstAccounts.Accounts
}

func FetchTransactions(accountId string, pageNumber int, start string, end string) ListTransaction {
	authHeader := req.Header{
		"Accept":        "application/json",
		"Authorization": "Bearer " + ACCESS_TOKEN,
	}
	param := req.QueryParam{
		"limit": 20,
		"page": pageNumber,
		"start": start,
		"end": end,
	}

	reqUrl := "transactions"
	if accountId != "0" {
		reqUrl = API_BASE_URL + "accounts/" + accountId + "/transactions"
	}

	rawResp, _ := req.Get(reqUrl, authHeader, param)

	var lstTransactions ListTransaction
	rawResp.ToJSON(&lstTransactions)

	return lstTransactions
}

func GetAssetAccount(accountName string) *Account {
	lstAccounts := FetchAccountList()

	for i := range lstAccounts {
		// fmt.Println("Account: ", lstAccounts[i].Attributes.Name)
		if lstAccounts[i].Attributes.Name == accountName {
			// Found!
			return &lstAccounts[i]
		}
	}

	return nil
}

func GetAllTransactions(accountId string, startDate time.Time, endDate time.Time) []ParentTransaction {
	var lstTransactions []ParentTransaction

	start := startDate.Format("2006-01-02")
	end := endDate.Format("2006-01-02")
	pageNum := 1
	for {
		pageTransactions := FetchTransactions(accountId, pageNum, start, end)

		pageData := pageTransactions.Meta.Pagination
		lstTransactions = append(lstTransactions, pageTransactions.ParentTransactions...)
		pageNum++
		if pageData.CurrentPage == pageData.TotalPages {
			break
		}
	}

	return lstTransactions
}
