package kotak

import (
	"fmt"
	"log"
	"os"

	"github.com/gocarina/gocsv"
)

type Transaction struct {
	SlNo    int     `csv:"Sl. No."` // .csv column headers
	Date    string  `csv:"Date"`
	Amount  float32 `csv:"Amount"`
	Type    string  `csv:"Dr / Cr"`
	Balance float32 `csv:"Balance"`
}

func ReadCSV(filePath string) interface{} {
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

	for _, transaction := range transactions {
		if transaction.Type == "DR" {
			fmt.Println("Gaya: ", transaction.Amount, "On: ", transaction.Date)
		} else {
			fmt.Println("Aaya: ", transaction.Amount, "On: ", transaction.Date)
		}

	}

	return transactions
}
