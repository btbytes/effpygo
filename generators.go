package effpygo

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

type Word struct {
	Index int
	Text  string
}

func isLetter(r rune) bool {
	return !(unicode.IsSpace(r) || unicode.IsPunct(r))
}

// ---

func IndexWords(text string) (result []Word) {
	var buf []rune
	var word Word

	checkFound := func() {
		if len(buf) > 0 {
			word.Text = string(buf)
			result = append(result, word)
			buf = make([]rune, 0)
		}
	}

	for i, r := range text {
		if isLetter(r) {
			// When the buffer was empty, but now we've found a
			// rune that's part of a word, mark the index as the
			// first character of a newly found word.
			if len(buf) == 0 {
				word.Index = i
			}
			// Append the rune to the current word's buffer.
			buf = append(buf, r)
		} else {
			// When the current rune is not part of a word, then
			// we may have reached the end of the word.
			checkFound()
		}
	}

	// Any runes remaining in the buffer should be returned as part of a
	// final found word.
	checkFound()

	return
}

// ---

// func getNext() rune {
// 	r, _, err := reader.ReadRune()
// 	if err != nil {
// 		return 0
// 	}
// }

// type statusChange struct {
// 	isText bool
// 	word   string
// }

// func generateStatusChanges(reader *bufio.Reader) <-chan statusChange {
// 	r, _, err := reader.ReadRune()
// 	if err != nil {
// 		return
// 	}
// 	nextStatus := isLetter(r)

// 	for {
// 		r, _, err := reader.ReadRune()
// 		if err != nil {
// 			return
// 		}
// 		next = isLetter(r)
// 		if next != current {

// 		}
// 	}
// }

func readUntil(targetStatus bool, reader *bufio.Reader) string {
	var buf bytes.Buffer
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		if isLetter(r) == targetStatus {
			reader.UnreadRune()
			break
		}
		buf.WriteRune(r)
	}
	return buf.String()
}

// func IndexWordsFromReader(in io.Reader) []Word {
// 	words := make([]Word, 0)
// 	reader := bufio.NewReader(in)
// 	index := 0
// 	for {
// 		whitespace := readUntil(true, reader)
// 		index += len(whitespace)
// 		text := readUntil(false, reader)
// 		if len(text) == 0 {
// 			break
// 		}
// 		words = append(words, Word{index, text})
// 		index += len(text)
// 	}
// 	return words
// }

// ---

func IndexWordsStream(in io.Reader) <-chan Word {
	out := make(chan Word)
	go func() {
		defer close(out)
		reader := bufio.NewReader(in)
		index := 0
		for {
			whitespace := readUntil(true, reader)
			index += len(whitespace)
			text := readUntil(false, reader)
			if len(text) == 0 {
				break
			}
			out <- Word{index, text}
			index += len(text)
		}
	}()
	return out
}
