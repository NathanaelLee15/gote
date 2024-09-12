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

func RunProgram(project string, does_auto_close bool) {
	log.Println("Running Program...")

	file := "main"
	str_cmd := fmt.Sprintf("pushd %s && go build -o ./%s && popd", project, file)
	cmd := exec.Command("bash", "-c", str_cmd)
	err := cmd.Run()
	if err != nil {
		log.Printf("Go Build Failed: %s --- %s\n", project, err.Error())
		return
	}
	log.Printf("Go Build Success: %s", cmd.String())

	eop := "&& "
	switch does_auto_close {
	case true:
		seconds := 3
		eop += fmt.Sprintf("sleep %d", seconds)
	case false:
		eop += "read -p 'press a key to exit...'"
	}
	str_cmd = fmt.Sprintf("gnome-terminal --geometry=136x43 -- bash -c \"%s/%s %s\"", project, file, eop)
	cmd = exec.Command("bash", "-c", str_cmd)
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to run program: %s --- %s\n", project+"/main", err.Error())
		return
	}
	log.Printf("Successfully ran program: %s", cmd.String())
}

func SaveFile(path, content string) {
	log.Printf("Saving (%s)...\n", path)

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

func LoadFileIntoTextArea(path string, textArea *tview.TextArea) {
	arr, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to load file content: %s\n", path)
	} else {
		sty := tcell.Style{}.Background(tcell.ColorMediumBlue)
		textArea.SetSelectedStyle(sty)

		text := string(arr)
		textArea.SetText(text, false)
	}
}

func main() {
	vsCodeBgColor := tcell.NewHexColor(0x303446)
	// vsCodeKeywordColor := tcell.NewHexColor(0xca9ee6)
	// vsCodeIdentiferColor := tcell.NewHexColor(0x87a4e5)
	// vsCodeVarColor := tcell.NewHexColor(0xea999c)
	// vsCodeStrColor := tcell.NewHexColor(0x9dc583)

	f, err := os.OpenFile("./demo/testlogfile", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	current_project := "./demo"
	current_file := current_project + "/main.go"

	app := tview.NewApplication()

	editorBgColor := vsCodeBgColor
	textArea := tview.NewTextArea().
		SetWrap(true).
		SetPlaceholder("Enter text here...")
	textArea.SetTitle("[yellow] " + current_file + " ").SetBorder(true)
	textArea.SetBackgroundColor(editorBgColor)
	textArea.SetTextStyle(textArea.GetTextStyle().Background(editorBgColor))

	LoadFileIntoTextArea(current_file, textArea)

	helpInfo := tview.NewTextView().
		SetText(" F1 help, Ctrl-C exit, Ctrl-S save")
	helpInfo.SetBackgroundColor(vsCodeBgColor)
	position := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)
	position.SetBackgroundColor(vsCodeBgColor)

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

	tvSuggestions := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[blue]At Cursor\n[blue]From AI")
	tvSuggestions.SetBackgroundColor(vsCodeBgColor)

	frSuggestions := tview.NewFrame(tvSuggestions).
		SetBorders(1, 1, 0, 0, 2, 2)
	frSuggestions.SetBorder(true).
		SetTitle(" Suggestions ").
		SetTitleColor(tcell.ColorLightGray).
		SetBorderColor(tcell.ColorBisque)
	frSuggestions.SetBackgroundColor(vsCodeBgColor)

	tvOutput := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[blue]log 1\n[blue]log 2")
	tvOutput.SetBackgroundColor(vsCodeBgColor)

	frOutput := tview.NewFrame(tvOutput).
		SetBorders(1, 1, 0, 0, 2, 2)
	frOutput.SetBorder(true).
		SetTitle(" Output ").
		SetTitleColor(tcell.ColorLightGray).
		SetBorderColor(tcell.ColorBisque)
	frOutput.SetBackgroundColor(vsCodeBgColor)

	ROWS, COLS := 12, 6

	editorRows := 8

	mainView := tview.NewGrid().
		SetRows(0, 12).
		AddItem(textArea, 0, 0, editorRows, COLS, 0, 0, true).
		AddItem(helpInfo, ROWS, 0, 1, COLS-2, 0, 0, false).
		AddItem(position, ROWS, COLS-2, 1, 2, 0, 0, false).
		AddItem(frSuggestions, editorRows, 0, 4, 3, 0, 0, false).
		AddItem(frOutput, editorRows, 3, 4, 3, 0, 0, false)

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

	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

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
			LoadFileIntoTextArea(current_file, textArea)
			// close explorer
			pages.SwitchToPage("main")
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

	explorer := tview.NewFrame(tree).
		SetBorders(1, 1, 0, 0, 2, 2)
	explorer.SetBorder(true).
		SetTitle(" Explorer ").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				pages.SwitchToPage("main")
				return nil
			} else if event.Key() == tcell.KeyEnter {
				return nil
			}
			return event
		})

	pages.AddAndSwitchToPage("main", mainView, true).
		AddPage("help", tview.NewGrid().
			SetColumns(0, 64, 0).
			SetRows(0, 22, 0).
			AddItem(help, 1, 1, 1, 1, 0, 0, true), true, false).
		AddPage("explorer", tview.NewGrid().
			SetColumns(0, 64, 0).
			SetRows(0, 22, 0).
			AddItem(explorer, 1, 1, 1, 1, 0, 0, true), true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			pages.ShowPage("help") //TODO: Check when clicking outside help window with the mouse. Then clicking help again.
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
		} else if event.Key() == tcell.KeyCtrlR {
			RunProgram(current_project, true)
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
