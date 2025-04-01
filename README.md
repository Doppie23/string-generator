# String Generator

Creates a random string based on a provided [formal grammar](https://en.wikipedia.org/wiki/Formal_grammar). Only works with context free grammars. Uses `S` as the start symbol. All other nonterminals are automatically determined based on the rules in the provided grammar file.

## Usage

```
go run main.go grammar.txt
```

## Example

Generating a random Dutch postal code:

`postal_code.txt`

```
S -> NumNumNumNumCharChar
Num -> 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9
Char -> A | B | C | D | E | F | G | H | I | J | K | L | M | N | O | P | Q | R | T | S | U | V | W | X | Y | Z
```

Note: we use `Num` and `Char` as nonterminals instead of `N` and `C` because `N` and `C` are reserved as string values in this grammar format.

```
$ go run main.go postal_code.txt
7798BH
```
