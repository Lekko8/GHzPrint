package main

import "log"

//go-winres simply --icon icon.png для fyne
//rsrc -ico icon.ico -o rsrc.syso для walk
//-ldflags="-s -w -H windowsgui -buildmode=exe"
//rsrc -manifest walk.manifest -ico icon.ico -o rsrc.syso

var NewGUI func() GUI

func main() {
	log.Println("Приложение запущено")
	//setupLogging() // логгирование (включить для отладки)
	if NewGUI == nil {
		log.Fatal("Не выбран GUI-фреймворк. Используйте тег сборки")
	}
	gui := NewGUI()
	if err := gui.Run(); err != nil {
		log.Fatal(err)
	}

}
