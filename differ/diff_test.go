package differ

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestChunkChange(t *testing.T){
	txt1:= "This is a Rolling hash file diff algorithm. It should check for changes in file and text"
	txt2:= "This is a Rolling hash file diff algorithm. It should check for changes in file and text. This is some extra text that should also be present"

	differInstance := New(16)
	reader1 := bytes.NewReader([]byte(txt1))
	buffReader1 := bufio.NewReader(reader1)

	reader2 := bytes.NewReader([]byte(txt2))
	buffReader2 := bufio.NewReader(reader2)

	signatures := differInstance.GenerateSignatures(buffReader1)
	deltas:=differInstance.GenerateDelta(signatures,buffReader2)
	fmt.Println(deltas,"Deltas")
} 