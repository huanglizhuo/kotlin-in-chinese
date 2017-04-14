package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var listFile []string

func main() {
	filepath.Walk(".", listFun)
	//	readLine("/Users/lizhuohuang/Workspace/b.md")
	for _, name := range listFile {
		fmt.Println(name)
		addSpace2Sharp(name)
	}
}
func listFun(path string, f os.FileInfo, err error) error {
	strRet, _ := os.Getwd()
	strRet += path
	if strings.HasSuffix(strRet, ".md") {
		listFile = append(listFile, strRet)
	}
	return nil
}
func addSpace2Sharp(filename string) {
	tmpname := filename + ".tmp"
	inFile, _ := os.Open(filename)
	defer inFile.Close()
	//	reader := bufio.NewReader(inFile)
	outfile, _ := os.Create(tmpname)
	writer := bufio.NewWriter(outfile)
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	reg := regexp.MustCompile(`^[#]{1,4}`)
	for scanner.Scan() {
		if len(reg.FindAllString(scanner.Text(), -1)) != 0 {
			src := reg.FindString(scanner.Text())
			writer.WriteString(reg.ReplaceAllString(scanner.Text(), src+" "))
		} else {
			writer.WriteString(scanner.Text() + "\n")
		}
	}
	writer.Flush()
	exec.Command("rm", filename).Run()
	exec.Command("mv", tmpname, filename).Run()
	fmt.Println("modify " + filename + "finish")
}
