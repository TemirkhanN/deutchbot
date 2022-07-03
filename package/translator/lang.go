package translator

type Language struct {
	value string
}

func (l Language) String() string {
	return l.value
}

var (
	EN = Language{value: "english"}
	DE = Language{value: "german"}
)
