package main

import "log"

type GUI interface {
	Run() error
}

//go-winres simply --icon icon.png для fyne
//rsrc -ico icon.ico -o rsrc.syso для walk
//-ldflags="-s -w -H windowsgui -buildmode=exe"
//rsrc -manifest walk.manifest -ico icon.ico -o rsrc.syso

var NewGUI func() GUI
var filePath = "C:\\P3-34 measurments"

func main() {
	log.Println("Приложение запущено")
	//setupLogging() // логирование (для отладки, также раскомментировать функцию)
	if NewGUI == nil {
		log.Fatal("Не выбран GUI-фреймворк. Используйте тег сборки")
	}
	gui := NewGUI()
	if err := gui.Run(); err != nil {
		log.Fatal(err)
	}

}
