package fuzzystrmatcher

import (
	"hash"
	"hash/fnv"
	"strings"
	"unicode/utf8"
	"unsafe"

	"github.com/sajari/fuzzy"
)

// Matcher is a fuzzy string matcher
type Matcher struct {
	empty                   bool
	minLen                  int
	minWordLen              int
	minWordLenWithSpellErr  int
	options                 []matcherOption
	wordHasher              hash.Hash64
	model                   *fuzzy.Model
	allOptionWordsWithCount map[uint64]bool
	ignoredWordsCacheMap    map[uint64]uint64
}

type matcherOption struct {
	minWords int
	// key = word hash, value = used by the match method to determin if a match was made
	wordsHashes []uint64
}

// Compile creates a new Matcher for the input string
func Compile(ins ...string) *Matcher {
	if len(ins) == 0 {
		return &Matcher{
			empty: true,
		}
	}

	m := &Matcher{
		minLen:                  999,
		minWordLen:              999,
		options:                 []matcherOption{},
		wordHasher:              fnv.New64(),
		model:                   fuzzy.NewModel(),
		allOptionWordsWithCount: map[uint64]bool{},
		ignoredWordsCacheMap:    map[uint64]uint64{},
	}
	m.model.SetThreshold(1)
	m.model.SetDepth(2)
	m.model.SetUseAutocomplete(false)

	wordsToTrain := map[string]bool{}
	for _, in := range ins {
		normalizedIn := NormalizeString(in)

		wordHashes := []uint64{}
		wordsCount := 0
		sentenceLen := 0

	outerLoop:
		for _, word := range strings.Split(normalizedIn, " ") {
			for _, blackListedWord := range allBlackListedWords {
				if word == blackListedWord {
					continue outerLoop
				}
			}
			wordsToTrain[word] = true
			wordLen := len(word)
			if wordLen < m.minWordLen {
				m.minWordLen = wordLen
			}
			wordsCount++
			sentenceLen += wordLen + 1

			m.wordHasher.Reset()
			m.wordHasher.Write(s2b(word))
			hashedWord := m.wordHasher.Sum64()

			wordHashes = append(wordHashes, hashedWord)
			m.allOptionWordsWithCount[hashedWord] = false
		}

		m.options = append(m.options, matcherOption{
			minWords:    wordsCount,
			wordsHashes: wordHashes,
		})

		inMinLen := (sentenceLen - 1) / 4 * 3
		if inMinLen < m.minLen {
			m.minLen = inMinLen
		}
	}

	m.minWordLenWithSpellErr = m.minWordLen / 5 * 4
	wordsToTrainList := []string{}
	for word := range wordsToTrain {
		wordsToTrainList = append(wordsToTrainList, word)
	}
	m.model.Train(wordsToTrainList)

	return m
}

