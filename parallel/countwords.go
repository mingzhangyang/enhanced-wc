package parallel

func countWords(buf []byte) (count int64) {
	var preIsSpace bool
	var i int
Loop:
	for ; i < len(buf); i++ {
		switch buf[i] {
		case ' ', '\n', '\t', '\r', '\v', '\f':
			preIsSpace = true
			break Loop
		}
	}
	if i != 0 {
		count++
	}
	i++
	for ; i < len(buf); i++ {
		switch buf[i] {
		case ' ', '\n', '\t', '\r', '\v', '\f':
			preIsSpace = true
		default:
			if preIsSpace {
				count++
				preIsSpace = false
			}
		}
	}
	return
}
