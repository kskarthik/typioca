package main

import (
	"fmt"
	"os"
	"time"

	input "github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func main() {
	termenv.ClearScreen()
	termenv.SetWindowTitle("typioca")

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	termenv.Reset()
	println("bye!")
}

func initTimerBasedTest(settings TimerBasedTestSettings) TimerBasedTest {
	return TimerBasedTest{
		settings: settings,
		timer: myTimer{
			timer:     timer.NewWithInterval(settings.timeSelections[settings.timeCursor], time.Second),
			duration:  settings.timeSelections[settings.timeCursor],
			isRunning: false,
			timedout:  false,
		},
		wordsToEnter: NewGenerator().Generate(settings.wordListSelections[settings.wordListCursor]),
		inputBuffer:  make([]rune, 0),
		rawInputCnt:  0,
		mistakes: mistakes{
			mistakesAt:     make(map[int]bool, 0),
			rawMistakesCnt: 0,
		},
		completed: false,
		cursor:    0,
	}
}

func initTimerBasedTestSelection() TimerBasedTestSettings {
	return TimerBasedTestSettings{
		timeSelections:     []time.Duration{time.Second * 120, time.Second * 60, time.Second * 30, time.Second * 15},
		timeCursor:         2,
		wordListSelections: []string{"dorian-gray", "frankenstein", "common-words", "pride-and-prejudice"},
		wordListCursor:     2,
		cursor:             0,
	}
}

func initWordCountBasedTestSelection() WordCountBasedTestSettings {
	return WordCountBasedTestSettings{
		wordCountSelections: []int{100, 50, 25, 10},
		wordCountCursor:     1,
		wordListSelections:  []string{"dorian-gray", "frankenstein", "common-words", "pride-and-prejudice"},
		wordListCursor:      2,
		cursor:              0,
	}
}

func initMainMenu() MainMenu {
	return MainMenu{
		choices: []MainMenuSelection{initTimerBasedTestSelection(), initWordCountBasedTestSelection()},
		cursor:  0,
	}
}

func initialModel() model {
	profile := termenv.ColorProfile()
	fore := termenv.ForegroundColor()

	return model{
		state: initMainMenu(),
		styles: styles{
			correct: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore)
			},
			toEnter: func(str string) termenv.Style {
				return termenv.String(str).Foreground(fore).Faint()
			},
			mistakes: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("1")).Underline()
			},
			cursor: func(str string) termenv.Style {
				return termenv.String(str).Reverse().Bold()
			},
			runningTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2"))
			},
			stoppedTimer: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("2")).Faint()
			},
			greener: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("6")).Faint()
			},
			magenta: func(str string) termenv.Style {
				return termenv.String(str).Foreground(profile.Color("10")).Faint()
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		input.Blink, //we probably don't need this anymore
	)
}
