package main

import (
	// "flag"
	"fmt"
	"log"
	"os"
	// "sort"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
	// "github.com/olekukonko/tablewriter"

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

	// 1. Ask what user wants to to today: Audit Transactions, Audit Final Balance

	// 1.1 If Transactions: Fetch list of all accounts from Firefly, then ask which account user wants to audit
	// 1.2 Ask source of transactions: Show last 5 files(?)
	// 1.3 Then audit

	w := wow.New(os.Stdout, spin.Get(spin.Dots), " Loading All Accounts!")
	w.Start()
	lstAccounts := firefly.GetAllAssetAccounts()
	fmt.Println("")
	w.Stop()

	templates := &promptui.SelectTemplates{
		// Label:    "{{ . }}",
		Active:   "▸ {{ .Attributes.Name | cyan }}",
		Inactive: "{{ .Attributes.Name | cyan }}",
		Selected: "✔ {{ .Attributes.Name | white }}",
	}

	prompt := promptui.Select{
		Label:     "Select Account",
		Items:     lstAccounts,
		Templates: templates,
		Size:      14,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	yoAccount := lstAccounts[i]

	promptAction := promptui.Select{
		Label: "What would you like to do with " + yoAccount.Attributes.Name + " today?",
		Items: []string{"Audit Transactions", "Audit Balance"},
	}

	idxSelect, _, err := promptAction.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch idxSelect {
	case 0:
		AuditTransactions(yoAccount)
	case 1:
		fmt.Println("Balance")
	default:
		fmt.Println("What?!?")
	}

	// for i := range lstAccounts {
	// 	// fmt.Println("Account: ", lstAccounts[i].Attributes.Name)
	// 	if lstAccounts[i].Attributes.Name == accountName {
	// 		// Found!
	// 		return &lstAccounts[i]
	// 	}
	// }

	// Read command line arguments
	// ptrTransactionFile := flag.String("transactions", "", "Path to the file with list of transactions")
	// ptrAccountName := flag.String("account", "", "Name of asset account in Firefly")
	// flag.Parse()

	// accountName := *ptrAccountName
	// transactionFile := *ptrTransactionFile

	// // Read file with help of bank manager
	// bankTransactions := ReadTransactions(transactionFile, accountName)

	// if len(bankTransactions) == 0 {
	// 	exitWithMessage("No transactions found in CSV file")
	// }

	// sort.Slice(bankTransactions, func(i, j int) bool {
	// 	return bankTransactions[i].Date.Before(bankTransactions[j].Date)
	// })

	// firstTransactionOn := bankTransactions[0].Date
	// lastTransactionOn := bankTransactions[len(bankTransactions)-1].Date

	// ptrAssetAccount := firefly.GetAssetAccount(accountName)
	// if ptrAssetAccount == nil {
	// 	exitWithMessage("No account found with that name. Are you sure you want to audit \"" + accountName + "\"?")
	// }

	// assetAccount := *ptrAssetAccount
	// fireflyTransactions := firefly.GetAllTransactions(assetAccount.Id, firstTransactionOn, lastTransactionOn)
	// sort.Slice(fireflyTransactions, func(i, j int) bool {
	// 	return fireflyTransactions[i].Attributes.Transactions[0].Date.Before(fireflyTransactions[j].Attributes.Transactions[0].Date)
	// })

	// table := tablewriter.NewWriter(os.Stdout)
	// table.SetHeader([]string{"Date", "Description", "Amount", "CR/DR"})

	// for bankTIdx := range bankTransactions {
	// 	didILogIt := false
	// 	bankTransDate := bankTransactions[bankTIdx].Date.Format("2016-01-02")
	// 	for fireflyTIdx := range fireflyTransactions {
	// 		fireflyTransDate := fireflyTransactions[fireflyTIdx].Attributes.Transactions[0].Date.Format("2016-01-02")
	// 		if bankTransDate == fireflyTransDate {
	// 			if bankTransactions[bankTIdx].Amount == fireflyTransactions[fireflyTIdx].Attributes.Transactions[0].Amount {
	// 				didILogIt = true
	// 				break
	// 			}
	// 		}
	// 	}

	// 	if didILogIt {
	// 		//fmt.Println("Yes!")
	// 	} else {
	// 		yoRow := []string{bankTransDate, bankTransactions[bankTIdx].Description, fmt.Sprintf("%f", bankTransactions[bankTIdx].Amount), bankTransactions[bankTIdx].Type}
	// 		table.Append(yoRow)
	// 	}
	// }

	// table.Render()
}
