package zenbun

type runeNGramID uint64

const invalidRuneNGramID runeNGramID = 0xffffffffffffffff

type runemap struct {
	children map[rune]*runemap
	id       runeNGramID
}

func newRunemap(id runeNGramID) *runemap {
	return &runemap{
		children: map[rune]*runemap{},
		id:       id,
	}
}

func (m *runemap) get(runes []rune) (runeNGramID, bool) {
	c, ok := m.children[runes[0]]
	if !ok {
		return invalidRuneNGramID, false
	}

	if len(runes) == 1 {
		return c.id, true
	}
	return c.get(runes[1:])
}

func (m *runemap) register(nextID runeNGramID, runes []rune) runeNGramID {
	c, ok := m.children[runes[0]]
	if !ok {
		c = newRunemap(nextID)
		m.children[runes[0]] = c
		nextID++
	}

	if len(runes) == 1 {
		return nextID
	}
	return c.register(nextID, runes[1:])
}
