package internal

import (
	t "DeutschBot/package/translator"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type translationCache struct {
	gorm.Model
	Word     string
	FromLang string
	ToLang   string
	Cache    datatypes.JSON
}

type meaningDto struct {
	Word   string `json:"Word"`
	Gender string `json:"Gender"`
}

type exampleDto struct {
	Usage   string `json:"Usage"`
	Meaning string `json:"Meaning"`
}

type translationDto struct {
	Meanings []meaningDto `json:"Meanings"`
	Examples []exampleDto `json:"Examples"`
}

func newTranslationCache(translation t.Translation) translationCache {
	payload := translationDto{}

	for _, meaning := range translation.Meanings() {
		payload.Meanings = append(payload.Meanings, meaningDto{
			Word:   meaning.Word(),
			Gender: meaning.Gender(),
		})
	}

	for _, example := range translation.Examples() {
		payload.Examples = append(payload.Examples, exampleDto{
			Usage:   example.Usage(),
			Meaning: example.Meaning(),
		})
	}

	cache := translationCache{
		Word:     translation.Word(),
		FromLang: translation.FromLanguage().String(),
		ToLang:   translation.ToLanguage().String(),
		Cache:    Serialize(payload),
	}

	return cache
}

func (tc translationCache) extract() t.Translation {
	var payload translationDto
	Deserialize(tc.Cache, &payload)

	var meanings []t.Meaning
	for _, meaning := range payload.Meanings {
		meanings = append(meanings, t.NewMeaning(meaning.Word, meaning.Gender))
	}

	var examples []t.Example
	for _, example := range payload.Examples {
		examples = append(examples, t.NewExample(example.Usage, example.Meaning))
	}

	return t.NewTranslation(tc.Word, t.ToLang(tc.FromLang), t.ToLang(tc.ToLang), meanings, examples)
}

type cachingTranslator struct {
	cache      *gorm.DB
	translator t.Translator
}

func initTranslator(cacheSource *gorm.DB) cachingTranslator {
	cacheSource.AutoMigrate(&translationCache{})

	return cachingTranslator{
		cache:      cacheSource,
		translator: t.NewReverso(3),
	}
}

func (ct cachingTranslator) Translate(word string, fromLang t.Language, toLang t.Language) (t.Translation, error) {
	var cache translationCache
	ct.cache.
		Model(&translationCache{}).
		Where("word = ? AND from_lang = ? AND to_lang = ?", word, fromLang.String(), toLang.String()).
		Find(&cache)

	// If cache does not exist
	if cache.ID == 0 {
		result, err := ct.translator.Translate(word, fromLang, toLang)

		if err == nil {
			cache = newTranslationCache(result)
			ct.cache.Create(&cache)
		}

		return result, err
	}

	return cache.extract(), nil
}

var (
	Translator = initTranslator(Db)
)
