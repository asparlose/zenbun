package zenbun

type wordNGramID uint64

const invalidWordNGramID wordNGramID = 0xffffffffffffffff

type wordmap struct {
	children map[string]*wordmap
	id       wordNGramID
}

func newWordmap(id wordNGramID) *wordmap {
	return &wordmap{
		children: map[string]*wordmap{},
		id:       id,
	}
}

func (m *wordmap) get(words []string) (wordNGramID, bool) {
	c, ok := m.children[words[0]]
	if !ok {
		return invalidWordNGramID, false
	}

	if len(words) == 1 {
		return c.id, true
	}
	return c.get(words[1:])
}

func (m *wordmap) register(nextID wordNGramID, words []string) wordNGramID {
	c, ok := m.children[words[0]]
	if !ok {
		c = newWordmap(nextID)
		m.children[words[0]] = c
		nextID++
	}

	if len(words) == 1 {
		return nextID
	}
	return c.register(nextID, words[1:])
}
