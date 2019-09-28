package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/smurfpandey/firefly-auditor/accounts/kotak"
	"github.com/smurfpandey/firefly-auditor/firefly"
)

type Transaction struct {
	Date    string
	Amount  float32
	Type    string
	Balance float32
}

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// set variables from Environment Variables
	firefly.ACCESS_TOKEN = os.Getenv("FIREFLY_ACCESS_TOKEN")
	firefly.API_BASE_URL = os.Getenv("FIREFLY_API_BASE_URL")

	// Read command line arguments
	ptrTransactionFile := flag.String("transactions", "", "Path to the file with list of transactions")
	ptrAccountName := flag.String("account", "", "Name of asset account in Firefly")
	flag.Parse()

	// Read file with help of bank manager
	transactionFile := *ptrTransactionFile
	transactions := kotak.ReadTransactions(transactionFile)

	accountName := *ptrAccountName
	ptrAssetAccount := firefly.GetAssetAccount(accountName)
	if ptrAssetAccount == nil {
		exitWithMessage("No account found with that name. Are you sure you want to audit \"" + accountName + "\"?")
	}

	assetAccount := *ptrAssetAccount
	fmt.Println("================================================")
	fmt.Println(assetAccount)
	fmt.Println("================================================")
	fmt.Println(transactions)
	fmt.Println("================================================")

	fireflyTransactions := firefly.FetchTransactions(assetAccount.Id)
	fmt.Println(fireflyTransactions)
}
