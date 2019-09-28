package scrape

// IGrab ...
type IGrab interface {
	MainPage(url string)
	Sample(bool)
	Name() string
	Find(string) (IGrab, error)
	Decode(*[]*Message) error
}

// Sample ...
type Sample struct {
	Index int
	Thumb string
	Image string
	Title string
}

// GrabLanguage ...
type GrabLanguage int

// GrabLanguage detail ...
const (
	LanguageEnglish GrabLanguage = iota
	LanguageJapanese
	LanguageChinese
	LanguageKorea
)
