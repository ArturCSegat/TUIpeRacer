package main

import (
	"log"
	"net"
	"sync"
)

type input struct {
    msg []byte
    sender_id int
}

type Player struct {
    id int
    con net.Conn
}

type PlayerList struct {
    players []Player
    mutex sync.Mutex
}

func (pl * PlayerList) add (con net.Conn) int {
    pl.mutex.Lock()
    id := len(pl.players)
    p := Player{id, con}
    pl.players = append(pl.players, p)
    pl.mutex.Unlock()
    return id
}

func (pl * PlayerList) remove (id int) {
    n_p := make([]Player, 0, len(pl.players))
    for _, p := range pl.players {
        if p.id == id {
            p.con.Close()
            continue
        }
        n_p = append(n_p, p)
    }
    pl.players = n_p
}

func (pl * PlayerList) spread (in input) {
    log.Printf("%d sent %s\n", in.sender_id, string(in.msg))
    for _, p := range pl.players {
        if p.id == in.sender_id {
            continue
        }

        _, err := p.con.Write(in.msg)
        if err != nil {
            pl.remove(p.id)
        }
    }
}

var Players PlayerList
var Input chan input = make(chan input)

func main() {
    text := RandomSampleText(10, "ptbr")
    // tt := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
    to_type := ""
    i := 0
    for _, c := range text {
        if i == 100 {
            to_type += "\n"
            i = 0
        }
        to_type += string(c)
        i++
    }

    listener, err := net.Listen("tcp", "localhost:3000")
    if err != nil {
        log.Fatalf("err main con: %v", err)
    }
    // spreads player input
    go func () {
        for {
            Players.spread(<- Input)
        }
    } ()
    
    // listens for new players
    for {
        con, err := listener.Accept()
        if err != nil {
            log.Printf("err main accept: %v", err)
            continue
        }

        _, err = con.Write([]byte(to_type))
        if err != nil {
            log.Printf("err main send: %v", err)
            continue
        }

        id := Players.add(con)
        go func () {
            for { 
                inp := make([]byte, 512)
                i, err := con.Read(inp)
                if err != nil {
                    log.Printf("err main input: %v", err)
                    Players.remove(id)
                    return
                }
                Input <- input{sender_id:id, msg:inp[:i]}
            }
        } ()
    }

}
