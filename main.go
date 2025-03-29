package main

import (
	"fmt"
	"log"
	"math/rand"
)

const fileName = "test.gram"

type Identifier string
type CustomString string

type Value interface{}

type Word []Value

type Grammar map[Identifier][]Word

func (g Grammar) GenerateRandomString(i Identifier) string {
	ws, ok := g[i]
	if !ok {
		log.Fatalf("Missing symbol '%s' in grammar\n", i)
	}

	vs := ws[rand.Intn(len(ws))]

	r := ""
	for _, v := range vs {
		switch v := v.(type) {
		case Identifier:
			r += g.GenerateRandomString(v)
		case CustomString:
			r += string(v)
		}
	}

	return r
}

func main() {
	g := make(Grammar)

	g[Identifier("S")] = []Word{[]Value{CustomString("test")}, []Value{Identifier("S"), Identifier("S")}}
	fmt.Println(g.GenerateRandomString("S"))

	// data, err := os.ReadFile(fileName)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// m := make(map[string][]string)

	// content := string(data)
	// lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
}
