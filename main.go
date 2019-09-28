package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/joho/godotenv"

	"github.com/smurfpandey/firefly-auditor/accounts"
	"github.com/smurfpandey/firefly-auditor/accounts/hdfc"
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

func ReadTransactions(filePath string, accountName string) []accounts.Transaction {
	switch accountName {
	case "Kotak Mahindra Bank":
		return kotak.ReadTransactions(filePath)
	case "HDFC Bank":
		return hdfc.ReadTransactions(filePath)
	default:
		return []accounts.Transaction{}
	}
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

	accountName := *ptrAccountName
	transactionFile := *ptrTransactionFile

	// Read file with help of bank manager
	bankTransactions := ReadTransactions(transactionFile, accountName)

	if len(bankTransactions) == 0 {
		exitWithMessage("No transactions found in CSV file")
	}

	sort.Slice(bankTransactions, func(i, j int) bool {
		return bankTransactions[i].Date.Before(bankTransactions[j].Date)
	})

	firstTransactionOn := bankTransactions[0].Date
	lastTransactionOn := bankTransactions[len(bankTransactions)-1].Date

	ptrAssetAccount := firefly.GetAssetAccount(accountName)
	if ptrAssetAccount == nil {
		exitWithMessage("No account found with that name. Are you sure you want to audit \"" + accountName + "\"?")
	}

	assetAccount := *ptrAssetAccount
	fireflyTransactions := firefly.GetAllTransactions(assetAccount.Id, firstTransactionOn, lastTransactionOn)
	sort.Slice(fireflyTransactions, func(i, j int) bool {
		return fireflyTransactions[i].Attributes.Transactions[0].Date.Before(fireflyTransactions[j].Attributes.Transactions[0].Date)
	})

	for bankTIdx := range bankTransactions {
		didILogIt := false
		bankTransDate := bankTransactions[bankTIdx].Date.Format("2016-01-02")
		for fireflyTIdx := range fireflyTransactions {
			fireflyTransDate := fireflyTransactions[fireflyTIdx].Attributes.Transactions[0].Date.Format("2016-01-02")
			if bankTransDate == fireflyTransDate {
				if bankTransactions[bankTIdx].Amount == fireflyTransactions[fireflyTIdx].Attributes.Transactions[0].Amount {
					didILogIt = true
					break
				}
			}
		}

		if didILogIt {
			//fmt.Println("Yes!")
		} else {
			fmt.Println("Transaction not logged in Firefly! :(", bankTransactions[bankTIdx])
		}
	}
}
