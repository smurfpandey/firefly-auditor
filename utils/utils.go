package utils

import (
	"time"
	"fmt"
	"sort"
	"io/ioutil"
)

type TransactionFile struct {
	Name             string
	Path             string
	LastModifiedTime time.Time // modification time
}

func ListFiles(BASE_FOLDER_PATH string) []TransactionFile {
	fmt.Println(BASE_FOLDER_PATH)
	files, err := ioutil.ReadDir(BASE_FOLDER_PATH)

	if err != nil {
		return []TransactionFile{}
	}

	// TODO: handle the error!
	sort.Slice(files, func(i,j int) bool{
		return files[i].ModTime().After(files[j].ModTime())
	})

	var lstTransactions []TransactionFile

	for _, file := range files {
		transFile := TransactionFile{
			Name:             file.Name(),
			Path:             BASE_FOLDER_PATH + file.Name(),
			LastModifiedTime: file.ModTime(),
		}

		lstTransactions = append(lstTransactions, transFile)
	}

	return lstTransactions
}