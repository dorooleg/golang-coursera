package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func dirTreeImpl(out io.Writer, path string, printFiles bool, prefix string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	listFiles, err := file.Readdir(0)
	if err != nil {
		return err
	}

	sort.Slice(listFiles, func(i, j int) bool {
		lName := listFiles[i].Name()
		rName := listFiles[j].Name()
		if !printFiles {
			if !listFiles[i].IsDir() {
				lName = ""
			}

			if !listFiles[j].IsDir() {
				rName = ""
			}
		}
		return lName < rName
	})

	for i, file := range listFiles {
		if !printFiles && !file.IsDir() {
			continue
		}

		currentPrefix := prefix + "├───"
		nextLevelPrefix := prefix + "│	"
		if i+1 == len(listFiles) {
			currentPrefix = prefix + "└───"
			nextLevelPrefix = prefix + "	"
		}

		suffix := " (empty)"
		if file.IsDir() {
			suffix = ""
		} else if file.Size() > 0 {
			suffix = " (" + strconv.FormatInt(file.Size(), 10) + "b)"
		}

		fmt.Fprintln(out, currentPrefix+file.Name()+suffix)
		if file.IsDir() {
			err = dirTreeImpl(out, path+string(os.PathSeparator)+file.Name(), printFiles, nextLevelPrefix)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeImpl(out, path, printFiles, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
