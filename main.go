package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func exitWithMessage(message string) {
	fmt.Println("Error:", message)
	os.Exit(1)
}

func main() {
	// Init
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	FIREFLY_AUTH_TOKEN = os.Getenv("FIREFLY_ACCESS_TOKEN")
	FIREFLY_API_BASE_URL = os.Getenv("FIREFLY_API_BASE_URL")

	transactionFile := flag.String("transactions", "", "Path to the file with list of transactions")
	accountName := flag.String("account", "", "Name of asset account in Firefly")
	flag.Parse()

	fmt.Println(*transactionFile, *accountName)

	ReadCSV(*transactionFile)

	GetAssetAccount(*accountName)
}
