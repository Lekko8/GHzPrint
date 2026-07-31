package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Заполняет ВСЕ файлы данными в FileWData.data
func dataCreate(files []FileWData) []FileWData {
	var result []FileWData

	for _, file := range files {
		filledFile := readCSVData(file)
		result = append(result, filledFile)
	}

	return result
}

// Читает данные из .csv и заполняет FileWData.data
func readCSVData(file FileWData) FileWData {
	csvPath := filepath.Join(filePath, file.fileName)
	csvFile, err := os.Open(csvPath)
	if err != nil {
		log.Printf("ОШИБКА: не удалось открыть файл %s: %v", file.fileName, err)
		return file
	}
	defer func() {
		if err := csvFile.Close(); err != nil {
			log.Printf("ОШИБКА: не удалось закрыть файл %s: %v", file.fileName, err)
		}
	}()

	reader := csv.NewReader(csvFile)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	var dataList []Data

	for {
		record, err := reader.Read()
		if err == io.EOF {
			log.Printf("ПРЕДУПРЕЖДЕНИЕ: файл %s не содержит данных", file.fileName)
			return file
		}
		if err != nil {
			log.Printf("ОШИБКА: %v", err)
			return file
		}
		if len(record) > 0 && strings.Contains(strings.ToLower(record[0]), "interval_start") {
			break
		}
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("ОШИБКА при чтении файла %s: %v", file.fileName, err)
			continue
		}

		if isEmptyRecord(record) {
			dataList = append(dataList, Data{
				avgValue: "",
				maxValue: "",
				isValid:  false,
			})
			continue
		}

		if len(record) < 6 || !strings.Contains(record[0], "202") {
			continue
		}

		avgVal := strings.TrimSpace(record[3])
		maxVal := strings.TrimSpace(record[5])

		_, errAvg := strconv.ParseFloat(strings.ReplaceAll(avgVal, ",", "."), 64)
		_, errMax := strconv.ParseFloat(strings.ReplaceAll(maxVal, ",", "."), 64)

		if errAvg == nil && errMax == nil {
			dataList = append(dataList, Data{
				avgValue: avgVal,
				maxValue: maxVal,
				isValid:  true,
			})
		} else {
			dataList = append(dataList, Data{
				avgValue: "",
				maxValue: "",
				isValid:  false,
			})
		}
	}

	file.data = dataList
	return file
}

// Проверяем, что запись пустая
func isEmptyRecord(record []string) bool {
	if len(record) == 0 {
		return true
	}
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}
	return true
}
