package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
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

type IdentLoc struct {
	loc    int
	length int
}

func indexAll(s string, substr string) []int {
	var r []int
	offset := 0
	for {
		loc := strings.Index(s[offset:], substr)
		if loc == -1 {
			break
		}
		offset += loc
		r = append(r, offset)
		offset += len(substr)
	}
	return r
}

func newGrammarFromFile(filepath string) Grammar {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	content := string(data)
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")

	t := make(map[string][]string)

	// TODO: support `|` syntax
	for _, line := range lines {
		ps := strings.Split(line, "->")
		if len(ps) != 2 {
			log.Fatalln("Invalid grammar file") // TODO: better error message...
		}
		ident := strings.TrimSpace(ps[0])
		rest := strings.TrimSpace(ps[1])

		if v, ok := t[ident]; ok {
			t[ident] = append(v, rest)
		} else {
			t[ident] = []string{rest}
		}
	}

	idents := make([]string, len(t))

	i := 0
	for k := range t {
		idents[i] = k
		i++
	}

	g := make(Grammar)

	for _, ident := range idents {
		v := t[ident]
		key := Identifier(ident)
		g[key] = make([]Word, 0)
		for _, tw := range v {
			identLocs := make([]IdentLoc, 0)

			for _, id := range idents {
				locs := indexAll(tw, id)
				if len(locs) == 0 {
					continue
				}
				for loc := range locs {
					identLocs = append(identLocs, IdentLoc{loc, len(id)})
				}
			}

			if len(identLocs) == 0 {
				g[key] = append(g[key], Word{CustomString(tw)})
				continue
			}

			sort.Slice(identLocs, func(i int, j int) bool {
				return identLocs[i].loc <= identLocs[j].loc
			})

			word := make(Word, 0)

			// if first indent loc is 0, we first need to get that
			//
			// than we look for the next word and get it if its there and look if there is another ident, if so grab it aswell
			// repeat

			g[key] = append(g[key], word)
		}
	}

	return g
}

func main() {
	// g := make(Grammar)

	// g[Identifier("S")] = []Word{[]Value{CustomString("test")}, []Value{Identifier("S"), Identifier("S")}}
	// fmt.Println(g.GenerateRandomString("S"))

	g := newGrammarFromFile(fileName)
	fmt.Print(g)

}
