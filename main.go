package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

const (
	RamSize = 30000
	BFExt   = ".bf"
)

const (
	OP_MV_R       = ">"
	OP_MV_L       = "<"
	OP_INC_VAL    = "+"
	OP_DEC_VAL    = "-"
	OP_OUT        = "."
	OP_IN         = ","
	OP_START_LOOP = "["
	OP_END_LOOP   = "]"
)

var RAM = make([]int, RamSize)

func processError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func validateFile(file *string) error {
	if ext := path.Ext(*file); ext != BFExt {
		return fmt.Errorf("Unsupported file format: %s", ext)
	}
	if _, err := os.Stat(*file); err != nil {
		return fmt.Errorf("No such file: %s", *file)
	}
	return nil
}

func executeFile(file *string) error {
	ptr := 0

	f, err := os.Open(*file)
	processError(err)
	defer f.Close()

	buffer := make([]byte, 1)
	for {
		_, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		processError(err)

		switch string(buffer) {
		case OP_MV_R:
			ptr++
		case OP_MV_L:
			ptr--
		case OP_INC_VAL:
			RAM[ptr]++
		case OP_DEC_VAL:
			RAM[ptr]--
		case OP_OUT:
			fmt.Printf("%s", string(RAM[ptr]))
		case OP_IN:
			fmt.Scanf("%d", &RAM[ptr])
		case OP_START_LOOP:
			if RAM[ptr] == 0 {
				for {
					_, err := f.Read(buffer)
					processError(err)
					if string(buffer) == OP_END_LOOP {
						break
					}
				}
			}
		case OP_END_LOOP:
			if RAM[ptr] != 0 {
				for {
					_, err := f.Read(buffer)
					processError(err)
					if string(buffer) == OP_START_LOOP {
						break
					}
					_, err = f.Seek(-2, io.SeekCurrent)
					processError(err)
				}
			}
		}
	}
	return nil
}

func getFilename(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("No source file specified.")
	}
	return args[1], nil
}

func main() {
	file, err := getFilename(os.Args)
	processError(err)

	err = validateFile(&file)
	processError(err)

	err = executeFile(&file)
	processError(err)
}
