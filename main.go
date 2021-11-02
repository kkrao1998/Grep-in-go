package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		fmt.Println(f.Name)
		if f.Name == name {
			found = true
		}
	})
	return found
}

var root, query string
var found = 1
var wg sync.WaitGroup
var ans string

func readFile(path string) {

	file, err := os.Open(path)

	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		if strings.Contains(scanner.Text(), query) {
			found = 0
			ans += root + "/" + path + ":" + strconv.Itoa(i) + " : " + scanner.Text() + "\n"
			//fmt.Printf("%s/%s:%d: %s\n", root, path, i, scanner.Text())
		}
	}
}

func readDir(wg *sync.WaitGroup, path string) {
	defer wg.Done()

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return
	}
	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		if strings.Contains(scanner.Text(), query) {
			found = 0
			ans += root + "/" + path + ":" + strconv.Itoa(i) + " : " + scanner.Text() + "\n"
			//fmt.Printf("%s/%s:%d: %s\n", root, path, i, scanner.Text())
		}
	}
}

func cmdPrint() {
	for i := 2; i < len(os.Args); i++ {

		if strings.Contains(os.Args[i], query) {
			ans += os.Args[i] + "\n"
			//fmt.Println(os.Args[i])
		}
	}
	return
}

func checkForFile() bool {
	query = "-o"
	for i := 1; i < len(os.Args); i++ {

		if strings.Contains(os.Args[i], query) {
			return true
		}
	}
	return false
}

func writingToFile() {
	output := checkForFile()
	if output { //check if I have to write to a file
		outputFile := os.Args[len(os.Args)-1]
		b := []byte(ans)
		ioutil.WriteFile(outputFile, b, 0644)
		fmt.Println("Written to the file")

	} else { //just print the solution
		fmt.Println(ans)
	}
}

func main() {
	flag.Parse()
	query = flag.Arg(0)
	root = flag.Arg(1)

	fileInfo, err := os.Stat(root)

	if os.IsNotExist(err) { //for command line comparison
		cmdPrint()
		writingToFile()
		return
	}

	if err != nil {
		panic(err)
	}

	if fileInfo.IsDir() {
		//if its directory then explore all of it
		filepath.Walk(root, func(path string, file os.FileInfo, err error) error {

			if !file.IsDir() {
				wg.Add(1)
				go readDir(&wg, path)
			}
			return nil
		})
		wg.Wait()
		defer os.Exit(found)

	} else {
		//if its a file then just read it
		readFile(root)

	}

	writingToFile()
}
