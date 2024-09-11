package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RunProgram() {

}

func SaveFile(path, content string) {
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

func LoadFileIntoTextView(path string, textArea *tview.TextArea) {
	arr, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to load file content: %s\n", path)
	} else {
		textArea.SetText(string(arr), false)
	}
}

func main() {
	f, err := os.OpenFile("./demo/testlogfile", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	current_file := "./demo/demo.go"

	app := tview.NewApplication()

	textArea := tview.NewTextArea().
		SetWrap(false).
		SetPlaceholder("Enter text here...")
	textArea.SetTitle(current_file).SetBorder(true)

	LoadFileIntoTextView(current_file, textArea)

	helpInfo := tview.NewTextView().
		SetText(" Press F1 for help, press Ctrl-C to exit")
	position := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)
	pages := tview.NewPages()

	updateInfos := func() {
		fromRow, fromColumn, toRow, toColumn := textArea.GetCursor()
		if fromRow == toRow && fromColumn == toColumn {
			position.SetText(fmt.Sprintf("Row: [yellow]%d[white], Column: [yellow]%d ", fromRow, fromColumn))
		} else {
			position.SetText(fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn))
		}
	}

	textArea.SetMovedFunc(updateInfos)
	updateInfos()

	mainView := tview.NewGrid().
		SetRows(0, 1).
		AddItem(textArea, 0, 0, 1, 2, 0, 0, true).
		AddItem(helpInfo, 1, 0, 1, 1, 0, 0, false).
		AddItem(position, 1, 1, 1, 1, 0, 0, false)

	help1, help2, help3 := CreateHelpItems()

	help := tview.NewFrame(help1).
		SetBorders(1, 1, 0, 0, 2, 2)
	help.SetBorder(true).
		SetTitle("Help").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				pages.SwitchToPage("main")
				return nil
			} else if event.Key() == tcell.KeyEnter {
				switch {
				case help.GetPrimitive() == help1:
					help.SetPrimitive(help2)
				case help.GetPrimitive() == help2:
					help.SetPrimitive(help3)
				case help.GetPrimitive() == help3:
					help.SetPrimitive(help1)
				}
				return nil
			}
			return event
		})

	pages.AddAndSwitchToPage("main", mainView, true).
		AddPage("help", tview.NewGrid().
			SetColumns(0, 64, 0).
			SetRows(0, 22, 0).
			AddItem(help, 1, 1, 1, 1, 0, 0, true), true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			pages.ShowPage("help") //TODO: Check when clicking outside help window with the mouse. Then clicking help again.
			return nil
		} else if event.Key() == tcell.KeyCtrlS {
			log.Println("Saving...")
			SaveFile(current_file, textArea.GetText())
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
