package main

import (
	"log"
	"time"
    "net"

	tea "github.com/charmbracelet/bubbletea"
)


func initialModel() Game {
    con, err := net.Dial("tcp", "localhost:3000")
    if err != nil {
        log.Fatal(err)
    }
    buf := make([]byte, 512)
    i, err := con.Read(buf)
    if err != nil {
        log.Fatalf("Couldn't get text from server: %s\n", err)
    }
    t := string(buf[:i])

    r1 := Racer{
        to_type: t,
        input: "",
        typed: 0,
        start_time: time.Time{},
        end_time: time.Time{},
        con: &con,
    }

    r2 := Racer {
        to_type: t,
        input: "",
        typed:  0,
        start_time: time.Time{},
        end_time: time.Time{},
        con: nil,
    }

    return Game{
        player: r1,
        opponent: r2,
    }
}

func main() {
    m := initialModel()
	p := tea.NewProgram(m)

    go func () {
        buf := make([]byte, 512)
        for {
            i, err := (*m.player.con).Read(buf)
            if err != nil {
                log.Fatalf("failed to read: %v\n", err)
            }
            send := buf[:i]
            p.Send(string(send))
        }
    } ()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