// Match matches the input against the compiled value
// The input is expected to be normalized already using fuzzystrmatcher.NormalizeString(..)
func (m *Matcher) Match(in ...string) bool {
	if m.empty {
		return false
	}

	// Cleanup the ignoredWordsCacheMap
	if len(m.ignoredWordsCacheMap) > 10.000 {
		removeCount := 0
		for key, hitCount := range m.ignoredWordsCacheMap {
			if hitCount < 10 {
				delete(m.ignoredWordsCacheMap, key)
				removeCount++
				if removeCount == 2000 {
					break
				}
			}
		}
	}

	for _, inputStr := range in {
		if len(inputStr) < m.minLen {
			continue
		}

		inputBytes := []byte(inputStr)
		startOfLastWord := 0
		inputBytesLastIdx := len(inputBytes) - 1

		// m.words = m.words[:0]
		wordsCount := 0
		var word []byte
		for idx, c := range inputBytes {
			if idx == inputBytesLastIdx {
				word = inputBytes[startOfLastWord:]
			} else if c == ' ' {
				word = inputBytes[startOfLastWord:idx]
			} else {
				continue
			}
			startOfLastWord = idx + 1
			if len(word) < m.minWordLenWithSpellErr {
				continue
			}
			m.wordHasher.Reset()
			m.wordHasher.Write(word)
			hashedWord := m.wordHasher.Sum64()
			if _, ok := m.allOptionWordsWithCount[hashedWord]; ok {
				m.allOptionWordsWithCount[hashedWord] = true
			} else if len(word) > 4 && m.ignoredWordsCacheMap[hashedWord] == 0 {
				// Only do a spell check once we check if the word didn't appear in the m.allOptionWordsWithCount
				// Also skip spell checking if the word is smaller than 5 characters
				suggestion := m.model.SpellCheck(b2s(word))
				if len(suggestion) != 0 {
					m.wordHasher.Reset()
					m.wordHasher.Write(s2b(suggestion))
					hashedWord = m.wordHasher.Sum64()
					if _, ok := m.allOptionWordsWithCount[hashedWord]; ok {
						m.allOptionWordsWithCount[hashedWord] = true
					}
				} else {
					m.ignoredWordsCacheMap[hashedWord]++
				}
			}
			wordsCount++
		}

	optionsLoop:
		for _, option := range m.options {
			if option.minWords > wordsCount {
				continue
			}

			for _, wordHash := range option.wordsHashes {
				if !m.allOptionWordsWithCount[wordHash] {
					continue optionsLoop
				}
			}

			// Revert changes made to m.allOptionWordsWithCount
			for key, matched := range m.allOptionWordsWithCount {
				if matched {
					m.allOptionWordsWithCount[key] = false
				}
			}

			return true
		}

		// Revert changes made to m.allOptionWordsWithCount
		for key, matched := range m.allOptionWordsWithCount {
			if matched {
				m.allOptionWordsWithCount[key] = false
			}
		}
	}

	return false
}

// NormalizeString normalizes the input.
// It changes the following things:
// - Space like characters are converted to spaces
// - Duplicated spaces are removed
// - Spaces around the string are removed
// - Non number and non letter characters are removed
// - Uppercase characters are converted to lowercase
func NormalizeString(inStr string) string {
	if len(inStr) == 0 {
		return ""
	}

	inBytes := []byte(inStr)

	for idx := len(inBytes) - 1; idx >= 0; idx-- {
		c := inBytes[idx]
		if c >= '0' && c <= '9' || c >= 'a' && c <= 'z' {
			continue
		}

		switch c {
		case '\n', '\r', '\t', ' ':
			if idx == len(inBytes)-1 {
				// Trim the space like characters around the input
				inBytes = inBytes[:idx]
			} else if idx != 0 && (inBytes[idx-1] == '\n' || inBytes[idx-1] == '\r' || inBytes[idx-1] == '\t' || inBytes[idx-1] == ' ') {
				// The character to the left is also a whitespace character, so we can remove this char
				// By doing this we remove the duplicated spaces
				inBytes = append(inBytes[:idx], inBytes[idx+1:]...)
			} else if c != ' ' {
				inBytes[idx] = ' '
			} else if idx == 0 {
				// The first character is a space, trim the front
				// We don't have to worry if the next character where spaces because
				// they would be already removed by a previous if else check
				inBytes = inBytes[1:]
			}
		default:
			if c >= 'A' && c <= 'Z' {
				// Convert uppercase to lowercase
				inBytes[idx] += 'a' - 'A'
			} else if c < utf8.RuneSelf {
				// Remove all other special ascii characters
				inBytes = append(inBytes[:idx], inBytes[idx+1:]...)
			}
		}
	}

	// Convert the inBytes to a string without copying the data
	return b2s(inBytes)
}

// b2s converts a byte slice to a string without copying
// Note that this will mean that changes made after to the byte slice will be reflected in the string and visa versa
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// s2b converts a string to a byte slice without copying
// Note that this will mean that changes made after to the string will be reflected in the byte slice and visa versa
func s2b(b string) []byte {
	return *(*[]byte)(unsafe.Pointer(&b))
}
