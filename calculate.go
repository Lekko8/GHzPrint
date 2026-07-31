package main

import (
	"strconv"
	"strings"
)

type Data struct {
	avgValue string
	maxValue string
	isValid  bool
}

type FinalData struct {
	sampleName string
	avgValue   float32
	maxValue   float32
}

// Считает все данные и выводит полный список
func calculateFiles(fileMass []FileWData) []FinalData {
	var calculatedFiles []FinalData

	for _, file := range fileMass {
		fileData := calculateData(file)

		for _, datas := range fileData {
			calculatedFiles = append(calculatedFiles, FinalData{
				sampleName: file.sample,
				avgValue:   datas.avgValue,
				maxValue:   datas.maxValue,
			})
		}
	}
	return calculatedFiles
}

// Подсчёт данных для R, выводит массив средних значений
func calculateData(datafile FileWData) []FinalData {
	if len(datafile.data) == 0 {
		return []FinalData{}
	}

	blockSize := 5
	dataLen := len(datafile.data)

	numBlocks := dataLen / blockSize
	if dataLen%blockSize != 0 {
		numBlocks++
	}

	calculatedData := make([]FinalData, numBlocks)

	for a := 0; a < numBlocks; a++ {
		var sumAvg, sumMax float32
		var count int

		start := a * blockSize
		end := min(start+blockSize, dataLen)

		for i := start; i < end; i++ {
			if !datafile.data[i].isValid {
				continue
			}
			num := stringToNum(datafile.data[i])
			sumAvg += num.avgValue
			sumMax += num.maxValue
			count++
		}

		if count > 0 {
			calculatedData[a].avgValue = sumAvg / float32(count)
			calculatedData[a].maxValue = sumMax / float32(count)
		} else {
			calculatedData[a].avgValue = 0
			calculatedData[a].maxValue = 0
		}
	}

	return calculatedData
}

// Конвертер значений из string в float32
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

	Output.avgValue = parseValue(Input.avgValue)
	Output.maxValue = parseValue(Input.maxValue)

	return Output
}
