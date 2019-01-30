package zenbun

import (
	"errors"
	"sort"
	"strings"
	"sync"

	lsd "github.com/mattn/go-lsd"
)

type DocumentID uint64

var (
	DocumentNameAlreadyExists = errors.New("document name already exists")
)

const InvalidDocumentID DocumentID = 0xffffffffffffffff

type DB struct {
	documents         map[DocumentID]string
	documentIndex     map[string]DocumentID
	blob              map[DocumentID]string
	wordmap           *wordmap
	wordIndex         *wordIndex
	runemap           *runemap
	runeIndex         *runeIndex
	words             map[wordNGramID]string
	ngram             int
	didMutex          *sync.RWMutex
	splitter          WordSplitterGenerator
	maxWordCandidates int
	nextDocumentID    DocumentID
	nextWordNGramID   wordNGramID
	nextRuneNGramID   runeNGramID
	sampleCount       int
	wordRating        func(v float64, wordLength int) float64
	runeRating        func(v float64, runeLength int) float64
	normalizer        Normalizer
}

func New() *DB {
	return &DB{
		documents:         map[DocumentID]string{},
		documentIndex:     map[string]DocumentID{},
		blob:              map[DocumentID]string{},
		words:             map[wordNGramID]string{},
		wordmap:           newWordmap(invalidWordNGramID),
		wordIndex:         newWordIndex(),
		runemap:           newRunemap(invalidRuneNGramID),
		runeIndex:         newRuneIndex(),
		ngram:             4,
		maxWordCandidates: 100,
		sampleCount:       100,
		didMutex:          new(sync.RWMutex),
		splitter:          WordSplitterGeneratorFunc(newSplitter),
		normalizer:        NormalizeFunc(func(d *string) string { return strings.ToLower(*d) }),
		wordRating: func(v float64, wordLength int) float64 {
			return v * float64(wordLength*wordLength)
		},
		runeRating: func(v float64, runeLength int) float64 {
			return v * float64(runeLength*runeLength)
		},
	}
}

func (db *DB) Splitter() WordSplitterGenerator {
	return db.splitter
}

func (db *DB) SetSplitter(splitter WordSplitterGenerator) {
	db.splitter = splitter
}

func (db *DB) registerWord(word string) {
	wordid, _ := db.wordmap.get([]string{word})
	if _, ok := db.words[wordid]; ok {
		return
	}
	db.words[wordid] = word

	runes := []rune(word)
	for j := 0; j < len(runes); j++ {
		l := db.ngram
		if l+j > len(runes) {
			l = len(runes) - j
		}
		for i := 1; i <= l; i++ {
			db.nextRuneNGramID = db.runemap.register(db.nextRuneNGramID, runes[j:j+i])
			runeid, _ := db.runemap.get(runes[j : j+i])
			db.runeIndex.add(runeid, wordid, 1.0)
		}
	}
	db.runeIndex.scale(wordid, 1.0/float64(len(runes)))
}

func (db *DB) Index(name, document string) (DocumentID, error) {
	db.didMutex.Lock()
	defer db.didMutex.Unlock()

	if _, ok := db.documentIndex[name]; ok {
		return InvalidDocumentID, DocumentNameAlreadyExists
	}
	did := db.nextDocumentID
	db.nextDocumentID++

	db.documents[did] = name
	db.documentIndex[name] = did
	db.blob[did] = document

	document = db.normalizer.Normalize(&document)
	splitter := db.splitter.New(&document)

	words := []string{}
	for word := splitter.Next(); word != ""; word = splitter.Next() {
		words = append(words, word)
		l := db.ngram
		if len(words) < l {
			l = len(words)
		}
		for i := 1; i <= l; i++ {
			db.nextWordNGramID = db.wordmap.register(db.nextWordNGramID, words[len(words)-i:])
			wordid, _ := db.wordmap.get(words[len(words)-i:])
			db.wordIndex.add(wordid, did, 1.0)
		}
		db.registerWord(word)
	}
	db.wordIndex.scale(did, 1.0/float64(len(words)))
	return did, nil
}

func (db *DB) Find(text string) []Candidate {
	text = db.normalizer.Normalize(&text)
	splitter := db.Splitter().New(&text)
	words := []string{}
	docs := map[DocumentID]float64{}

	for samp := 0; samp < db.sampleCount; samp++ {
		for word := splitter.Next(); word != ""; word = splitter.Next() {
			wz := db.FindWord(word)
			if len(wz) > 0 {
				if samp == 0 {
					word = wz[0].Word
				} else {
					word = wz.Sample().Word
				}
			}
			words = append(words, word)
			l := db.ngram
			if len(words) < l {
				l = len(words)
			}
			for i := 1; i <= l; i++ {
				wordid, ok := db.wordmap.get(words[len(words)-i:])
				if ok {
					m := db.wordIndex.get(wordid)
					for k, v := range m {
						if _, ok := docs[k]; !ok {
							docs[k] = 0
						}
						docs[k] += db.wordRating(v, i)
					}
				}
			}
		}
	}

	docs2 := []Candidate{}
	for k, v := range docs {
		docs2 = append(docs2, Candidate{
			DocumentID:   k,
			DocumentName: db.documents[k],
			Score:        v,
		})
	}

	sort.Slice(docs2, func(i, j int) bool { return docs2[i].Score > docs2[j].Score })
	return docs2
}

func (db *DB) FindWord(word string) WordCandidates {
	word = db.normalizer.Normalize(&word)
	runes := []rune(word)
	words := map[wordNGramID]float64{}

	for j := 0; j < len(runes); j++ {
		l := db.ngram
		if l+j > len(runes) {
			l = len(runes) - j
		}
		for i := 1; i <= l; i++ {
			runeid, ok := db.runemap.get(runes[j : j+i])
			if ok {
				m := db.runeIndex.get(runeid)
				for k, v := range m {
					if _, ok := words[k]; !ok {
						words[k] = 0
					}
					words[k] += db.runeRating(v, i)
				}
			}
		}
	}

	docs2 := []WordCandidate{}
	for k, v := range words {
		docs2 = append(docs2, WordCandidate{
			Word:   db.words[k],
			WordID: uint64(k),
			Score:  v,
		})
	}

	sort.Slice(docs2, func(i, j int) bool { return docs2[i].Score > docs2[j].Score })

	if len(docs2) > db.maxWordCandidates {
		docs2 = docs2[:db.maxWordCandidates]
		if wid, ok := db.wordmap.get([]string{word}); ok {
			flag := false
			for _, v := range docs2 {
				if v.Word == word {
					flag = true
					break
				}
			}
			if !flag {
				docs2 = append(docs2, WordCandidate{
					Word:   word,
					WordID: uint64(wid),
				})
			}
		}
	}

	for i, v := range docs2 {
		docs2[i].Score = float64(lsd.StringDistance(word, v.Word)) / float64(len(v.Word))
	}

	sort.Slice(docs2, func(i, j int) bool { return docs2[i].Score < docs2[j].Score })

	if len(docs2) > db.maxWordCandidates {
		docs2 = docs2[:db.maxWordCandidates]
	}
	return docs2
}
