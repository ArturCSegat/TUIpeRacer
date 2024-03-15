package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
)

func RandomSampleText(word_count int, language string) string {
    f, err := os.Open(fmt.Sprintf("words/%s.txt", language))
    if err != nil {
        log.Fatalf("Canot get words from language: %v", err)
    }
    content, err := io.ReadAll(f)
    if err != nil {
        log.Fatalf("Canot get words from language: %v", err)
    }

    words := bytes.Split(content, []byte{'\n'})

    text := ""
    for range word_count {
        text += string(words[rand.Intn(len(words))])
        text += " "
    }
    text = strings.TrimRight(text, " ")
    return text
}
