package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
)

type Nonterminal string
type CustomString string

type Value interface{}

type Word []Value

type Grammar map[Nonterminal][]Word

func (g Grammar) GenerateRandomString(i Nonterminal) string {
	ws, ok := g[i]
	if !ok {
		log.Fatalf("Missing symbol '%s' in grammar\n", i)
	}

	w := ws[rand.Intn(len(ws))]

	r := ""
	for _, v := range w {
		switch v := v.(type) {
		case Nonterminal:
			r += g.GenerateRandomString(v)
		case CustomString:
			r += string(v)
		}
	}

	return r
}

type NonterminalLoc struct {
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

	t := make(map[Nonterminal][]string)

	for _, line := range lines {
		ps := strings.Split(line, "->")
		if len(ps) != 2 {
			log.Fatalln("Invalid grammar file, a line should be in the following format:\nS -> string")
		}
		nonterminal := Nonterminal(strings.TrimSpace(ps[0]))
		vs := strings.Split(ps[1], "|")
		for i, v := range vs {
			vs[i] = strings.TrimSpace(v)
		}

		if v, ok := t[nonterminal]; ok {
			t[nonterminal] = append(v, vs...)
		} else {
			t[nonterminal] = vs
		}
	}

	nonterminals := make([]Nonterminal, len(t))

	i := 0
	for k := range t {
		nonterminals[i] = Nonterminal(k)
		i++
	}

	g := make(Grammar)

	for _, nonterminal := range nonterminals {
		strs := t[nonterminal]
		g[nonterminal] = make([]Word, 0)
		for _, str := range strs {
			nonterminalLocs := make([]NonterminalLoc, 0)

			for _, otherNonterminal := range nonterminals {
				strOtherNonterminal := string(otherNonterminal)
				locs := indexAll(str, strOtherNonterminal)
				if len(locs) == 0 {
					continue
				}
				for _, loc := range locs {
					nonterminalLocs = append(nonterminalLocs, NonterminalLoc{loc, len(strOtherNonterminal)})
				}
			}

			if len(nonterminalLocs) == 0 {
				g[nonterminal] = append(g[nonterminal], Word{CustomString(str)})
				continue
			}

			sort.Slice(nonterminalLocs, func(i int, j int) bool {
				return nonterminalLocs[i].loc <= nonterminalLocs[j].loc
			})

			word := make(Word, 0)

			offset := 0
			i := 0

			for {
				for i < len(nonterminalLocs) && offset == nonterminalLocs[i].loc {
					word = append(word, Nonterminal(str[offset:offset+nonterminalLocs[i].length]))
					offset += nonterminalLocs[i].length
					i++
				}

				if offset < len(str) {
					var to int
					if i < len(nonterminalLocs) {
						to = nonterminalLocs[i].loc
					} else {
						to = len(str)
					}

					word = append(word, CustomString(str[offset:to]))
					offset = to
				} else {
					break
				}
			}

			g[nonterminal] = append(g[nonterminal], word)
		}
	}

	return g
}

func main() {
	var filename string
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Usage: %s input.txt\n", args[0])
		os.Exit(0)
	} else {
		filename = args[1]
	}

	g := newGrammarFromFile(filename)
	fmt.Println(g.GenerateRandomString("S"))
}
