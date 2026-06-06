package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
)

// заполняет ВСЕ файлы данными в FileWData.data
func dataCreate(files []FileWData) []FileWData {
	var filledfiles []FileWData
	for _, file := range files {
		filledfiles = append(filledfiles, readExcelData(file))
	}
	return filledfiles
}

// читает данные из .csv и заполняет FileWData.data
func readExcelData(file FileWData) FileWData {
	fileR, err := os.Open(filepath.Join("C:\\P3-34 measurments", file.filename))
	defer fileR.Close()
	if err != nil {
		log.Printf("ОШИБКА: %v", err)
		return file
	}
	reader := csv.NewReader(fileR)
	reader.Comma = ';'
	for {
		var str, _ = reader.Read()
		if str == nil || len(str) == 0 || str[0] != "interval_start" {
			continue
		}
		break
	}

	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("ОШИБКА: %v", err)
			return file
		}
		var temp Data
		temp.avgvalue = data[3]
		temp.maxvalue = data[5]
		file.data = append(file.data, temp)
	}
	return file
}
