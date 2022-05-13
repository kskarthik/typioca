package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

func (m model) View() string {
	var s string

	switch state := m.state.(type) {
	case MainMenu:
		s := style("typioca", m.styles.magenta)
		s += "\n\n\n"

		for i, choice := range state.choices {
			cursor := " "
			cursorClose := " "
			if state.cursor == i {
				cursor = style(">", m.styles.runningTimer)
				cursorClose = style("<", m.styles.runningTimer)
			}

			// Render the row
			s += fmt.Sprintf("%s %s%s\n\n", cursor, choice.show(m.styles), cursorClose)
		}
		termWidth, termHeight, _ := term.GetSize(0)

		s = lipgloss.NewStyle().Align(lipgloss.Left).Render(s)

		return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)
	case WordCountTestResults:
		termenv.Reset()

		rawWpmShow := "raw: " + style(strconv.Itoa(state.results.rawWpm), m.styles.greener)
		cpm := "cpm: " + style(strconv.Itoa(state.results.cpm), m.styles.greener)
		wpm := "wpm: " + style(strconv.Itoa(state.results.wpm), m.styles.runningTimer)
		givenTime := "time: " + style(state.results.time.String(), m.styles.greener)
		wordCnt := "cnt: " + style(strconv.Itoa(state.wordCnt), m.styles.greener)
		accuracy := "accuracy: " + style(fmt.Sprintf("%.1f", state.results.accuracy), m.styles.greener)
		words := "words: " + style(state.results.wordList, m.styles.greener)

		content := wpm + "\n\n" + accuracy + " " + rawWpmShow + " " + cpm + "\n" + givenTime + " " + wordCnt + " " + words

		var style = lipgloss.NewStyle().
			Align(lipgloss.Center).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(5).
			PaddingRight(5).
			BorderStyle(lipgloss.HiddenBorder()).
			BorderForeground(lipgloss.Color("2"))

		termWidth, termHeight, _ := term.GetSize(0)

		s = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, style.Render(content))

	case TimerBasedTestResults:
		termenv.Reset()

		rawWpmShow := "raw: " + style(strconv.Itoa(state.results.rawWpm), m.styles.greener)
		cpm := "cpm: " + style(strconv.Itoa(state.results.cpm), m.styles.greener)
		wpm := "wpm: " + style(strconv.Itoa(state.results.wpm), m.styles.runningTimer)
		givenTime := "time: " + style(state.results.time.String(), m.styles.greener)
		accuracy := "accuracy: " + style(fmt.Sprintf("%.1f", state.results.accuracy), m.styles.greener)
		words := "words: " + style(state.results.wordList, m.styles.greener)

		content := wpm + "\n\n" + accuracy + " " + rawWpmShow + " " + cpm + "\n" + givenTime + " " + words

		var style = lipgloss.NewStyle().
			Align(lipgloss.Center).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(5).
			PaddingRight(5).
			BorderStyle(lipgloss.HiddenBorder()).
			BorderForeground(lipgloss.Color("2"))

		termWidth, termHeight, _ := term.GetSize(0)

		s = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, style.Render(content))

	case WordCountBasedTest:
		termWidth, termHeight, _ := term.GetSize(0)

		var lineLenLimit int = 40

		reactiveLimit := (termWidth / 10) * 6
		if reactiveLimit < lineLenLimit {
			lineLenLimit = reactiveLimit
		}

		var coloredStopwatch string
		if state.stopwatch.isRunning {
			coloredStopwatch = style(state.stopwatch.stopwatch.View(), m.styles.runningTimer)
		} else {
			coloredStopwatch = style(state.stopwatch.stopwatch.View(), m.styles.stoppedTimer)
		}

		paragraphView := m.paragraphViewWordCount(lineLenLimit, state)
		lines := strings.Split(paragraphView, "\n")
		cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), state.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		// Vertical positioning
		for i := 0; i < termHeight/2-3; i++ {
			s += "\n"
		}

		avgLineLen := averageStringLen(lines[:len(lines)-1])
		indentBy := uint(termWidth/2) - (uint(avgLineLen) / 2)

		s += m.indent(coloredStopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

		if !state.stopwatch.isRunning {
			s += "\n\n\n"
			s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart", m.styles.toEnter))
		}

	case TimerBasedTest:
		termWidth, termHeight, _ := term.GetSize(0)

		var lineLenLimit int = 40

		reactiveLimit := (termWidth / 10) * 6
		if reactiveLimit < lineLenLimit {
			lineLenLimit = reactiveLimit
		}

		var coloredTimer string
		if state.timer.isRunning {
			coloredTimer = style(state.timer.timer.View(), m.styles.runningTimer)
		} else {
			coloredTimer = style(state.timer.timer.View(), m.styles.stoppedTimer)
		}

		paragraphView := m.paragraphView(lineLenLimit, state)
		lines := strings.Split(paragraphView, "\n")
		cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), state.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		// Vertical positioning
		for i := 0; i < termHeight/2-3; i++ {
			s += "\n"
		}

		avgLineLen := averageStringLen(lines[:len(lines)-1])
		indentBy := uint(termWidth/2) - (uint(avgLineLen) / 2)

		s += m.indent(coloredTimer, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

		if !state.timer.isRunning {
			s += "\n\n\n"
			s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart", m.styles.toEnter))
		}
	}

	return s
}

