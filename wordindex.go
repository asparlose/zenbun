package zenbun

type wordIndex struct {
	m   map[wordNGramID]map[DocumentID]float64
	rev map[DocumentID]map[wordNGramID]struct{}
}

func newWordIndex() *wordIndex {
	return &wordIndex{
		m:   map[wordNGramID]map[DocumentID]float64{},
		rev: map[DocumentID]map[wordNGramID]struct{}{},
	}
}

func (i *wordIndex) get(k wordNGramID) map[DocumentID]float64 {
	r, ok := i.m[k]
	if !ok {
		r = map[DocumentID]float64{}
	}
	return r
}

func (i *wordIndex) add(k wordNGramID, doc DocumentID, v float64) {
	if _, ok := i.rev[doc]; !ok {
		i.rev[doc] = map[wordNGramID]struct{}{}
	}
	if _, ok := i.m[k]; !ok {
		i.m[k] = map[DocumentID]float64{}
	}
	if _, ok := i.m[k][doc]; !ok {
		i.m[k][doc] = 0
		i.rev[doc][k] = struct{}{}
	}
	i.m[k][doc] += v
}

func (i *wordIndex) scale(doc DocumentID, v float64) {
	if _, ok := i.rev[doc]; !ok {
		return
	}
	for k := range i.rev[doc] {
		i.m[k][doc] *= v
	}
}

func (i *wordIndex) clear(doc DocumentID) {
	if _, ok := i.rev[doc]; !ok {
		return
	}
	for k := range i.rev[doc] {
		delete(i.m[k], doc)
	}
}
