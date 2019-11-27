package main

import (
	"os"
	"time"
	"wc/concurrent"
)

func main() {
	if len(os.Args) < 2 {
		panic("no file path specified")
	}
	filePath := os.Args[1]

	start := time.Now()
	concurrent.Wc(filePath)
	//parallel.Wc(filePath)
	end := time.Now()
	println(end.Sub(start).Seconds())

	println("=========================")

	//start = time.Now()
	//concurrent.Wc(filePath)
	//end = time.Now()
	//println(end.Sub(start).Seconds())
}


