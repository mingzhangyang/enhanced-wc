package concurrent

import (
	"io"
	"os"
	"runtime"
	"sync"
)

// The codes below are from https://ajeetdsouza.github.io/blog/posts/beating-c-with-70-lines-of-go/
// it can be found at https://github.com/ajeetdsouza/blog-wc-go/blob/master/wc-mutex/main.go
// It is used as a comparison to the real parallel version.

func Wc(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileReader := &FileReader{
		File:            file,
		LastCharIsSpace: true,
	}
	counts := make(chan Count)

	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		go FileReaderCounter(fileReader, counts)
	}

	totalCount := Count{}
	for i := 0; i < numWorkers; i++ {
		count := <-counts
		totalCount.LineCount += count.LineCount
		totalCount.WordCount += count.WordCount
	}
	close(counts)

	fileStat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	byteCount := fileStat.Size()

	println(totalCount.LineCount, totalCount.WordCount, byteCount, file.Name())
}

type FileReader struct {
	File            *os.File
	LastCharIsSpace bool
	mutex           sync.Mutex
}

func (fileReader *FileReader) ReadChunk(buffer []byte) (Chunk, error) {
	fileReader.mutex.Lock()
	defer fileReader.mutex.Unlock()

	bytes, err := fileReader.File.Read(buffer)
	if err != nil {
		return Chunk{}, err
	}

	chunk := Chunk{fileReader.LastCharIsSpace, buffer[:bytes]}
	fileReader.LastCharIsSpace = IsSpace(buffer[bytes-1])

	return chunk, nil
}

func FileReaderCounter(fileReader *FileReader, counts chan Count) {
	const bufferSize = 16 * 1024
	buffer := make([]byte, bufferSize)

	totalCount := Count{}

	for {
		chunk, err := fileReader.ReadChunk(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		count := GetCount(chunk)
		totalCount.LineCount += count.LineCount
		totalCount.WordCount += count.WordCount
	}

	counts <- totalCount
}

type Chunk struct {
	PrevCharIsSpace bool
	Buffer          []byte
}

type Count struct {
	LineCount int
	WordCount int
}

func GetCount(chunk Chunk) Count {
	count := Count{}

	prevCharIsSpace := chunk.PrevCharIsSpace
	for _, b := range chunk.Buffer {
		switch b {
		case '\n':
			count.LineCount++
			prevCharIsSpace = true
		case ' ', '\t', '\r', '\v', '\f':
			prevCharIsSpace = true
		default:
			if prevCharIsSpace {
				prevCharIsSpace = false
				count.WordCount++
			}
		}
	}

	return count
}

func IsSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\v' || b == '\f'
}
