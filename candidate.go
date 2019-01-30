package zenbun

import "math/rand"

type Candidate struct {
	DocumentName string
	DocumentID   DocumentID
	Score        float64
}

type WordCandidate struct {
	Word   string
	WordID uint64
	Score  float64
}

type WordCandidates []WordCandidate

func (c WordCandidates) Sample() WordCandidate {
	var sum float64
	a := make([]float64, len(c))
	for i, v := range c {
		if v.Score == 0 {
			a[i] = float64(len(c))
		} else {
			a[i] = 1.0 / v.Score
		}
		sum += a[i]
	}

	k := rand.Float64() * sum
	for i, v := range a {
		k -= v
		if k < 0 {
			return c[i]
		}
	}
	return c[0]
}
