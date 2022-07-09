package translator

type grammaticalGender string

// Grammatical genders when applicable
const (
	GGenderNeuter    grammaticalGender = "n"
	GGenderMasculine grammaticalGender = "m"
	GGenderFeminine  grammaticalGender = "f"
)

type Example struct {
	usage   string
	meaning string
}

func NewExample(usage string, meaning string) Example {
	return Example{
		usage:   usage,
		meaning: meaning,
	}
}

func (e Example) Usage() string {
	return e.usage
}

func (e Example) Meaning() string {
	return e.meaning
}

type Meaning struct {
	word   string
	gender grammaticalGender
}

func NewMeaning(word string, gender string) Meaning {
	return Meaning{
		word:   word,
		gender: grammaticalGender(gender),
	}
}

func (m Meaning) Word() string {
	return m.word
}

func (m Meaning) IsNeuter() bool {
	return m.gender == GGenderNeuter
}

func (m Meaning) IsMasculine() bool {
	return m.gender == GGenderMasculine
}

func (m Meaning) IsFeminine() bool {
	return m.gender == GGenderFeminine
}

func (m Meaning) Gender() string {
	return string(m.gender)
}

type Translation struct {
	from     Language
	to       Language
	word     string
	meanings []Meaning
	examples []Example
}

func NewTranslation(word string, from Language, to Language, meanings []Meaning, examples []Example) Translation {
	return Translation{
		word:     word,
		from:     from,
		to:       to,
		meanings: meanings,
		examples: examples,
	}
}

func (t Translation) FromLanguage() Language {
	return t.from
}

func (t Translation) ToLanguage() Language {
	return t.to
}

func (t Translation) Word() string {
	return t.word
}

func (t Translation) Meanings() []Meaning {
	return t.meanings
}

func (t Translation) Examples() []Example {
	return t.examples
}

func (t Translation) IsSuccessful() bool {
	return len(t.meanings) != 0
}
