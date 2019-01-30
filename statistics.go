package zenbun

type Statistics struct {
	Documents  int
	Words      int
	WordNGrams int
	Runes      int
	RuneNGrams int
}

func (db *DB) Statistics() Statistics {
	return Statistics{
		Documents:  len(db.documents),
		Words:      len(db.wordmap.children),
		WordNGrams: len(db.wordIndex.m),
		Runes:      len(db.runemap.children),
		RuneNGrams: len(db.runeIndex.m),
	}
}
