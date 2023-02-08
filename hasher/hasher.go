package hasher

import (
	goAdler "hash/adler32"
	"math"

	"github.com/gulatinavneet/rollingHash/utils"
)

const ADLER_CONSTANT=65521

type adler32 struct{
	s1 uint
	s2 uint
	window []byte
}

func NewAdler32(chunkSize int) *adler32{
	return &adler32{
		s1: 1,
		s2: 0,
		window: make([]byte,0, chunkSize),
	}
}

func (ad *adler32) WindowLength() int{
	return len(ad.window)
}

func (ad *adler32) Sum(chunk []byte){
	ad.window = append(ad.window, chunk...)
}

func (ad *adler32)Write(chunk []byte){
	ad.window = chunk
}

func (ad *adler32) Hash() uint{
	for _, val := range ad.window{
		ad.s1+=uint(val)
		ad.s2+=ad.s1
	}
	ad.s1=ad.s1%ADLER_CONSTANT
	ad.s2=ad.s2%ADLER_CONSTANT
	return ad.s2 << 16 + ad.s1
}

func (ad *adler32) RollIn(c byte) (uint){
	ad.window = append(ad.window, c)
	ad.s1 = (ad.s1 + uint(c)) % ADLER_CONSTANT
	ad.s2 = (ad.s2 + ad.s1 ) % ADLER_CONSTANT
	return ad.s2 << 16 + ad.s1
}

func (ad *adler32) RollOut() (uint,byte){
	removed := ad.window[0]
	ad.s1 = (ad.s1 - uint(removed))%ADLER_CONSTANT
	ad.s2 = (ad.s2 - (uint(len(ad.window))*uint(removed))-1)%ADLER_CONSTANT
	ad.window = ad.window[1:]
	return ad.s2 << 16 + ad.s1,removed
}

func (ad *adler32) GetWindowLiterals()[]byte{
	return ad.window
}

func (ad *adler32) Reset(){
	utils.Clear(ad)
	ad.s1=1
}

// func (roll *RollingHash) RollingHash(data []byte) map[int]uint{
// 	var signatures map[int]uint
// 	hash := polynomialHash(data[:roll.WindowSize])
// 	signatures[0]=hash
// 	for i :=roll.WindowSize;i<len(data);i++{
// 		newHash := hash + uint(data[i])*uint(math.Pow(5,float64(roll.WindowSize-1)))
// 		signatures[i-roll.WindowSize+1]=newHash
// 	}
// 	return signatures
// }

func polynomialHash(chunk []byte) uint{
	var hash uint
	for index,val:=range chunk{
		hash += uint(val)*uint(math.Pow(5,float64(index)))
	}
	return hash
}

func GenerateAdlerHash(chunk []byte) uint32{
	return goAdler.Checksum(chunk)
}
