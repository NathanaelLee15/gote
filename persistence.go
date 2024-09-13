package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/nathanaellee15/tview"
)

func SaveFile(path, content string) {
	log.Printf("Saving %s\n", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("Failed to open file for saving: %s\n", path)
		return
	}
	c, err := f.WriteString(content)
	if err != nil {
		log.Println("Failed to write file contents...")
	}
	log.Printf("Wrote (%d) bytes\n", c)
}

func LoadFileIntoTextArea(path string, textArea *tview.TextArea, drawFunc func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int)) {
	arr, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to load file content: %s\n", path)
	} else {
		sty := tcell.Style{}.Background(tcell.ColorMediumBlue)
		textArea.SetSelectedStyle(sty)

		textArea.SetDrawFunc(drawFunc)

		text := string(arr)
		textArea.SetText(text, false)
	}
}
