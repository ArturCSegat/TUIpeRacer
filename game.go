package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Game struct {
    player Racer
    opponent Racer 
}

func (g Game) Init() tea.Cmd {
    return nil
}

func (g Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tea.KeyMsg:
        m, cmd  := g.player.Update(msg)
        v, ok := m.(Racer)
        if ok {
            g.player = v
        }
        return g, cmd
    case string:
        m, cmd  := g.opponent.Update(msg)
        v, ok := m.(Racer)
        if ok {
            g.opponent = v
        }
        return g, cmd
    }
    return g, nil
}

func (g Game) View() string {
    return fmt.Sprintf("%s\n\n\n%s\n", g.player.View(), g.opponent.View())
}

