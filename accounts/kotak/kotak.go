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
		transDate, _ := time.Parse("02/01/2006", transaction.Date)
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
