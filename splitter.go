package zenbun

import "unicode"

type WordSplitterGenerator interface {
	New(*string) WordSplitter
}

type WordSplitterGeneratorFunc func(*string) WordSplitter

func (f WordSplitterGeneratorFunc) New(document *string) WordSplitter {
	return f(document)
}

type WordSplitter interface {
	Next() string
}

type splitter struct {
	document []rune
	index    int
}

func newSplitter(document *string) WordSplitter {
	s := &splitter{
		document: []rune(*document),
	}
	for s.index = 0; s.index < len(s.document) && !unicode.IsLetter(s.document[s.index]); s.index++ {
	}
	return s
}

func (s *splitter) Next() string {
	w := []rune{}
	for ; s.index < len(s.document) && unicode.IsLetter(s.document[s.index]); s.index++ {
		w = append(w, s.document[s.index])
	}
	for ; s.index < len(s.document) && !unicode.IsLetter(s.document[s.index]); s.index++ {
	}
	return string(w)
}
