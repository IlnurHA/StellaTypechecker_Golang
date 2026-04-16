package main

import (
	"flag"
	"fmt"
	"os"

	utils "typechecker/internal/utils"
)

func main() {
	dirPathPtr := flag.String("dirPath", "", "Path to get tests from")
	filePathPtr := flag.String("filePath", "", "Path to source code on stella")
	noExitOnErrorPtr := flag.Bool("noExitOnError", false, "If false then exits with status code 1 on typecheck error. Otherwise, continues to typecheck (in case of dirPath)")

	flag.Parse()

	if *dirPathPtr == "" && *filePathPtr == "" {
		// Getting program from stdin
		p, err := utils.GetProgramFromStdin()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		utils.Typecheck(p, !(*noExitOnErrorPtr))
		return
	}

	if *filePathPtr != "" {
		exists, err := utils.FileExists(*filePathPtr)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		if !exists {
			fmt.Fprintf(os.Stderr, "File does not exist: %s\n", *filePathPtr)
			return
		}

		utils.TypecheckFromFile(*filePathPtr, !(*noExitOnErrorPtr))
	}

	if *dirPathPtr != "" {
		files, err := utils.GetTestPaths(*dirPathPtr)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		for _, file := range files {
			utils.TypecheckFromFile(file, !(*noExitOnErrorPtr))
		}
	}
}
