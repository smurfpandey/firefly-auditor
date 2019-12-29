package hdfc

import (
	"fmt"
	"log"
	"time"
	"strings"
	"sort"
	"os"
	"io/ioutil"

	"github.com/gocarina/gocsv"
	"github.com/smurfpandey/firefly-auditor/accounts"
	"github.com/smurfpandey/firefly-auditor/utils"
)

type Transaction struct {
	Date         string  `csv:"Date"`
	Description  string  `csv:"Narration"`
	CreditAmount float32 `csv:"Credit Amount"`
	DebitAmount  float32 `csv:"Debit Amount"`
	Balance      float32 `csv:"Closing Balance"`
}

func ReadTransactions(filePath string) []accounts.Transaction {
	transactions := []*Transaction{}

	in, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error loading .csv file")
	}

	defer in.Close()

	if err := gocsv.UnmarshalFile(in, &transactions); err != nil {
		fmt.Println(err)
		log.Fatal("Error parsing csv to struct")
	}

	var outTransactions []accounts.Transaction
	for _, transaction := range transactions {
		transDate, err := time.Parse("02/01/06", strings.TrimSpace(transaction.Date))
		if err != nil {
			log.Fatal("Error parsing date ", err)
		}

		var transAmount float32
		transType := ""
		if transaction.DebitAmount == 0 {
			transAmount = transaction.CreditAmount
			transType = "CR"
		} else {
			transAmount = transaction.DebitAmount
			transType = "DR"
		}
		outTransaction := accounts.Transaction{
			Date: transDate,
			Amount: transAmount,
			Type: transType,
			Balance: transaction.Balance,
			Description: strings.TrimSpace(transaction.Description),
		}

		outTransactions = append(outTransactions, outTransaction)
	}

	return outTransactions
}

func ListFiles() []utils.TransactionFile {
	BASE_FOLDER_PATH := os.Getenv("HDFC_FOLDER_BASE_PATH")

	files, err := ioutil.ReadDir(BASE_FOLDER_PATH)

	if err != nil {
		return []utils.TransactionFile{}
	}

	// TODO: handle the error!
	sort.Slice(files, func(i,j int) bool{
		return files[i].ModTime().After(files[j].ModTime())
	})

	var lstTransactions []utils.TransactionFile

	for _, file := range files {
		transFile := utils.TransactionFile{
			Name:             file.Name(),
			Path:             BASE_FOLDER_PATH + file.Name(),
			LastModifiedTime: file.ModTime(),
		}

		lstTransactions = append(lstTransactions, transFile)
	}

	return lstTransactions
}
