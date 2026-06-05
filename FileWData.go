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

type Data struct {
	avgvalue string
	maxvalue string
}

type FinalData struct {
	samplename string
	avgvalue   float32
	maxvalue   float32
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
		}
		var temp Data
		temp.avgvalue = data[3]
		temp.maxvalue = data[5]
		file.data = append(file.data, temp)
	}
	return file
}

// заполняет ВСЕ файлы данными в FileWData.data
func dataCreate(files []FileWData) []FileWData {
	var filledfiles []FileWData
	for _, file := range files {
		filledfiles = append(filledfiles, readExcelData(file))
	}
	return filledfiles
}

// StringToNum конвертер значений из string в float32
func StringToNum(Input Data) FinalData {
	var Output FinalData

	parseValue := func(s string) float32 {
		s = strings.TrimSpace(s)
		if s == "" {
			return 0
		}
		val, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
		if err != nil {
			return 0
		}
		return float32(val)
	}

	Output.avgvalue = parseValue(Input.avgvalue)
	Output.maxvalue = parseValue(Input.maxvalue)

	return Output
}

// подстчёт данных для R, выводит массив средних значений
func calculateData(datafile FileWData) []FinalData {
	convertedData := make([]FinalData, len(datafile.data))
	for i := range datafile.data {
		convertedData[i] = StringToNum(datafile.data[i])
	}

	calculatedData := make([]FinalData, 6)

	for a := 0; a < 6; a++ {
		var sumAvg, sumMax float32
		start := a * 5
		for i := start; i < start+5; i++ {
			sumAvg += convertedData[i].avgvalue
			sumMax += convertedData[i].maxvalue
		}
		calculatedData[a].avgvalue = sumAvg / 5
		calculatedData[a].maxvalue = sumMax / 5
	}
	return calculatedData
}

// считает все данные и выводит полный список
func calculateFiles(fileMass []FileWData) []FinalData {
	var calculatedFiles []FinalData

	for _, file := range fileMass {
		fileData := calculateData(file)

		for _, datas := range fileData {
			calculatedFiles = append(calculatedFiles, FinalData{
				samplename: file.sample,
				avgvalue:   datas.avgvalue,
				maxvalue:   datas.maxvalue,
			})
		}
	}
	return calculatedFiles
}
