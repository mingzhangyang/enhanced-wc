package parallel

import (
	"bytes"
	"io"
	"os"
	"runtime"
)

// BufferSize is the default batch size for reading
const BufferSize = 16 * 1024

// NewLineChar is a new line char
var NewLineChar = []byte{'\n'}

// Remainder is the remainder content which is part of a whole line
type Remainder struct {
	Head  []byte
	Tail  []byte
	Range [2]int64
}

// Counter contains the counter information
type Counter struct {
	TotalBytes, TotalLines, TotalWords int64
}

// Result is the final statisics and remainder of a worker
type Result struct {
	Count  Counter
	Remain Remainder
}

// Summary is the final summary of the results of all workers
type Summary struct {
	FileName string
	Counter  Counter
}

func findFirstAndLastNewlineChar(buf []byte) (int, int) {
	first, last := 0, 0
	n := len(buf)
	for i := 0; i < n; i++ {
		if buf[i] == '\n' {
			first = i
			break
		}
	}
	for j := n - 1; j >= 0; j-- {
		if buf[j] == '\n' {
			last = j
			break
		}
	}
	return first, last
}

func parallelRead(fp string, start, end int64, res chan Result) {
	f, err := os.Open(fp)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := make([]byte, BufferSize, BufferSize)
	pos := start
	holder := make([]byte, 0, BufferSize)
	result := Result{
		Count:  Counter{},
		Remain: Remainder{Range: [2]int64{start, end}},
	}

	var n, first, last int

Loop:
	for pos < end {
		n, err = f.ReadAt(buf, pos)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}
		result.Count.TotalBytes += int64(n)
		result.Count.TotalLines += int64(bytes.Count(buf[:n], NewLineChar))

		first, last = findFirstAndLastNewlineChar(buf[:n])

		holder = append(holder, buf[:first]...)
		if pos == start {
			result.Remain.Head = make([]byte, len(holder))
			copy(result.Remain.Head, holder)
		} else {
			result.Count.TotalWords += countWords(holder)
		}
		result.Count.TotalWords += countWords(buf[first:last])

		holder = holder[:0]
		holder = append(holder, buf[last:n]...)

		pos += BufferSize
		if err == io.EOF {
			break Loop
		}
	}

	result.Remain.Tail = holder
	res <- result
}

// Wc is exported function of word counting
func Wc(fp string) (Summary, error) {
	info, err := os.Stat(fp)
	if err != nil {
		return Summary{}, err
	}

	// TB: total bytes of the file
	var TB = info.Size()
	var N, batchSize int64 = 1, TB

	if TB > BufferSize*BufferSize {
		N = int64(runtime.NumCPU())
		Bl := N * BufferSize
		rem := TB % Bl
		n := (TB + Bl - rem) / Bl
		batchSize = n * BufferSize
	}

	res := make(chan Result, N)

	var i int64
	for ; i < N; i++ {
		go parallelRead(fp, i*batchSize, (i+1)*batchSize, res)
	}

	var m = make(map[int64]*Remainder)
	var tb, tl, tw int64

	for i = 0; i < N; i++ {
		r := <-res
		tb += r.Count.TotalBytes
		tl += r.Count.TotalLines
		tw += r.Count.TotalWords

		if _, ok := m[r.Remain.Range[0]]; !ok {
			m[r.Remain.Range[0]] = &Remainder{}
		}
		if _, ok := m[r.Remain.Range[1]]; !ok {
			m[r.Remain.Range[1]] = &Remainder{}
		}
		m[r.Remain.Range[0]].Tail = r.Remain.Head
		m[r.Remain.Range[1]].Head = r.Remain.Tail
	}

	//println(len(m))
	for _, v := range m {
		s := make([]byte, 0, len(v.Head)+len(v.Tail))
		s = append(s, v.Head...)
		s = append(s, v.Tail...)
		tw += countWords(s)
	}

	return Summary{
		FileName: info.Name(),
		Counter: Counter{
			TotalWords: tw,
			TotalLines: tl,
			TotalBytes: tb,
		},
	}, nil
}
