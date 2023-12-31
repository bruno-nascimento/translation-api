// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package api

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error message
	Message string `json:"message"`
}

// NewTranslation defines model for NewTranslation.
type NewTranslation struct {
	// From word and language of the word to be translated
	From *struct {
		// Language The language of the word to be translated
		Language *string `json:"language,omitempty"`

		// Word The word to be translated
		Word *string `json:"word,omitempty"`
	} `json:"from,omitempty"`

	// To word and language pair containing translation of the 'from' word
	To *struct {
		// Language The language of the word to be translated
		Language *string `json:"language,omitempty"`

		// Word The word to be translated
		Word *string `json:"word,omitempty"`
	} `json:"to,omitempty"`
}

// Translation defines model for Translation.
type Translation struct {
	Results *[]string `json:"results,omitempty"`
}

// TranslationSuggestions A list of suggestions containing words that are similar to the one provided in the 'word' parameter with the same 'language' parameter
type TranslationSuggestions struct {
	SimilarWords *[]string `json:"similar_words,omitempty"`
}

// FindTranslationParams defines parameters for FindTranslation.
type FindTranslationParams struct {
	// Word Word to be translated
	Word string `form:"word" json:"word"`

	// Language The actual language of the word we want to translate
	Language string `form:"language" json:"language"`

	// TargetLanguage The target language for translation of the 'word' parameter
	TargetLanguage string `form:"target_language" json:"target_language"`
}

// AddTranslationJSONRequestBody defines body for AddTranslation for application/json ContentType.
type AddTranslationJSONRequestBody = NewTranslation
