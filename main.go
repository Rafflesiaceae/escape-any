package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func assert(some bool) {
	if !some {
		panic(some)
	}
}

func main() {

	args := os.Args[1:]
	farg := args[0]

	formatStr := []rune(farg)
	towardsC := formatStr[0]
	separatorC := formatStr[1]
	// assert(formatStr[2] == '|')

	inputStrRaw := formatStr[3:]

	var inputStr []rune

	{ // interpret inputStrRaw
		for i := 0; i < len(inputStrRaw); i++ {
			c := inputStrRaw[i]
			hasNextC := i < len(inputStrRaw)-1
			switch c {
			case '\\':
				if hasNextC {
					nextC := inputStrRaw[i+1]
					switch nextC {
					case '\\':
						inputStr = append(inputStr, '\\')
						i++
					case 'n':
						inputStr = append(inputStr, '\n')
						i++
					case 't':
						inputStr = append(inputStr, '\t')
						i++
					case '"':
						inputStr = append(inputStr, '"')
						i++
					default:
						inputStr = append(inputStr, '\\', nextC)
					}
				} else {
					inputStr = append(inputStr, c)
				}
			default:
				inputStr = append(inputStr, c)
			}
		}
	}

	// parse formatStr
	var inputMap [][][]rune
	var inTowards = false
	from := make([]rune, 0)
	to := make([]rune, 0)

	snap := func() {
		inTowards = false
		inputMap = append(inputMap, [][]rune{from, to})
		from = make([]rune, 0)
		to = make([]rune, 0)
	}
	var lastC rune
	for _, c := range inputStr {
		lastC = c
		switch c {
		case towardsC:
			inTowards = true
		case separatorC:
			snap()
		default:
			if inTowards {
				to = append(to, c)
			} else {
				from = append(from, c)
			}
		}
	}
	if lastC != separatorC {
		snap()
	}

	// read stdin
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	inpStr := string(bytes)

	var result []rune

	// escape/interpret
	inp := []rune(inpStr)
	inpLen := len(inp)
	for i := 0; i < inpLen; i++ {
		c := inp[i]
		inputMapMatched := false
		for _, t := range inputMap {
			from := t[0]
			to := t[1]

			if (i + len(from)) > inpLen {
				continue
			}

			matches := true
			for j, k := range from {
				if inp[i+j] != k {
					matches = false
					break
				}
			}

			if matches {
				i = i + len(from) - 1
				result = append(result, to...)
				inputMapMatched = true
				break
			}
		}

		if !inputMapMatched {
			result = append(result, c)
		}
	}

	fmt.Fprint(os.Stdout, string(result))
}
