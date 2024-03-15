package main

import (
	"log"
	"time"

	txtInput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


func initialModel() Racer {
    tt := RandomSampleText(1, "ptbr")
    // tt := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
    t := ""
    i := 0
    for _, c := range tt {
        if i == 100 {
            t += "\n"
            i = 0
        }
        t += string(c)
        i++
    }

    in := txtInput.New()
    in.Focus()
    in.CharLimit = len(t)
    in.Width = len(t)

    return Racer{
        to_type: t,
        input: in.Value(),
        typed: 0,
        start_time: time.Time{},
        end_time: time.Time{},
    }
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
