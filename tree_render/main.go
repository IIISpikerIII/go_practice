package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

func splitFilesDirs(files []os.FileInfo) ([]os.FileInfo, []os.FileInfo) {
	fileInfo := []os.FileInfo{}
	dirInfo := []os.FileInfo{}

	for idx := 0; idx < len(files); idx++ {
		file := files[idx]

		if file.IsDir() == true {
			dirInfo = append(dirInfo, file)
		} else {
			fileInfo = append(fileInfo, file)
		}
	}

	return fileInfo, dirInfo
}

func printFilesStruct(output io.Writer, files []os.FileInfo) {
	var prefix string

	for idx, fileInfo := range files {
		if idx == len(files) {
			prefix = "└"
		} else {
			prefix = "├"
		}

		fmt.Fprintln(output, prefix+fileInfo.Name()+"("+fmt.Sprint(fileInfo.Size())+")")
	}
}

func printFileStruct(output io.Writer, file os.FileInfo, deep []bool, lastPosition bool) {
	var prefix, sizeFormat string

	for _, val := range deep {
		if val {
			prefix += "│	"
		} else {
			prefix += "	"
		}
	}

	if lastPosition {
		prefix += "└───"
	} else {
		prefix += "├───"
	}

	if !file.IsDir() {
		if file.Size() > 0 {
			sizeFormat = " (" + fmt.Sprint(file.Size()) + "b)"
		} else {
			sizeFormat = " (empty)"
		}
	}
	fmt.Fprintln(output, prefix+file.Name()+sizeFormat)
}

func printTree(output io.Writer, path string, printFiles bool, deep []bool) error {
	var counterDirs int

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	files, _ := file.Readdir(-1)
	_, sliceDirs := splitFilesDirs(files)

	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for idx, file := range files {

		lastPosition := false
		if file.IsDir() && !printFiles {
			counterDirs++
			lastPosition = counterDirs == len(sliceDirs)
		} else {
			lastPosition = idx == len(files)-1
		}

		if printFiles || file.IsDir() {
			printFileStruct(output, file, deep, lastPosition)
		}

		if file.IsDir() {
			printTree(output, path+string(os.PathSeparator)+file.Name(), printFiles, append(deep, !lastPosition))
		}
	}

	return err
}

func dirTree(output io.Writer, path string, printFiles bool) error {

	deep := []bool{}
	return printTree(output, path, printFiles, deep)
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
