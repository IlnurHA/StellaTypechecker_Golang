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

	flag.Parse()

	if *dirPathPtr == "" && *filePathPtr == "" {
		// Getting program from stdin
		p, err := utils.GetProgramFromStdin()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		utils.Typecheck(p)
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

		utils.TypecheckFromFile(*filePathPtr)
	}

	if *dirPathPtr != "" {
		files, err := utils.GetTestPaths(*dirPathPtr)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		for _, file := range files {
			utils.TypecheckFromFile(file)
		}
	}
}