func (selection TimerBasedTestSettings) show(s styles) string {
	var optionsStr string
	options := []string{selection.timeSelections[selection.timeCursor].String(), selection.wordListSelections[selection.wordListCursor]}
	for i, option := range options {
		if i+1 == selection.cursor {
			// optionsStr += style("[", s.greener) + style(option, s.runningTimer) + style("]", s.greener)
			optionsStr += "[" + style(option, s.runningTimer) + "]"
		} else {
			optionsStr += "[" + style(option, s.greener) + "]"
		}
		optionsStr += " "
	}
	return fmt.Sprintf("%s %s", "Timer run", optionsStr)
}

func (selection WordCountBasedTestSettings) show(s styles) string {
	var optionsStr string
	options := []string{fmt.Sprint(selection.wordCountSelections[selection.wordCountCursor]), selection.wordListSelections[selection.wordListCursor]}
	for i, option := range options {
		if i+1 == selection.cursor {
			optionsStr += "[" + style(option, s.runningTimer) + "]"
		} else {
			optionsStr += "[" + style(option, s.greener) + "]"
		}
		optionsStr += " "
	}
	return fmt.Sprintf("%s %s", "Word count run", optionsStr)
}

func getLinesAroundCursor(lines []string, cursorLine int) []string {
	cursor := cursorLine

	// 3 lines to show
	if cursorLine == 0 {
		cursor += 3
	} else {
		cursor += 2
	}

	low := int(math.Max(0, float64(cursorLine-1)))
	high := int(math.Min(float64(len(lines)), float64(cursor)))

	return lines[low:high]
}

func dropAnsiCodes(colored string) string {
	m := regexp.MustCompile("\x1b\\[[0-9;]*m")

	return m.ReplaceAllString(colored, "")
}

func (m model) indent(block string, indentBy uint) string {
	indentedBlock := indent.String(block, indentBy) // this crashes on small windows

	return indentedBlock
}

func (m model) paragraphViewWordCount(lineLimit int, test WordCountBasedTest) string {
	paragraph := m.colorInputWordCount(test)
	paragraph += m.colorCursorWordCount(test)
	paragraph += m.colorWordsToEnterWordCount(test)

	wrapped := wrapStyledParagraph(paragraph, lineLimit)

	return wrapped
}

func (m model) paragraphView(lineLimit int, test TimerBasedTest) string {
	paragraph := m.colorInput(test)
	paragraph += m.colorCursor(test)
	paragraph += m.colorWordsToEnter(test)

	wrapped := wrapStyledParagraph(paragraph, lineLimit)

	return wrapped
}

