package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// создание .xlsx файла
func createrXlsx(listOfFiles []FileWData, results []FinalData, zakaz string, targetDate string) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Ошибка при создании файла .xlsx: %v", err)
		}
	}()

	boldStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 2, // 2 = 0.00
		Font: &excelize.Font{
			Bold: true,
		},
	})
	numberStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 2,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})

	if err != nil {
		fmt.Println("Ошибка при создании стиля:", err)
		return
	}

	for idx := range listOfFiles { // создание листов с сырыми данными
		index, err := f.NewSheet(listOfFiles[idx].sample)
		if err != nil {
			log.Printf("Ошибка создания листа %s: %v", listOfFiles[idx].sample, err)
			continue
		}
		csvFile, err := os.Open(filepath.Join("C:\\P3-34 measurments", listOfFiles[idx].filename))
		if err != nil {
			log.Printf("Ошибка открытия файла %s: %v", listOfFiles[idx].filename, err)
			continue
		}

		defer csvFile.Close()

		reader := csv.NewReader(csvFile)
		reader.Comma = ';'
		reader.FieldsPerRecord = -1
		reader.LazyQuotes = true

		rows, err := reader.ReadAll()
		if err != nil {
			log.Printf("Ошибка чтения файла %s: %v", listOfFiles[idx].filename, err)
			continue
		}

		for rowIdx, row := range rows {
			rowNum := rowIdx + 1 // номера строк и столбцов в Excel начинаются с 1
			for colIdx, cellValue := range row {
				colName, err := excelize.CoordinatesToCellName(colIdx+1, rowNum)
				if err != nil {
					log.Printf("Ошибка конвертации координат: %v", err)
					continue
				}
				f.SetCellStr(listOfFiles[idx].sample, colName, cellValue)
			}
		}
		f.SetCellValue(listOfFiles[idx].sample, "F1", zakaz)                   // вставка номера заказа
		f.SetCellValue(listOfFiles[idx].sample, "G1", listOfFiles[idx].sample) // вставка имени образца

		f.SetActiveSheet(index) // активация листа
	}

	for _, sheetName := range f.GetSheetList() { // перевод данных из текста в числа
		if err != nil {
			fmt.Println("Ошибка при применении стиля к столбцу:", err)
			return
		}
		for row := 10; row <= 42; row++ {
			for col := 4; col <= 7; col++ {
				cell, _ := excelize.CoordinatesToCellName(col, row)

				if val, _ := f.GetCellValue(sheetName, cell); val != "" {
					if num, err := strconv.ParseFloat(strings.ReplaceAll(val, ",", "."), 64); err == nil {
						f.SetCellValue(sheetName, cell, num)
						f.SetCellStyle(sheetName, "D12", "G42", numberStyle)

					}
				}
			}
		}
		err = f.SetCellStyle(sheetName, "D12", "G42", numberStyle)
		err = f.SetColStyle(sheetName, "D", boldStyle)
		err = f.SetColStyle(sheetName, "F", boldStyle)
	}

	index, err := f.NewSheet("Результат")
	if err != nil {
		log.Printf("ОШИБКА: %v", err)
		return
	}
	f.SetCellValue("Результат", "A1", "Группа")
	f.SetCellValue("Результат", "B1", "average value мкВт/см2")
	f.SetCellValue("Результат", "D1", "Группа")
	f.SetCellValue("Результат", "E1", "max value мкВт/см2")
	f.SetCellValue("Результат", "G1", zakaz)

	for i, data := range results {
		row := i + 2 // начинаем со 2-й строки (1-я строка - заголовки)

		cellA := fmt.Sprintf("A%d", row)
		cellB := fmt.Sprintf("B%d", row)
		cellD := fmt.Sprintf("D%d", row)
		cellE := fmt.Sprintf("E%d", row)

		f.SetCellValue("Результат", cellA, data.samplename)
		f.SetCellValue("Результат", cellB, data.avgvalue)
		f.SetCellValue("Результат", cellD, data.samplename)
		f.SetCellValue("Результат", cellE, data.maxvalue)
	}
	styleID, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: &[]string{"0.00"}[0], // формат 0.00
		Alignment: &excelize.Alignment{
			Horizontal: "center", // выравнивание по центру
		},
	})
	styleTxt, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left", // выравнивание слева
		},
	})
	if err != nil {
		fmt.Println("Ошибка создания стиля:", err)
		return
	}

	if err := f.SetColStyle("Результат", "B", styleID); err != nil {
		fmt.Println("Ошибка применения стиля:", err)
		return
	}
	if err := f.SetColStyle("Результат", "E", styleID); err != nil {
		fmt.Println("Ошибка применения стиля:", err)
		return
	}
	if err := f.SetColStyle("Результат", "A", styleTxt); err != nil {
		fmt.Println("Ошибка применения стиля:", err)
		return
	}
	if err := f.SetColStyle("Результат", "D", styleTxt); err != nil {
		fmt.Println("Ошибка применения стиля:", err)
		return
	}
	f.SetActiveSheet(index)                                         // активация листа
	f.DeleteSheet("Sheet1")                                         // удаление лишнего листа появившегося при создании файла
	f.SaveAs("Результаты " + zakaz + " за " + targetDate + ".xlsx") // сохранение файла
}
