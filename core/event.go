package core

import (
	"github.com/TimothyYe/glance/lib"
	ui "github.com/gizak/termui/v3"
)

func handleEvents() {
	uiEvents := ui.PollEvents()
	defer ui.Close()
	defer persistPosition()

	for {
		e := <-uiEvents

		if showTOC {
			switch e.ID {
			case "q", "<C-c>":
				return
			case "m":
				displayTOC()
			case "j", "<C-n>":
				updateTOCSelection(1)
			case "k", "<C-p>":
				updateTOCSelection(-1)
			case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
				appendTOCNumber(e.ID)
			case "<Enter>":
				openSelectedTOCChapter()
			}

			renderUI()
			continue
		}

		switch e.ID {
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			applyLayout(payload.Width, payload.Height)
		case "?":
			// show the help menu
			displayHelp(currentReadingText())
		case "p":
			// show the progress
			displayProgress(currentReadingText(), r.GetProgress())
		case "m":
			displayTOC()
		case "f":
			// show the frame
			displayBorder()
		case "b":
			// boss key
			displayBossKey(currentReadingText())
		case "q", "<C-c>":
			// quit
			return
		case "<C-n>":
			// show the next page
			moveReading(pageStep())
		case "<C-p>":
			// show the previous page
			moveReading(-pageStep())
		case "j", "<Space>", "<Enter>":
			if rowNumber == "" {
				// show the next page
				moveReading(pageStep())
			} else {
				// parse the row number
				if num, err := lib.ParseRowNum(rowNumber); err != nil {
					updateParagraph(err.Error())
				} else {
					moveReading(num)
				}
				rowNumber = ""
			}
		case "k":
			if rowNumber == "" {
				// show the previous page
				moveReading(-pageStep())
			} else {
				// parse the row number
				if num, err := lib.ParseRowNum(rowNumber); err != nil {
					updateParagraph(err.Error())
				} else {
					moveReading(1 - num)
				}
				rowNumber = ""
			}
		case "[":
			r.PrevChapter()
			updateParagraph(currentReadingText())
		case "]":
			r.NextChapter()
			updateParagraph(currentReadingText())
		case "G":
			// jump to the specified row
			if rowNumber == "" {
				// jump to the last row
				r.Last()
				updateParagraph(currentReadingText())
			} else {
				// parse the row number
				if num, err := lib.ParseRowNum(rowNumber); err != nil {
					updateParagraph(err.Error())
				} else {
					r.Goto(num)
					updateParagraph(currentReadingText())
				}
				rowNumber = ""
			}
		case "g":
			if rowNumber == "g" {
				// jump to the first row
				r.First()
				updateParagraph(currentReadingText())
				rowNumber = ""
			} else {
				rowNumber = "g"
			}
		case "+":
			adjustDisplayLines(1)
		case "-":
			adjustDisplayLines(-1)
		case "c":
			color++
			// change front color
			switchColor()
		case "t":
			// timer
			setTimer()
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// jump to rows
			rowNumber += e.ID
			updateParagraph(rowNumber)
		}

		renderUI()
	}
}

func persistPosition() {
	if onExit != nil {
		onExit(r.CurrentPos())
	}
}
