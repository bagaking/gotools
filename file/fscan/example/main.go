package main

import (
	"fmt"

	"github.com/bagaking/gotools/file/fscan"
)

func main() {
	sr, err := fscan.NewAndScan(".",
		true,
		false,
		fscan.FilterIgnoreHiddenFile(),
		fscan.FilterIgnoreNameTableOfFile("vendor"),
	)
	if err != nil {
		fmt.Printf("recursive path=. failed, err= %v\n", err)
	}

	sr.WarmUp()
	fileInd, numShowForDebug := 1, 100
	for _, path := range sr.GetPaths() {
		se := sr.Get(path)
		hex, _ := se.GetSha1Hex()
		fmt.Printf("%d. %s => %s", fileInd, path, hex)
		fileInd++
		if fileInd > numShowForDebug {
			fmt.Printf("... and %d more files", sr.Len()-numShowForDebug)
			break
		}
	}

	for k, v := range sr.GetChildrenTable() {
		fmt.Println("=", k, v)
	}
}
