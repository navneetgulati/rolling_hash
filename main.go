package main

import (
	"fmt"
	"os"
	"path"

	"github.com/gulatinavneet/rollingHash/differ"
	"github.com/gulatinavneet/rollingHash/utils"
)

func main() {
	currDir, _ := os.Getwd()
	chunkSize := 64
	filepath := path.Join(currDir, "mock_data/testData.txt")
	if len(os.Args) > 1 {
		filepath = os.Args[1]
	}

	fileReader := utils.New(chunkSize)
	buffReader, err := fileReader.Open(filepath)
	if err != nil {
		panic(err)
	}

	updatedFilePath := path.Join(currDir, "mock_data/updatedTestData.txt")
	if len(os.Args) > 2 {
		updatedFilePath = os.Args[2]
	}
	updatedFileReader := utils.New(chunkSize)
	updatedBuffReader, err := updatedFileReader.Open(updatedFilePath)
	if err != nil {
		panic(err)
	}

	differIns := differ.New(chunkSize)
	signatures := differIns.GenerateSignatures(buffReader)

	deltas := differIns.GenerateDelta(signatures, updatedBuffReader)

	fmt.Println(differ.PrettifyDelta(deltas))
}
