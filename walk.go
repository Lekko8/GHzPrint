//go:build walk
// +build walk

package main

import (
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func init() {
	NewGUI = func() GUI { return &walkGUI{} }
}

type walkGUI struct{}

func (g *walkGUI) Run() error {

	var (
		inputData      *walk.LineEdit
		inputOrderName *walk.LineEdit
		statusFiles    *walk.TextEdit
		statusLabel    *walk.Label
		mw             *walk.MainWindow
	)

	targetDate := time.Now().Format("02.01.2006")

	// Переменные для захвата в замыкания
	var capturedFiles []FileWData
	var capturedDate = targetDate
	appIcon, _ := walk.NewIconFromResourceId(2)

	// Создаём окно
	mainWindow := MainWindow{
		AssignTo: &mw,
		Title:    "Обработка файлов ГГц",
		Icon:     appIcon,
		Size:     Size{Width: 350, Height: 400},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "Дата искомых файлов:"},
			LineEdit{AssignTo: &inputData},

			Label{Text: "Найденные файлы:"},
			TextEdit{AssignTo: &statusFiles, ReadOnly: true, VScroll: true},

			Label{Text: "Введите номер заказа:"},
			LineEdit{AssignTo: &inputOrderName},

			PushButton{
				Text: "Повторить поиск",
				OnClicked: func() {
					capturedDate = inputData.Text()
					capturedFiles = readFiles(capturedDate)
					statusFiles.SetText(buildNames(capturedFiles))
					statusLabel.SetText("Прочитаны файлы за " + capturedDate)
				},
			},

			PushButton{
				Text: "Создать .xlsx",
				OnClicked: func() {
					orderName := inputOrderName.Text()
					files := dataCreate(capturedFiles)
					results := calculateFiles(files)
					createrXlsx(files, results, orderName, capturedDate)
					statusLabel.SetText("Файл .xlsx успешно создан")
				},
			},

			Label{AssignTo: &statusLabel},
		},
	}

	// Создаём окно
	if err := mainWindow.Create(); err != nil {
		return err
	}

	// Устанавливаем значения
	inputData.SetText(targetDate)

	// ИНИЦИАЛИЗАЦИЯ ПОИСКА ПРИ ЗАПУСКЕ
	capturedDate = targetDate
	capturedFiles = readFiles(targetDate)
	statusFiles.SetText(buildNames(capturedFiles))
	statusLabel.SetText("Прочитаны файлы за " + targetDate)

	// Запускаем окно
	mw.Run()

	return nil
}
