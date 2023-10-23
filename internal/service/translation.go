package service

import (
	"context"

	"github.com/bruno-nascimento/translation-api/internal/config"
	"github.com/bruno-nascimento/translation-api/internal/entrypoint/http/api"
	"github.com/bruno-nascimento/translation-api/internal/repository"
)

type Translation interface {
	CreateTranslation(ctx context.Context, newMigration *NewTranslationParam) error
	FindTranslations(ctx context.Context, params *FindTranslationParam) (*TranslationResult, error)
}

type translation struct {
	repo repository.Querier
	cfg  *config.Config
}

func NewTranslation(cfg *config.Config, repo repository.Querier) Translation {
	return &translation{
		repo: repo,
		cfg:  cfg,
	}
}

func (t translation) CreateTranslation(ctx context.Context, newMigration *NewTranslationParam) error {
	wordParams := newMigration.ToInsertWordParam()
	for idx, wordParam := range wordParams {
		err := t.repo.InsertWord(ctx, wordParam)
		if err != nil {
			if repository.IsUniqueViolation(err) {
				wordParams[idx].ID, err = t.repo.SelectWordIdByWordLangRegion(ctx, newMigration.ToSelectWordParams(WordPairType(idx)))
				if err != nil {
					return err
				}
				continue
			}
			return err
		}
	}

	err := t.repo.InsertTranslation(ctx, repository.InsertTranslationParams{
		WordID:            wordParams[From.Idx()].ID,
		WordTranslationID: wordParams[To.Idx()].ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (t translation) FindTranslations(ctx context.Context, params *FindTranslationParam) (*TranslationResult, error) {
	var result []string
	translations, err := t.repo.SelectTranslations(ctx, params.ToSelectTranslationParams())
	if err != nil {
		return nil, err
	}
	if len(translations) > 0 {
		for _, tr := range translations {
			result = append(result, tr.Value)
		}
		return &TranslationResult{Result: &api.Translation{Results: &result}}, nil
	}
	if !t.cfg.Suggestions.Enabled {
		return &TranslationResult{}, nil
	}
	words, err := t.repo.SelectSimilarWords(ctx, params.ToSelectSimilarWordsParams(t.cfg))
	if err != nil {
		return nil, err
	}
	for _, wo := range words {
		if wo.StrictWordSimilarity > 0 {
			result = append(result, wo.Value)
		}
	}
	if len(result) > 0 {
		return &TranslationResult{Suggestion: &api.TranslationSuggestions{SimilarWords: &result}}, nil
	}
	return &TranslationResult{}, nil
}
