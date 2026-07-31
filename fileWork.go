package main

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"
)

type FileWData struct {
	fileName   string
	sample     string
	sampleTime time.Time
	data       []Data
}

// Читает файлы и собирает из них данные (кроме измерений)
func readFiles(targetDate string) []FileWData {

	var sampleFiles []FileWData
	var filteredFiles []os.DirEntry

	files, err := os.ReadDir(filePath) // читаем все файлы в папке
	if err != nil {
		return sampleFiles // если файлов нет, то возвращаем пустой список
	}
	for _, file := range files { // ищем подходящие файлы
		if !file.IsDir() && filepath.Ext(file.Name()) == ".csv" && slices.Contains(strings.Split(targetDate, " "), fileDateTime(file.Name()).Format("02.01.2006")) {
			filteredFiles = append(filteredFiles, file) // кладём подходящий
		}
	}
	sorting(filteredFiles) // сортируем
	for _, file := range filteredFiles {
		sampleFiles = append(sampleFiles, createFileWData(file)) // заполняем список файлов файлами
	}

	return sampleFiles // []FileWData без измеренных данных
}

func createFileWData(file os.DirEntry) FileWData {
	var first FileWData                          // выбираем файл
	first.fileName = file.Name()                 // заполняем имя файла
	first.sampleTime = fileDateTime(file.Name()) // заполняем время создания
	first.sample = sampleName(first.fileName)    // заполняем имя образца
	return first
}

// Берём имя образца
func sampleName(filename string) string {
	base := strings.TrimSuffix(filename, ".csv")
	parts := strings.Split(base, "_")

	if len(parts) < 5 {
		log.Printf("Предупреждение: неожиданный формат имени файла: %s", filename)
		return base
	}
	return strings.Join(parts[4:], "_")
}

// Берём дату и время файла (из имени)
func fileDateTime(filename string) time.Time {
	layout := "20060102150405"
	DateTime, err := time.Parse(layout, strings.Split(strings.Split(filename, "_")[2]+strings.Split(filename, "_")[3], ".")[0])
	if err != nil {
		log.Println("Ошибка чтения даты и времени из имени файла: ", err)
	}
	return DateTime
}

// Сортировка по дате
func sorting(files []os.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		info1 := fileDateTime(files[i].Name())
		info2 := fileDateTime(files[j].Name())
		return info1.Before(info2)
	})
}

// Собираем список файлов берущихся в работу
func buildNames(files []FileWData) string {
	if len(files) == 0 {
		return ""
	}
	var filenames strings.Builder
	for _, file := range files {
		if filenames.Len() > 0 {
			filenames.WriteString("\n")
		}
		filenames.WriteString(file.fileName)
	}
	return filenames.String()
}

/*
// Создание .log файла для записи ошибок
func setupLogging() {
	filename := fmt.Sprintf("GHz_%s.log", time.Now().Format("2006-01-02")) // создаём .log файл

	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Не могу создать лог-файл:", err)
	}

	log.SetOutput(logFile)
	log.Println("Запись в .log запущена")
}
*/
