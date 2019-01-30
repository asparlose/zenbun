package zenbun

type runeIndex struct {
	m   map[runeNGramID]map[wordNGramID]float64
	rev map[wordNGramID]map[runeNGramID]struct{}
}

func newRuneIndex() *runeIndex {
	return &runeIndex{
		m:   map[runeNGramID]map[wordNGramID]float64{},
		rev: map[wordNGramID]map[runeNGramID]struct{}{},
	}
}

func (i *runeIndex) get(k runeNGramID) map[wordNGramID]float64 {
	r, ok := i.m[k]
	if !ok {
		r = map[wordNGramID]float64{}
	}
	return r
}

func (i *runeIndex) add(k runeNGramID, doc wordNGramID, v float64) {
	if _, ok := i.rev[doc]; !ok {
		i.rev[doc] = map[runeNGramID]struct{}{}
	}
	if _, ok := i.m[k]; !ok {
		i.m[k] = map[wordNGramID]float64{}
	}
	if _, ok := i.m[k][doc]; !ok {
		i.m[k][doc] = 0
		i.rev[doc][k] = struct{}{}
	}
	i.m[k][doc] += v
}

func (i *runeIndex) scale(doc wordNGramID, v float64) {
	if _, ok := i.rev[doc]; !ok {
		return
	}
	for k := range i.rev[doc] {
		i.m[k][doc] *= v
	}
}

func (i *runeIndex) clear(doc wordNGramID) {
	if _, ok := i.rev[doc]; !ok {
		return
	}
	for k := range i.rev[doc] {
		delete(i.m[k], doc)
	}
}
