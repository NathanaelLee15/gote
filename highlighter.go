package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var vsCodeBgColor = tcell.NewHexColor(0x303446)

var vsCodeKeywordColor = tcell.NewHexColor(0xca9ee6)
var vsCodeIdentiferColor = tcell.NewHexColor(0x87a4e5)
var vsCodeSymbolColor = tcell.ColorYellow // TODO: get custom color
// var vsCodeVarColor = tcell.NewHexColor(0xea999c)
var vsCodeStrColor = tcell.NewHexColor(0x9dc583)

type HighlighterMap map[string]string

var highlighter_map = map[string]HighlighterMap{
	"go": {
		// built-in
		"package": vsCodeKeywordColor.CSS(),
		"func":    vsCodeKeywordColor.CSS(),
		"module":  vsCodeKeywordColor.CSS(),
		"import":  vsCodeIdentiferColor.CSS(),
		"main":    vsCodeIdentiferColor.CSS(),
		"(":       vsCodeSymbolColor.CSS(),
		")":       vsCodeSymbolColor.CSS(),
		"{":       vsCodeSymbolColor.CSS(),
		"}":       vsCodeSymbolColor.CSS(),
		"'":       vsCodeStrColor.CSS(),
		"\"":      vsCodeStrColor.CSS(),
		// std
		"Println":  vsCodeIdentiferColor.CSS(),
		"Printf":   vsCodeIdentiferColor.CSS(),
		"Sprintf":  vsCodeIdentiferColor.CSS(),
		"Sprintln": vsCodeIdentiferColor.CSS(),
	},
}

func HighlightParseText(text, format string) string {
	highlight_set := highlighter_map[format]
	parsed_text := text

	// TODO: process for strings ""
	// for _, c := range text {
	// 	strings.
	// }

	for token, color := range highlight_set {
		ntext := fmt.Sprintf("[%s]%s[-]", color, token)
		parsed_text = strings.ReplaceAll(parsed_text, token, ntext)
	}

	return parsed_text
}
