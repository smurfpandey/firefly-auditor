package main

import (
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
	"github.com/smurfpandey/firefly-auditor/accounts"
	"github.com/smurfpandey/firefly-auditor/accounts/hdfc"
	"github.com/smurfpandey/firefly-auditor/accounts/kotak"
	"github.com/smurfpandey/firefly-auditor/accounts/paytm"
	"github.com/smurfpandey/firefly-auditor/accounts/sbi_cc"
	"github.com/smurfpandey/firefly-auditor/firefly"
	"github.com/smurfpandey/firefly-auditor/utils"
)

func ReadTransactions(filePath string, accountName string) []accounts.Transaction {
	switch accountName {
	case "Kotak Mahindra Bank":
		return kotak.ReadTransactions(filePath)
	case "HDFC Bank":
		return hdfc.ReadTransactions(filePath)
	case "Paytm Wallet":
		return paytm.ReadTransactions(filePath)
	case "SBI CC":
		return sbi_cc.ReadTransactions(filePath)
	default:
		return []accounts.Transaction{}
	}
}

func ListFiles(accountName string) []utils.TransactionFile {
	switch accountName {
	case "Kotak Mahindra Bank":
		return utils.ListFiles(kotak.BASE_FOLDER_PATH())
	case "HDFC Bank":
		return utils.ListFiles(hdfc.BASE_FOLDER_PATH())
	case "Paytm Wallet":
		return utils.ListFiles(paytm.BASE_FOLDER_PATH())
	case "SBI CC":
		return utils.ListFiles(sbi_cc.BASE_FOLDER_PATH())
	default:
		return nil
	}
}

func ShowFiles(files []utils.TransactionFile) *utils.TransactionFile {
	templates := &promptui.SelectTemplates{
		// Label:    "{{ . }}",
		Active:   "▸ {{ .Name | cyan }}",
		Inactive: "{{ .Name | cyan }}",
		Selected: "✔ {{ .Name | white }}",
	}

	prompt := promptui.Select{
		Label:     "Select File",
		Items:     files,
		Templates: templates,
		Size:      10,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil
	}

	return &files[i]
}

func AuditTransactions(account firefly.Account) {
	accountName := account.Attributes.Name
	files := ListFiles(accountName)

	fmt.Printf("Found %d files\n", len(files))
	ptrSelectedFile := ShowFiles(files)
	if ptrSelectedFile == nil {
		exitWithMessage("Something went wrong")
	}
	selectedFile := *ptrSelectedFile
	fmt.Printf("Auditing %v with %v\n", accountName, selectedFile.Name)

	bankTransactions := ReadTransactions(selectedFile.Path, accountName)
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Description", "Amount", "CR/DR"})

	for bankTIdx := range bankTransactions {
		didILogIt := false
		bankTransDate := bankTransactions[bankTIdx].Date.Format("02-01-2006")
		for fireflyTIdx := range fireflyTransactions {
			fireflyTransDate := fireflyTransactions[fireflyTIdx].Attributes.Transactions[0].Date.Format("02-01-2006")
			if bankTransDate == fireflyTransDate {
				bankAmount := math.Round(float64(bankTransactions[bankTIdx].Amount))
				fireflyAmount := math.Round(float64(fireflyTransactions[fireflyTIdx].Attributes.Transactions[0].Amount))
				if bankAmount == fireflyAmount {
					didILogIt = true
					break
				}
			}
		}

		if didILogIt {
			//fmt.Println("Yes!")
		} else {
			yoRow := []string{bankTransDate, bankTransactions[bankTIdx].Description, fmt.Sprintf("%.2f", bankTransactions[bankTIdx].Amount), bankTransactions[bankTIdx].Type}
			table.Append(yoRow)
		}
	}

	table.Render()
}
