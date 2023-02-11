package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/gulatinavneet/rollingHash/differ"
	"github.com/gulatinavneet/rollingHash/utils"
)

func main() {
	currDir, _ := os.Getwd()
	chunkSize := 16
	if len(os.Args) > 1 {
		newChunkSize, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println(fmt.Errorf("expected integer for chunk size. Going with default value"))
		}
		chunkSize = newChunkSize
	}
	filepath := path.Join(currDir, "mock_data/testData.txt")
	if len(os.Args) > 2 {
		filepath = os.Args[2]
	}

	fileReader := utils.New(chunkSize)
	buffReader, err := fileReader.Open(filepath)
	if err != nil {
		panic(err)
	}

	updatedFilePath := path.Join(currDir, "mock_data/updatedTestData.txt")
	if len(os.Args) > 3 {
		updatedFilePath = os.Args[3]
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
