package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

type SelectSubcommand struct{}

func (sub *SelectSubcommand) Name() string {
	return "select"
}
func (sub *SelectSubcommand) Aliases() []string {
	return []string{}
}
func (sub *SelectSubcommand) Description() string {
	return "Extract specified columns."
}

func (sub *SelectSubcommand) Run(args []string) {
	fs := flag.NewFlagSet(sub.Name(), flag.ExitOnError)
	var columnsString string
	var exclude bool
	fs.StringVar(&columnsString, "columns", "", "Columns to select")
	fs.StringVar(&columnsString, "c", "", "Columns to select (shorthand)")
	fs.BoolVar(&exclude, "exclude", false, "Whether to exclude the specified columns")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
	if columnsString == "" {
		fmt.Fprintf(os.Stderr, "Missing required argument --columns")
		os.Exit(1)
	}
	columns := GetArrayFromCsvString(columnsString)

	inputCsvs, err := GetInputCsvs(fs.Args(), 1)
	if err != nil {
		panic(err)
	}

	if exclude {
		ExcludeColumns(inputCsvs[0], columns)
	} else {
		SelectColumns(inputCsvs[0], columns)
	}
}

func ExcludeColumns(inputCsv AbstractInputCsv, columns []string) {
	writer := csv.NewWriter(os.Stdout)

	// Get the column indices to exclude.
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}
	columnIndices := GetIndicesForColumnsOrPanic(header, columns)
	columnIndicesToExclude := make(map[int]bool)
	for _, columnIndex := range columnIndices {
		columnIndicesToExclude[columnIndex] = true
	}

	outrow := make([]string, len(header)-len(columnIndicesToExclude))

	// Write header
	curIdx := 0
	for index, elem := range header {
		_, exclude := columnIndicesToExclude[index]
		if !exclude {
			outrow[curIdx] = elem
			curIdx++
		}
	}

	writer.Write(outrow)
	writer.Flush()

	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		curIdx = 0
		for index, elem := range row {
			_, exclude := columnIndicesToExclude[index]
			if !exclude {
				outrow[curIdx] = elem
				curIdx++
			}
		}
		writer.Write(outrow)
		writer.Flush()
	}
}

func SelectColumns(inputCsv AbstractInputCsv, columns []string) {
	writer := csv.NewWriter(os.Stdout)

	// Get the column indices to write.
	header, err := inputCsv.Read()
	if err != nil {
		panic(err)
	}
	columnIndices := GetIndicesForColumnsOrPanic(header, columns)
	outrow := make([]string, len(columnIndices))
	for i, columnIndex := range columnIndices {
		outrow[i] = header[columnIndex]
	}
	writer.Write(outrow)
	writer.Flush()

	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		for i, columnIndex := range columnIndices {
			outrow[i] = row[columnIndex]
		}
		writer.Write(outrow)
		writer.Flush()
	}
}