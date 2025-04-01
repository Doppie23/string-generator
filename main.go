package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
)

const fileName = "test.txt"

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

	for _, line := range lines {
		ps := strings.Split(line, "->")
		if len(ps) != 2 {
			log.Fatalln("Invalid grammar file, a line should be in the following format:\nS -> string")
		}
		ident := strings.TrimSpace(ps[0])
		vs := strings.Split(ps[1], "|")
		for i, v := range vs {
			vs[i] = strings.TrimSpace(v)
		}

		if v, ok := t[ident]; ok {
			t[ident] = append(v, vs...)
		} else {
			t[ident] = vs
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
				for _, loc := range locs {
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

			offset := 0
			i := 0
			if identLocs[i].loc == 0 {
				word = append(word, Identifier(tw[0:identLocs[i].length]))
				offset = identLocs[i].length
				i++
			}

			for {
				if offset < len(tw) {
					var to int
					if i < len(identLocs) {
						to = identLocs[i].loc
					} else {
						to = len(tw)
					}
					word = append(word, CustomString(tw[offset:to]))
					offset = to
				} else {
					break
				}

				if i < len(identLocs) {
					word = append(word, Identifier(tw[offset:offset + identLocs[i].length]))
					offset += identLocs[i].length
					i++
				}
			}

			g[key] = append(g[key], word)
		}
	}

	return g
}

func main() {
	g := newGrammarFromFile(fileName)
	fmt.Println(g.GenerateRandomString("S"))
}
