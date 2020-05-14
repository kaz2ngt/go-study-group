package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Fields []int

func (f *Fields) String() string {
	return fmt.Sprint(*f)
}
func (f *Fields) Set(s string) error {
	fieldRanges := strings.Split(s, ",")
	for _, fieldRange := range fieldRanges {
		fields := strings.Split(fieldRange, "-")
		switch len(fields) {
		case 1:
			field, err := strconv.Atoi(fields[0])
			if err != nil {
				return err
			}
			*f = append(*f, field)
		case 2:
			field1, err := strconv.Atoi(fields[0])
			if err != nil {
				return err
			}
			field2, err := strconv.Atoi(fields[1])
			if err != nil {
				return err
			}
			if field1 > field2 {
				return fmt.Errorf("invalid field1(%d) > field2(%d)", field1, field2)
			}
			for i := field1; i <= field2; i++ {
				*f = append(*f, i)
			}
		default:
			return fmt.Errorf("invalid fieldRange", fieldRange)
		}
	}
	return nil
}

// go-cutコマンドを実装しよう
func main() {
	_, err := cut()
	if err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		return
	}
}

func cut() (string, error) {
	var (
		fields    = make(Fields, 0, 1)
		delimiter string
	)
	flagSet := flag.NewFlagSet("cut", flag.ExitOnError)
	flagSet.Var(&fields, "f", "output fields")
	flagSet.StringVar(&delimiter, "d", ",", "text delimiter")
	flagSet.Parse(os.Args[1:])

	if flagSet.NArg() == 0 {
		flagSet.Usage()
		return "", nil
	}

	filePath := flagSet.Arg(0)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	indexList := make([]int, 0, len(fields))
	for _, field := range fields {
		index := field - 1
		if index < 0 {
			return "", fmt.Errorf("fields may not include zero")
		}
		indexList = append(indexList, index)
	}

	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		if len(indexList) == 0 {
			// 指定がない場合行をそのまま出力
			sb.WriteString(row)
		} else {
			splitRow := strings.Split(row, delimiter)
			pickList := make([]string, 0, len(indexList))
			for _, index := range indexList {
				if len(splitRow) > index {
					// cutコマンドがindex範囲外のとき空行を出力するため、範囲内の時のみ出力する
					pickList = append(pickList, splitRow[index])
				}
			}
			sb.WriteString(strings.Join(pickList, delimiter))
		}
		sb.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	result := sb.String()
	fmt.Print(result)
	return result, nil
}
