package main

import (
	"fmt"
	"io/ioutil"

	"github.com/imroc/req"
)

var (
	FIREFLY_AUTH_TOKEN   string
	FIREFLY_API_BASE_URL string
)

func GetAssetAccount(accountName string) {
	authHeader := req.Header{
		"Accept":        "application/json",
		"Authorization": "Bearer " + FIREFLY_AUTH_TOKEN,
	}
	param := req.Param{
		"type": "asset",
	}
	reqUrl := FIREFLY_API_BASE_URL + "accounts"

	rawResp, _ := req.Get(reqUrl, authHeader, param)
	resp := rawResp.Response()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	bodyString := string(bodyBytes)

	fmt.Println(bodyString)
}
