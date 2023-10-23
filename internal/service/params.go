package service

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oklog/ulid/v2"
	"golang.org/x/text/language"

	"github.com/bruno-nascimento/translation-api/internal/config"
	"github.com/bruno-nascimento/translation-api/internal/entrypoint/http/api"
	"github.com/bruno-nascimento/translation-api/internal/repository"
)

type WordPairType int

func (w WordPairType) Idx() int {
	return int(w)
}

const From, To = WordPairType(0), WordPairType(1)

type NewTranslationParam struct {
	From Pair
	To   Pair
}

type Pair struct {
	ID       string
	Language language.Tag
	Word     string
}

func (p Pair) ValidateLanguage() error {
	_, confidence := p.Language.Base()
	if confidence != language.Exact {
		return fmt.Errorf("from language is not valid: %s", p.Language.String())
	}
	return nil
}

func NewTranslationParamFromAPI(newTranslation api.NewTranslation) (*NewTranslationParam, error) {
	fromLang, err := language.Parse(*newTranslation.From.Language)
	if err != nil {
		return nil, err
	}

	toLang, err := language.Parse(*newTranslation.To.Language)
	if err != nil {
		return nil, err
	}

	return &NewTranslationParam{
		From: Pair{
			ID:       ulid.Make().String(),
			Language: fromLang,
			Word:     *newTranslation.From.Word,
		},
		To: Pair{
			ID:       ulid.Make().String(),
			Language: toLang,
			Word:     *newTranslation.To.Word,
		},
	}, nil
}

func (n NewTranslationParam) ValidateLanguages() error {
	var err error
	if e := n.From.ValidateLanguage(); e != nil {
		err = e
	}
	if e := n.To.ValidateLanguage(); e != nil {
		err = fmt.Errorf("to language is not valid: %s; %w", n.To.Language.String(), err)
	}
	return err
}

func (n NewTranslationParam) ToInsertWordParam() [2]repository.InsertWordParams {
	return [2]repository.InsertWordParams{
		{
			ID:         ulid.Make().String(),
			Value:      strings.ToLower(n.From.Word),
			Lang:       getLang(n.From.Language),
			LangRegion: getRegion(n.From.Language),
		},
		{
			ID:         ulid.Make().String(),
			Value:      strings.ToLower(n.To.Word),
			Lang:       getLang(n.To.Language),
			LangRegion: getRegion(n.To.Language),
		},
	}
}

func (n NewTranslationParam) ToSelectWordParams(fromOrTo WordPairType) repository.SelectWordIdByWordLangRegionParams {
	if fromOrTo == From {
		return repository.SelectWordIdByWordLangRegionParams{
			Value:      n.From.Word,
			Lang:       getLang(n.From.Language),
			LangRegion: getRegion(n.From.Language),
		}
	}
	return repository.SelectWordIdByWordLangRegionParams{
		Value:      n.To.Word,
		Lang:       getLang(n.To.Language),
		LangRegion: getRegion(n.To.Language),
	}
}

func getLang(lang language.Tag) string {
	base, _ := lang.Base()
	return base.String()
}

func getRegion(lang language.Tag) pgtype.Text {
	region, confidence := lang.Region()
	if confidence != language.Exact {
		return pgtype.Text{
			String: "",
			Valid:  false,
		}
	}
	return pgtype.Text{
		String: region.String(),
		Valid:  true,
	}
}

type FindTranslationParam struct {
	From Pair
	To   Pair
}

func (p FindTranslationParam) ToSelectTranslationParams() repository.SelectTranslationsParams {
	return repository.SelectTranslationsParams{
		FromWord:       p.From.Word,
		FromLang:       getLang(p.From.Language),
		FromLangRegion: getRegion(p.From.Language),
		ToLang:         getLang(p.To.Language),
		ToLangRegion:   getRegion(p.To.Language),
	}
}

func (p FindTranslationParam) ToSelectSimilarWordsParams(cfg *config.Config) repository.SelectSimilarWordsParams {
	return repository.SelectSimilarWordsParams{
		StrictWordSimilarity: p.From.Word,
		Lang:                 getLang(p.From.Language),
		LangRegion:           getRegion(p.From.Language),
		Limit:                int32(cfg.Suggestions.Limit),
	}
}

func NewFindTranslationParamFromAPI(params api.FindTranslationParams) (*FindTranslationParam, error) {
	fromLang, err := language.Parse(params.Language)
	if err != nil {
		return nil, err
	}

	toLang, err := language.Parse(params.TargetLanguage)
	if err != nil {
		return nil, err
	}

	return &FindTranslationParam{
		From: Pair{
			Language: fromLang,
			Word:     params.Word,
		},
		To: Pair{
			Language: toLang,
			Word:     "",
		},
	}, nil
}

type TranslationResult struct {
	Result     *api.Translation
	Suggestion *api.TranslationSuggestions
}
