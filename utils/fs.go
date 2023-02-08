package utils

import (
	"bufio"
	"fmt"
	"os"
)

type FileReader struct{
	chunkSize int
}

func New(chunkSize int) FileReader{
	return FileReader{
		chunkSize: chunkSize,
	}
}

func (f FileReader) Open(filepath string) (*bufio.Reader,error){
	file,err:=os.Open(filepath)
	if err!=nil{
		return nil, fmt.Errorf("Error in reading file. Error Details %v",err)
	}

	fileInfo,_ := file.Stat()
	size := fileInfo.Size()
	
	chunks:= size/int64(f.chunkSize)
	if chunks <2{
		return nil, fmt.Errorf("File is not of appropriate size of generate enough chunks")
	}
	return bufio.NewReader(file),nil
}