package main

import (
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"

	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Racer struct {
    to_type string
    input string
    typed uint
    start_time time.Time
    end_time time.Time
    err error
}

func (m * Racer) nth_utf8_char(n uint) string {
    for i, c := range m.to_type {
        if i == int(n) {
            return string(c)
        }
        i++
    }
    return ""
}

func (m * Racer) from_utf8_char(n uint) string {
    r := ""
    for i, c := range m.to_type {
        if i < int(n) {
            continue
        }
        r += string(c)
    }
    return r
}

func (m Racer) Init() tea.Cmd {
    return nil
}

func (m Racer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.err = nil
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
        case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
        case tea.KeyBackspace:
            m.input = m.input[:len(m.input) - 1]
            m.typed --
            return m, nil
		}


        // if next char to type is an '\n' auto insert it
        if m.nth_utf8_char(m.typed) == "\n" {
            m.input += "¬"
            m.typed ++ 
        }

        // log.Println(m.nth_utf8_char(m.typed) + "e depois '" + m.nth_utf8_char(m.typed + 1) + "'" + " e entao '" + m.nth_utf8_char(m.typed + 2) + "'")
        
        if msg.String() != m.nth_utf8_char(m.typed) {
            m.err = errors.New("Typed wrong char")
            return m, nil
        }
        
        if m.input == "" && m.start_time.IsZero() {
            m.start_time = time.Now()
        }

        m.input += msg.String()
        m.typed++

        if m.nth_utf8_char(m.typed) == "" {
            m.typed++
        }
        
        if m.typed >= uint(len(m.to_type)) {
            m.end_time = time.Now()
            return m, tea.Quit
        }

        return m, nil

	case error:
		log.Fatal(msg)
		return m, nil
	}
    return m, nil

}

func (m Racer) View() string {
    var not_typed string
    not_typed = m.from_utf8_char(m.typed)

    var time_end string
    if !m.end_time.IsZero() {
        elapsed := m.end_time.Sub(m.start_time)
        time_end = fmt.Sprintf("\nYou took %v to type\nThat counts as %v WPM",
            m.end_time.Sub(m.start_time), 
            float64(len(strings.Split(string(m.to_type), " ")))  * 60.0 / elapsed.Seconds(),
        )
    }

    pre_view := m.input
    b := make([]byte, len(pre_view))
    for _, c := range []byte(pre_view) {
        if c == '¬' {
            b = append(b, '\n')
        } else {
            b = append(b, c)
        }
    }
    view := string(b)

    middle_char := ""
    if not_typed != ""{
        if m.err != nil {
            middle_char = color.New(color.FgHiRed).Sprint(string(not_typed[0]))
        } else {
            middle_char = color.New(color.FgWhite).Sprint(string(not_typed[0]))
        }
    }
    middle_char = "|" + middle_char
    rest := not_typed
    if rest != "" {
        rest = rest[1:]
    }

	return fmt.Sprintf(
        "Type something\n\n%s%s%s\n\n%s\n%s",
        color.New(color.FgHiCyan).Sprint(view),
        middle_char,
        color.New(color.FgHiBlack).Sprint(rest),
		"(esc to quit)",
       time_end,
	) + "\n"
}

