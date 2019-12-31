package kotak

import (
	"fmt"
	"log"
	"time"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/smurfpandey/firefly-auditor/accounts"
)

type Transaction struct {
	SlNo        int     `csv:"Sl. No."` // .csv column headers
	Date        string  `csv:"Date"`
	Description string `csv:"Description"`
	Amount      float32 `csv:"Amount"`
	Type        string  `csv:"Dr / Cr"`
	Balance     float32 `csv:"Balance"`
}

const DATE_FORMAT = "02/01/2006"
func BASE_FOLDER_PATH() string {
	return os.Getenv("KOTAK_FOLDER_BASE_PATH")
}


func ReadTransactions(filePath string) []accounts.Transaction {
	fileTransactions := []*Transaction{}

	in, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error loading .csv file")
	}

	defer in.Close()

	if err := gocsv.UnmarshalFile(in, &fileTransactions); err != nil {
		fmt.Println(err)
		log.Fatal("Error parsing csv to struct")
	}

	var outTransactions []accounts.Transaction
	for _, transaction := range fileTransactions {
		transDate, _ := time.Parse(DATE_FORMAT, transaction.Date)
		outTransaction := accounts.Transaction{
			Date: transDate,
			Amount: transaction.Amount,
			Type: transaction.Type,
			Balance: transaction.Balance,
			Description: transaction.Description,
		}

		outTransactions = append(outTransactions, outTransaction)
	}

	return outTransactions
}