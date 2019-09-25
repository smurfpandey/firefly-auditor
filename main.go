package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/smurfpandey/firefly-auditor/firefly"
	"github.com/smurfpandey/firefly-auditor/kotak"
)

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
	transactionFile := flag.String("transactions", "", "Path to the file with list of transactions")
	accountName := flag.String("account", "", "Name of asset account in Firefly")
	flag.Parse()

	// Read file with help of bank manager
	kotak.ReadCSV(*transactionFile)

	yoAccount := firefly.GetAssetAccount(*accountName)
	if yoAccount == nil {
		fmt.Println("nahi mila")
	} else {
		fmt.Println(*yoAccount)
	}

}
