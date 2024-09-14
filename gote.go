package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/nathanaellee15/tview"
)

func main() {
	/// setup logger
	log_handle, err := os.OpenFile("./logs.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer log_handle.Close()
	log.SetOutput(log_handle)
	///

	/// Suggestions Section
	tvSuggestions := tview.NewTextView().
		SetDynamicColors(true)

	tvSuggestions.SetBackgroundColor(vsCodeBgColor)

	frSuggestions := tview.NewFrame(tvSuggestions).
		SetBorders(1, 1, 0, 0, 2, 2)
	frSuggestions.SetBorder(true).
		SetTitle(" Suggestions ").
		SetTitleColor(tcell.ColorLightGray).
		SetBorderColor(tcell.ColorBisque)
	frSuggestions.SetBackgroundColor(vsCodeBgColor)
	///

	/// Output Section
	tvOutput := tview.NewTextView().
		SetDynamicColors(true)
	tvOutput.SetBackgroundColor(vsCodeBgColor)

	frOutput := tview.NewFrame(tvOutput).
		SetBorders(1, 1, 0, 0, 2, 2)
	frOutput.SetBorder(true).
		SetTitle(" Output ").
		SetTitleColor(tcell.ColorLightGray).
		SetBorderColor(tcell.ColorBisque)
	frOutput.SetBackgroundColor(vsCodeBgColor)
	///

	/// Suggestions & Output Helpers

	SetSuggestionsText := func(text string) {
		tvSuggestions.ScrollToBeginning()
		tvSuggestions.SetText(text)
	}
	// PushSuggestionsText := func(text string) {
	// 	tvSuggestions.SetText(fmt.Sprintf("%s%s\n", tvSuggestions.GetText(false), text))
	// 	tvSuggestions.ScrollToEnd()
	// }

	SetOutputText := func(text string) {
		tvOutput.ScrollToBeginning()
		tvOutput.SetText(text)
	}
	PushOutputText := func(text string) {
		tvOutput.SetText(fmt.Sprintf("%s%s\n", tvOutput.GetText(false), text))
		tvOutput.ScrollToEnd()
	}
	///

	/// Project Path and Initial file target
	current_project := "."
	current_file := current_project + "/main.go"

	/// Editor Section
	editorBgColor := vsCodeBgColor
	textArea := tview.NewTextArea().
		SetWrap(true).
		SetPlaceholder("Write some code or whatever!")
	textArea.SetTitle("[yellow] " + current_file + " ").SetBorder(true)
	textArea.SetBackgroundColor(editorBgColor)
	textArea.SetTextStyle(textArea.GetTextStyle().Background(editorBgColor))
	///

	/// Highlighter Callback
	testFlag := false
	GoInsightCb := func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		if testFlag {
			testFlag = false
			log.Println("Test Event!")
			if textArea.HasSelection() {
				text, start, end := textArea.GetSelection()
				if start != end {
					log.Printf("Selected Text: %s\n", text)
					parsed_text := HighlightParseText(text, "go")
					SetSuggestionsText("[:::https://github.com/nathanaellee15/tview]Hyperlinks[:::-]\n[green]Info for:[blue]\n" + parsed_text)
					screen.Beep()
					tvSuggestions.Draw(screen)
				}
			}
		}
		// +1 and -2, for padding because of borders
		return x + 1, y + 1, width - 2, height - 2
	}
	///

	/// Setup App
	app := tview.NewApplication()
	pages := tview.NewPages()

	LoadFileIntoTextArea(current_file, textArea, GoInsightCb)
	///

	/// Simple Help Info & Cursor Info
	helpInfo := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Auto-Save: On[-]  F1 help, Ctrl-C exit, Ctrl-S save")
	helpInfo.SetBackgroundColor(vsCodeBgColor)
	position := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)
	position.SetBackgroundColor(vsCodeBgColor)

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
	///

	/// Editor Config
	ROWS, COLS := 12, 6
	editorRows := 8
	///

	/// Layout Editor
	mainView := tview.NewGrid().
		SetRows(0, 12).
		AddItem(textArea, 0, 0, editorRows, COLS, 0, 0, true).
		AddItem(helpInfo, ROWS, 0, 1, COLS-2, 0, 0, false).
		AddItem(position, ROWS, COLS-2, 1, 2, 0, 0, false).
		AddItem(frSuggestions, editorRows, 0, 4, 3, 0, 0, false).
		AddItem(frOutput, editorRows, 3, 4, 3, 0, 0, false)
	///

	/// Load Help Page and Items
	help1, help2, help3 := CreateHelpItems()
	help1.SetBackgroundColor(vsCodeBgColor)
	help2.SetBackgroundColor(vsCodeBgColor)
	help3.SetBackgroundColor(vsCodeBgColor)

	help := tview.NewFrame(help1).
		SetBorders(1, 1, 0, 0, 2, 2)
	help.SetBackgroundColor(vsCodeBgColor)
	help.SetBorder(true).
		SetTitle(" Help ").
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
	///

	/// Setup Explorer
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	tree.SetBackgroundColor(vsCodeBgColor)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(file.IsDir() || strings.Contains(file.Name(), "."))
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			} else {
				node.SetColor(tcell.ColorCornflowerBlue)
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}

		if strings.Contains(node.GetText(), ".") {
			log.Printf("explorer: opening file: ./%s\n", reference.(string))

			// save current file
			SaveFile(current_file, textArea.GetText())
			// change current file
			current_file = fmt.Sprintf("./%s", reference.(string))
			// set editor title new file's path/name
			new_title := fmt.Sprintf("[yellow] %s ", current_file)
			textArea.SetTitle(new_title)
			// load new file's contents
			LoadFileIntoTextArea(current_file, textArea, GoInsightCb)
			// close explorer
			pages.SwitchToPage("main")

			SetOutputText(fmt.Sprintf("[purple]Switching to %s\n", current_file))
			return
		}

		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})
	///

	/// Explorer Section
	explorer := tview.NewFrame(tree).
		SetBorders(1, 1, 0, 0, 2, 2)
	explorer.SetBackgroundColor(vsCodeBgColor)
	explorer.SetBorder(true).
		SetTitle(" Explorer ").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				pages.SwitchToPage("main")
				return nil
			} else if event.Key() == tcell.KeyEnter {
				// TODO: get cmd_str from an inputfield
				cmd_str := "echo 'hello'"

				log.Printf("Executing CMD: %s\n", cmd_str)
				cmd := exec.Command("bash", "-c", cmd_str)
				output, err := cmd.Output()
				var output_str string
				if err != nil {
					output_str = err.Error()
					log.Printf("Error running cmd: %s\n", output_str)
				} else {
					output_str = string(output)
					log.Printf("StdOut from cmd: %s\n", output_str)
				}
				PushOutputText(output_str)

				return nil
			}
			return event
		})
	///

	/// Attach Pages
	pages.AddAndSwitchToPage("main", mainView, true).
		AddPage("help", tview.NewGrid().
			SetColumns(0, 64, 0).
			SetRows(0, 22, 0).
			AddItem(help, 1, 1, 1, 1, 0, 0, true), true, false).
		AddPage("explorer", tview.NewGrid().
			SetColumns(0, 64, 0).
			SetRows(0, 22, 0).
			AddItem(explorer, 1, 1, 1, 1, 0, 0, true), true, false)
	///

	/// Handle Key Events
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			pages.ShowPage("help") //TODO: Check when clicking outside help window with the mouse. Then clicking help again.
			return nil
		} else if event.Key() == tcell.KeyCtrlP {
			testFlag = true

			return nil
		} else if event.Key() == tcell.KeyCtrlSpace {
			name, _ := pages.GetFrontPage()
			if name != "explorer" {
				pages.ShowPage("explorer")
			} else {
				pages.SwitchToPage("main")
			}
			return nil
		} else if event.Key() == tcell.KeyCtrlS {
			SaveFile(current_file, textArea.GetText())
			PushOutputText(fmt.Sprintf("[green]Saved %s", current_file))
		} else if event.Key() == tcell.KeyCtrlR {
			RunProgram(current_project, true)
		} else if event.Key() == tcell.KeyCtrlC {
			log.Println("Saving On Exit...")
			SaveFile(current_file, textArea.GetText())
		}
		return event
	})
	///

	/// Launch
	if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