func (m model) colorInputWordCount(test WordCountBasedTest) string {
	mistakes := toKeysSlice(test.mistakes.mistakesAt)
	sort.Ints(mistakes)

	coloredInput := ""

	if len(mistakes) == 0 {

		coloredInput += styleAllRunes(test.inputBuffer, m.styles.correct)

	} else {

		previousMistake := -1

		for _, mistakeAt := range mistakes {
			sliceUntilMistake := test.inputBuffer[previousMistake+1 : mistakeAt]
			mistakeSlice := test.wordsToEnter[mistakeAt : mistakeAt+1]

			coloredInput += styleAllRunes(sliceUntilMistake, m.styles.correct)
			coloredInput += style(mistakeSlice, m.styles.mistakes)

			previousMistake = mistakeAt
		}

		inputAfterLastMistake := test.inputBuffer[previousMistake+1:]
		coloredInput += styleAllRunes(inputAfterLastMistake, m.styles.correct)
	}

	return coloredInput
}

func (m model) colorInput(test TimerBasedTest) string {
	mistakes := toKeysSlice(test.mistakes.mistakesAt)
	sort.Ints(mistakes)

	coloredInput := ""

	if len(mistakes) == 0 {

		coloredInput += styleAllRunes(test.inputBuffer, m.styles.correct)

	} else {

		previousMistake := -1

		for _, mistakeAt := range mistakes {
			sliceUntilMistake := test.inputBuffer[previousMistake+1 : mistakeAt]
			mistakeSlice := test.wordsToEnter[mistakeAt : mistakeAt+1]

			coloredInput += styleAllRunes(sliceUntilMistake, m.styles.correct)
			coloredInput += style(mistakeSlice, m.styles.mistakes)

			previousMistake = mistakeAt
		}

		inputAfterLastMistake := test.inputBuffer[previousMistake+1:]
		coloredInput += styleAllRunes(inputAfterLastMistake, m.styles.correct)
	}

	return coloredInput
}

func (m model) colorCursorWordCount(test WordCountBasedTest) string {
	cursorLetter := test.wordsToEnter[len(test.inputBuffer) : len(test.inputBuffer)+1]

	return style(cursorLetter, m.styles.cursor)
}

func (m model) colorCursor(test TimerBasedTest) string {
	cursorLetter := test.wordsToEnter[len(test.inputBuffer) : len(test.inputBuffer)+1]

	return style(cursorLetter, m.styles.cursor)
}

func (m model) colorWordsToEnterWordCount(test WordCountBasedTest) string {
	wordsToEnter := test.wordsToEnter[len(test.inputBuffer)+1:] // without cursor

	return style(wordsToEnter, m.styles.toEnter)
}

func (m model) colorWordsToEnter(test TimerBasedTest) string {
	wordsToEnter := test.wordsToEnter[len(test.inputBuffer)+1:] // without cursor

	return style(wordsToEnter, m.styles.toEnter)
}

func wrapStyledParagraph(paragraph string, lineLimit int) string {
	// XXX: Replace spaces, because wordwrap trims them out at the ends
	paragraph = strings.Replace(paragraph, " ", "·", -1)

	f := wordwrap.NewWriter(lineLimit)
	f.Breakpoints = []rune{'·'}
	f.KeepNewlines = false
	f.Write([]byte(paragraph))
	f.Close()

	paragraph = strings.Replace(f.String(), "·", " ", -1)

	return paragraph
}

func findCursorLine(lines []string, cursorAt int) int {
	lenAcc := 0
	cursorLine := 0

	for _, line := range lines {
		lineLen := len(dropAnsiCodes(line))

		lenAcc += lineLen

		if cursorAt <= lenAcc-1 {
			return cursorLine
		} else {
			cursorLine += 1
		}
	}

	return cursorLine
}

func style(str string, style StringStyle) string {
	return style(str).String()
}

func styleAllRunes(runes []rune, style StringStyle) string {
	acc := ""

	for _, char := range runes {
		acc += style(string(char)).String()
	}

	return acc
}
