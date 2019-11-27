package parallel

func ReadLine(c <-chan []byte) Foo {
	res := Foo{}
	for line := range c {
		res.Counter.TotalLines += 1
		res.Counter.TotalBytes += int64(len(line) + 1)
		res.Counter.TotalWords += CountWords(line)
	}
	return res
}
