//go:build fyne
// +build fyne

package main

import (
	_ "embed"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

//go:embed icons/icon.jpg
var iconData []byte

func init() {
	NewGUI = func() GUI { return &fyneGUI{} }
}

type fyneGUI struct{}

func (g *fyneGUI) Run() error {

	var targetDate = time.Now().Format("02.01.2006")
	var orderName string

	myApp := app.New()                                // создание общего окна
	window := myApp.NewWindow("Обработка файлов ГГц") // имя

	icon := fyne.NewStaticResource("icon.jpg", iconData) // создание иконки
	if icon == nil {
		log.Println("Не удалось создать ресурс иконки")
	} else {
		window.SetIcon(icon)
		myApp.SetIcon(icon)
	}

	window.Resize(fyne.NewSize(350, 300)) // размер

	inputData := widget.NewEntry() // ввод поля даты
	inputData.SetText(targetDate)  // имя поля даты

	inputOrderName := widget.NewEntry() // ввод поля заказа
	inputOrderName.SetText(orderName)   // имя поля заказа

	var listOfFiles = readFiles(targetDate) // читаем файлы за сегодня

	statusFiles := widget.NewLabel(buildNames(listOfFiles)) // собираем список файлов

	statusLabel := widget.NewLabel("Прочитаны файлы за сегодня") // строка результата

	researchBtn := widget.NewButton("Повторить поиск", func() { // ищем файлы за сегодняшнюю дату
		targetDate = inputData.Text
		listOfFiles = readFiles(targetDate)
		statusFiles.SetText(buildNames(listOfFiles))
		statusLabel.SetText("Прочитаны файлы за " + targetDate)
	})

	startBtn := widget.NewButton("Создать .xlsx", func() { // формируем файл
		orderName = inputOrderName.Text
		files := dataCreate(listOfFiles)
		results := calculateFiles(files)
		createXlsx(files, results, orderName, targetDate)
		statusLabel.SetText("Файл .xlsx успешно создан")
	})

	window.SetContent(container.NewVBox( // разметка интерфейса
		widget.NewLabel("Дата искомых файлов:"),
		inputData,
		widget.NewSeparator(),
		researchBtn,
		statusFiles,
		widget.NewLabel("Введите номер заказа:"),
		inputOrderName,
		widget.NewSeparator(),
		startBtn,
		widget.NewSeparator(),
		statusLabel,
	))
	window.ShowAndRun()
	return nil
}
