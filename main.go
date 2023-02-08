package main

import (
	"fmt"
	"os"
	"path"

	"github.com/gulatinavneet/rollingHash/differ"
	"github.com/gulatinavneet/rollingHash/utils"
)

func main() {
	currDir,_ := os.Getwd()
	chunkSize := 4
	filepath := path.Join(currDir,"mock_data/testData.txt")
	if len(os.Args)>1{
		filepath = os.Args[1]
	}

	fileReader := utils.New(chunkSize)
	buffReader,err:=fileReader.Open(filepath)
	if err!=nil{
		panic(err)
	}

	differ := differ.New(chunkSize)
	signatures:=differ.GenerateSignatures(buffReader)

	fmt.Println(signatures)

}
//shj -> shT
//whd -> whD
//wuh -> wQh
