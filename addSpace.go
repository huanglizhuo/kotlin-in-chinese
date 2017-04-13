package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	readLine("/Users/lizhuohuang/b.md")
}
func readLine(path string) {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	reg := regexp.MustCompile(`^##`)
	for scanner.Scan() {
		if len(reg.FindAllString(scanner.Text(), -1)) != 0 {
			fmt.Println(scanner.Text())
		}
		//fmt.Println(reg.FindAllString(scanner.Text(), -1))
		//fmt.Println(scanner.Text())
	}
}
