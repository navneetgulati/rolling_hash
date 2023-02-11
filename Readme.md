# File Diff using Rolling Hash Algorithm

This is an implementation of file diffing algorithm using adler32 rolling hash. 

It first generates a signatures file from the original file contents. This signature file contains array of adler32 hashes constructed from the byte chunks of the original file.
Then this signature file as a reference to identify any changes in the updated files.

## How to Run

`go run main.go CHUNK_SIZE PATH_OF_ORIGINAL_FILE PATH_OF_UPDATED_FILE`

where

`CHUNK_SIZE (Optional): Size of chunks that will be used for calculating signatures and rolling hash. Default value : 64 `

`PATH_OF_ORIGINAL_FILE (Optional): Path to the file which will be used as base for comparison. Signatures will be generated on this file. Default value: mock_data/testData.txt`

`PATH_OF_UPDATED_FILE (Optional): Path to the file against which the original file will be compared. Default value: mock_data/updatedTestData.txt`

> P.S: Make sure that the order of the arguments need to be maintained. Eg: If you want to add PATH_OF_ORIGINAL_FILE, you will need to add CHUNK_SIZE too because it is the first argument`

## Output

It returns a map of chunk index (starting from zero) to Delta struct defined below
```
type Delta struct{
	startIndex // The index in original file from where the replacement needs to begin

	endIndex // The index in the original file till where the replacement needs to happen

	deleted // A boolean variable it specifies if the chunk has been deleted

	updatedLiterals // string which will replace the existing ones between startIndex and endIndex

}
```
The output contains the list of chunks that have been deleted and the chunks that have been found but have some updated literals surronding them. Eg:

```
map[
	5:{641 768 true } 
	6:{641 768 false  the readable content of a page when looking at its layout. The point of using Ipsum Lorem is that it has a more-or-less normal } 
	9:{1153 1280 true } 
	10:{1153 1280 false pose
There are many variations of passages of Lorem Ipsum available, but the majority have suffe}
]
```
The output above can be read as:

- Chunk index 5 has not been found
- Chunk index 6 has been found but the content starting from index 641 and ending at 768 need to replaced with the text ` the readable content of a page when looking at its layout. The point of using Ipsum Lorem is that it has a more-or-less normal`
- Chunk index 9 has not been found
- Chunk index 10 has been found but the content starting from index 1153 to index 1280 needs to replaced with the text `pose
There are many variations of passages of Lorem Ipsum available, but the majority have suffe`

