package differ

import (
	"bufio"
	"io"

	"github.com/gulatinavneet/rollingHash/hasher"
)

type Delta struct{
	startIndex int
	endIndex int
	deleted bool
	updatedLiterals string
}

type Differ struct{
	chunkSize int //ChunkSize in bytes
}

func New(chunkSize int) *Differ{
	return &Differ{ chunkSize: chunkSize}
}

func (diff *Differ)GenerateSignatures(reader *bufio.Reader) map[uint]int{
	var signatures = make(map[uint]int)
	var count=0
	for {
		chunk := make([]byte,diff.chunkSize)
		bytes,err:=reader.Read(chunk)
		if bytes ==0 || err == io.EOF{
			break
		}
		if bytes<diff.chunkSize{
			chunk=chunk[:bytes]
		}
		adler32:= hasher.NewAdler32(diff.chunkSize)
		adler32.Write(chunk)
		signature:= adler32.Hash()
		signatures[signature]=count
		count++
	}
	return signatures
}

func (diff *Differ) GenerateDelta(signatures map[uint]int,reader *bufio.Reader) map[int]Delta{
	deltas:= make(map[int]Delta)
	adler32 := hasher.NewAdler32(diff.chunkSize)

	var diffingLiterals []byte
	lastFoundIndex:=-1

	//Write logic for finding it in existing signature file
	for{
		c,err:= reader.ReadByte()
		if err == io.EOF || err!=nil {
			break
		}
		hash:=adler32.RollIn(c)

		if adler32.WindowLength()<diff.chunkSize{
			//Check if this is not the last byte before continuing
			if next,_:=reader.Peek(1);len(next)>0{
				continue
			}
		}
		index := findIndex(signatures,hash)

		if index != -1{
			adler32.Reset()
			deltas[index]=Delta{
				startIndex: index*diff.chunkSize+1,
				endIndex: (index+1)*diff.chunkSize,
				deleted: false,
				updatedLiterals: string(diffingLiterals),
			}
			diffingLiterals=diffingLiterals[:0]
			lastFoundIndex=index
			continue
		}
			_,removed:=adler32.RollOut()
			diffingLiterals = append(diffingLiterals, removed)
	}
	//If 
	if len(diffingLiterals)>0{
		diffingLiterals = append(diffingLiterals, adler32.GetWindowLiterals()...)
		deltas[lastFoundIndex+1]=Delta{
			startIndex: (lastFoundIndex+1)*diff.chunkSize,
			endIndex: (lastFoundIndex+1)*diff.chunkSize+len(diffingLiterals),
			deleted: false,
			updatedLiterals: string(diffingLiterals),
		}
	}
	deltas=diff.updateDeletedChunks(signatures,deltas)
	return deltas
}

func (diff *Differ)updateDeletedChunks(signatures map[uint]int, delta map[int]Delta)map[int]Delta{
	for _,index:=range signatures{
		if value,ok:=delta[index];!ok{
			delta[index]=Delta{
				startIndex: index*diff.chunkSize+1,
				endIndex: (index+1)*diff.chunkSize,
				deleted: true,
				updatedLiterals: "",
			}
		}else if len(value.updatedLiterals)==0{
			delete(delta,index)
		}
	}
	return delta
}

func findIndex(signatures map[uint]int,hash uint) int{
	if index,found:= signatures[hash];found{
		return index
	}
	return -1
}