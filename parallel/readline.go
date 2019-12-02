package parallel

// ReadLine read one line and count lines, bytes and words
func ReadLine(c <-chan []byte) Summary {
	res := Summary{}
	for line := range c {
		res.Counter.TotalLines++
		res.Counter.TotalBytes += int64(len(line) + 1)
		res.Counter.TotalWords += countWords(line)
	}
	return res
}
