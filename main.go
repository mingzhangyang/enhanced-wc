package main

import (
	"bufio"
	"github.com/mingzhangyang/fast-wc/parallel"
	"fmt"
	"log"
	"os"
)

func main() {
	var res parallel.Foo
	var err error
	switch len(os.Args) {
	case 1:
		res = parallel.ReadLine(CollectInput())
		fmt.Println("\t",
			res.Counter.TotalLines,
			"\t", res.Counter.TotalWords,
			"\t", res.Counter.TotalBytes)
	case 2:
		switch os.Args[1] {
		case "-l", "--lines":
			res = parallel.ReadLine(CollectInput())
			fmt.Println(res.Counter.TotalLines)
		case "-c", "--bytes":
			res = parallel.ReadLine(CollectInput())
			fmt.Println(res.Counter.TotalBytes)
		case "-w", "--words":
			res = parallel.ReadLine(CollectInput())
			fmt.Println(res.Counter.TotalWords)
		case "-h", "--help":
			printHelpInfo()
		default:
			res, err = parallel.Wc(os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("\t", res.Counter.TotalLines,
				"\t", res.Counter.TotalWords,
				"\t", res.Counter.TotalBytes,
				"\t", res.FileName)
		}
	case 3:
		var opt int
		switch os.Args[1] {
		case "-l", "--lines":
			opt = 1
		case "-c", "--bytes":
			opt = 3
		case "-w", "--words":
			opt = 2
		default:
			printHelpInfo()
			return
		}
		res, err = parallel.Wc(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		switch opt {
		case 1:
			fmt.Println(res.Counter.TotalLines)
		case 2:
			fmt.Println(res.Counter.TotalWords)
		case 3:
			fmt.Println(res.Counter.TotalBytes)
		}
	default:
		printHelpInfo()
	}

}

func printHelpInfo() {
	fmt.Println("fast-wc: an alternative to wc utility")
	fmt.Println("")
	fmt.Println("Synopsis: ")
	fmt.Println("\tfast-wc [OPTION] FILE")
	fmt.Println("\t... | fast-wc [OPTION]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("\t -c, --bytes,\tprint the byte count")
	fmt.Println("\t -w, --words,\tprint the word counts")
	fmt.Println("\t -l, --lines,\tprint the line counts")
	fmt.Println("\t -h, --help, \tprint the help information")
}

func CollectInput() <-chan []byte {
	var input = make(chan []byte, 16)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input <- scanner.Bytes()
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		close(input)
	}()
	return input
}
