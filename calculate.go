package main

import (
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

// подстчёт данных для R, выводит массив средних значений
func calculateData(datafile FileWData) []FinalData {
	convertedData := make([]FinalData, len(datafile.data))
	for i := range datafile.data {
		convertedData[i] = stringToNum(datafile.data[i])
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

// конвертер значений из string в float32
func stringToNum(Input Data) FinalData {
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
