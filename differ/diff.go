package differ

import (
	"bufio"
	"io"

	"github.com/gulatinavneet/rollingHash/hasher"
	"github.com/gulatinavneet/rollingHash/utils"
)

type Delta struct {
	startIndex      int
	endIndex        int
	deleted         bool
	updatedLiterals []byte
}
type PrettyDelta struct {
	startIndex      int
	endIndex        int
	deleted         bool
	updatedLiterals string
}
type Differ struct {
	chunkSize int //ChunkSize in bytes
}

func New(chunkSize int) *Differ {
	return &Differ{chunkSize: chunkSize}
}

func (d *Differ) GenerateSignatures(reader *bufio.Reader) map[uint]int {
	//Map of hash to the index of hash
	signatures := make(map[uint]int)
	count := 0
	for {
		chunk := make([]byte, d.chunkSize)
		bytes, err := reader.Read(chunk)
		if bytes == 0 || err == io.EOF {
			break
		}
		/*
			If the chunk cannot be filled with the remaining bytes in the file,
			only pick the slice that is filled
		*/
		if bytes < d.chunkSize {
			chunk = chunk[:bytes]
		}
		adler32 := hasher.NewAdler32(d.chunkSize)
		adler32.Write(chunk)
		signature := adler32.Hash()
		signatures[signature] = count
		count++
	}
	return signatures
}

func (d *Differ) GenerateDelta(signatures map[uint]int, reader *bufio.Reader) map[int]Delta {
	// Map of index to Delta struct
	deltas := make(map[int]Delta)
	adler32 := hasher.NewAdler32(d.chunkSize)

	var diffingLiterals []byte
	lastFoundIndex := -1

	for {
		c, err := reader.ReadByte()
		if err == io.EOF || err != nil {
			break
		}
		hash := adler32.RollIn(c)
		if adler32.WindowLength() < d.chunkSize {
			//Check if this is not the last byte before continuing
			if next, _ := reader.Peek(1); len(next) > 0 {
				continue
			}
		}
		index := findIndex(signatures, hash)

		if index != -1 {
			adler32.Reset()
			deltas[index] = Delta{
				startIndex:      (index-1)*d.chunkSize + 1,
				endIndex:        (index) * d.chunkSize,
				deleted:         false,
				updatedLiterals: diffingLiterals,
			}
			utils.Clear(&diffingLiterals)
			lastFoundIndex = index
			continue
		}
		_, removed := adler32.RollOut()
		diffingLiterals = append(diffingLiterals, removed)
	}
	if len(diffingLiterals) > 0 {
		diffingLiterals = append(diffingLiterals, adler32.GetWindowLiterals()...)
		deltas[lastFoundIndex+1] = Delta{
			startIndex:      (lastFoundIndex + 1) * d.chunkSize,
			endIndex:        (lastFoundIndex + 2) * d.chunkSize,
			deleted:         false,
			updatedLiterals: diffingLiterals,
		}
	}
	//Update the indexes of chunks that are not found
	deltas = d.updateDeletedChunks(signatures, deltas)
	return deltas
}

func (d *Differ) updateDeletedChunks(signatures map[uint]int, delta map[int]Delta) map[int]Delta {
	for _, index := range signatures {
		if value, ok := delta[index]; !ok {
			delta[index] = Delta{
				startIndex:      (index)*d.chunkSize + 1,
				endIndex:        (index + 1) * d.chunkSize,
				deleted:         true,
				updatedLiterals: []byte{},
			}
		} else if len(value.updatedLiterals) == 0 {
			/*
				If chunk is found and has no changes surronding it, then delete that entry from the map
				since we only need to present the diffing part
			*/
			delete(delta, index)
		}
	}
	return delta
}

func findIndex(signatures map[uint]int, hash uint) int {
	if index, found := signatures[hash]; found {
		return index
	}
	return -1
}

// convert byte array to string so that it is easy to read
func PrettifyDelta(deltas map[int]Delta) map[int]PrettyDelta {
	prettyDelta := make(map[int]PrettyDelta)
	for key, value := range deltas {
		prettyDelta[key] = PrettyDelta{
			startIndex:      value.startIndex,
			endIndex:        value.endIndex,
			deleted:         value.deleted,
			updatedLiterals: string(value.updatedLiterals),
		}
	}
	return prettyDelta
}
