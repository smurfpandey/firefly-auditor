package firefly

import (
	// "fmt"
	// "io/ioutil"

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

type ListAccount struct {
	Accounts []Account `json:"data"`
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
	// resp := rawResp.Response()

	// bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// bodyString := string(bodyBytes)

	var lstAccounts ListAccount
	rawResp.ToJSON(&lstAccounts)

	return lstAccounts.Accounts
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
