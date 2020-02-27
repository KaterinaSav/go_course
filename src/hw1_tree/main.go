package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"path/filepath"
	"strings"
)

const defaultCounter = 3

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


func dirTree(out io.Writer, path string, printFiles bool) error {
    level := 0
    err := visitAllDir(out, path, printFiles, level)
    if err != nil {
        return err
    }
    return nil
}

func visitAllDir(out io.Writer, path string, printFiles bool, level int) error {
    var files []os.FileInfo
    files, err := visitDir(path, printFiles)
    if err != nil {
        return err
    }
    sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
    filesLength := len(files)
    for index, file := range files {
        if (file.IsDir()) {
            lastFileIndex := filesLength - 1 == index
            filename := filepath.Join(path, "/", file.Name())
            if index == 0 {
                level++
            }
            printFile(file.Name(), level, lastFileIndex)
            err := visitAllDir(out, filename, printFiles, level)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func printFile(name string, level int, lastFileIndex bool) {
    levels := make([]int, level)
    var lastResult string
    var firstResult string
    for index, _ := range levels {
        lastIndex := index + 1 == level
        firstIndex := index == 0 && !lastIndex
        beforeLastIndex := index + 2 == level && !firstIndex
        if beforeLastIndex {
            firstResult = firstResult + strings.Repeat(" ", defaultCounter)
        }

        if lastIndex {
            if lastFileIndex {
                lastResult = "└"
            } else {
                lastResult = "├"
            }
            lastResult = lastResult + strings.Repeat("─", defaultCounter)
        }
        if !beforeLastIndex && !lastIndex {
            firstResult = firstResult + " │" + strings.Repeat(" ", defaultCounter)
        }
    }

    fmt.Println(firstResult, lastResult, name, level)
}

func visitDir(path string, printFiles bool) (files []os.FileInfo, err error) {
    f, err := os.Open(path)
    if err != nil {
        return files, err
    }
    files, err = f.Readdir(-1)
    f.Close()
    if err != nil {
        return files, err
    }
    if printFiles {
        return files, err
    }
    var onlyDirFiles []os.FileInfo
    for _, file := range files {
        if (file.IsDir()) {
            onlyDirFiles = append(onlyDirFiles, file)
        }
    }
    return onlyDirFiles, err
}
