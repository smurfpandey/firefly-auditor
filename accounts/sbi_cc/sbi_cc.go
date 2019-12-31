package sbi_cc

import (
	"fmt"
	"time"
	"os"
	"encoding/json"
	"log"
	"os/exec"
    "bytes"

	"github.com/smurfpandey/firefly-auditor/accounts"
)

type Transaction struct {
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Amount      float32 `json:"amount"`
	Type        string  `json:"type"`
}

const DATE_FORMAT = "02/01/2006"
func BASE_FOLDER_PATH() string {
	return os.Getenv("SBI_CC_FOLDER_BASE_PATH")
}


func ReadTransactions(filePath string) []accounts.Transaction {
	fileTransactions := []*Transaction{}

	// execute nodejs script
	cmd := "node"
    args := []string{"parse_sbicc_pdf.js", filePath}
	process := exec.Command(cmd, args...)
	stdin, err := process.StdinPipe()
    if err != nil {
        fmt.Println(err)
    }
    defer stdin.Close()
    buf := new(bytes.Buffer) // THIS STORES THE NODEJS OUTPUT
    process.Stdout = buf
    process.Stderr = os.Stderr

    if err = process.Start(); err != nil {
		fmt.Println("An error occured: ", err)
		log.Fatal("Error reading .pdf file")
    }

    process.Wait()
    err = json.Unmarshal(buf.Bytes(), &fileTransactions)
	if err != nil {
		fmt.Println("error:", err)
		log.Fatal("Error parsing pdf output")
	}

	var outTransactions []accounts.Transaction
	for _, transaction := range fileTransactions {
		transDate, _ := time.Parse(DATE_FORMAT, transaction.Date)
		outTransaction := accounts.Transaction{
			Date: transDate,
			Amount: transaction.Amount,
			Type: transaction.Type,
			Description: transaction.Description,
		}

		outTransactions = append(outTransactions, outTransaction)
	}

	return outTransactions
}

