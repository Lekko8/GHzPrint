package main

import (
	"strconv"
	"strings"
)

type Data struct {
	avgvalue string
	maxvalue string
	isValid  bool
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
		end := start + blockSize
		if end > dataLen {
			end = dataLen
		}

		for i := start; i < end; i++ {
			if !datafile.data[i].isValid {
				continue
			}
			num := stringToNum(datafile.data[i])
			sumAvg += num.avgvalue
			sumMax += num.maxvalue
			count++
		}

		if count > 0 {
			calculatedData[a].avgvalue = sumAvg / float32(count)
			calculatedData[a].maxvalue = sumMax / float32(count)
		} else {
			calculatedData[a].avgvalue = 0
			calculatedData[a].maxvalue = 0
		}
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
