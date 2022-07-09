package translator

type Language struct {
	value string
}

func (l Language) String() string {
	return l.value
}

func ToLang(lang string) Language {
	return Language{value: lang}
}

var (
	EN = ToLang("english")
	DE = ToLang("german")
)
