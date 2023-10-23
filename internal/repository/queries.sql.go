// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: queries.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const insertTranslation = `-- name: InsertTranslation :exec
INSERT INTO translations (word_id, word_translation_id) VALUES ($1, $2)
`

type InsertTranslationParams struct {
	WordID            string `db:"word_id" json:"word_id"`
	WordTranslationID string `db:"word_translation_id" json:"word_translation_id"`
}

func (q *Queries) InsertTranslation(ctx context.Context, arg InsertTranslationParams) error {
	_, err := q.db.Exec(ctx, insertTranslation, arg.WordID, arg.WordTranslationID)
	return err
}

const insertWord = `-- name: InsertWord :exec
INSERT INTO words (id, "value", lang, lang_region) VALUES ($1, $2, $3, $4)
`

type InsertWordParams struct {
	ID         string      `db:"id" json:"id"`
	Value      string      `db:"value" json:"value"`
	Lang       string      `db:"lang" json:"lang"`
	LangRegion pgtype.Text `db:"lang_region" json:"lang_region"`
}

func (q *Queries) InsertWord(ctx context.Context, arg InsertWordParams) error {
	_, err := q.db.Exec(ctx, insertWord,
		arg.ID,
		arg.Value,
		arg.Lang,
		arg.LangRegion,
	)
	return err
}

const selectSimilarWords = `-- name: SelectSimilarWords :many
select value, strict_word_similarity($1, value) from words where lang = $2 and lang_region = $3 order by 2 desc limit $4
`

type SelectSimilarWordsParams struct {
	StrictWordSimilarity string      `db:"strict_word_similarity" json:"strict_word_similarity"`
	Lang                 string      `db:"lang" json:"lang"`
	LangRegion           pgtype.Text `db:"lang_region" json:"lang_region"`
	Limit                int32       `db:"limit" json:"limit"`
}

type SelectSimilarWordsRow struct {
	Value                string  `db:"value" json:"value"`
	StrictWordSimilarity float32 `db:"strict_word_similarity" json:"strict_word_similarity"`
}

func (q *Queries) SelectSimilarWords(ctx context.Context, arg SelectSimilarWordsParams) ([]SelectSimilarWordsRow, error) {
	rows, err := q.db.Query(ctx, selectSimilarWords,
		arg.StrictWordSimilarity,
		arg.Lang,
		arg.LangRegion,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SelectSimilarWordsRow{}
	for rows.Next() {
		var i SelectSimilarWordsRow
		if err := rows.Scan(&i.Value, &i.StrictWordSimilarity); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectTranslations = `-- name: SelectTranslations :many
select distinct word_tr.id, word_tr.value, word_tr.lang, word_tr.lang_region from words
join words word_tr on word_tr.lang = $1 and word_tr.lang_region = $2
join translations on (words.id = translations.word_id AND word_tr.id = translations.word_translation_id) OR (words.id = translations.word_translation_id AND word_tr.id = translations.word_id)
where words.value = $3 and words.lang = $4 and words.lang_region = $5
`

type SelectTranslationsParams struct {
	ToLang         string      `db:"to_lang" json:"to_lang"`
	ToLangRegion   pgtype.Text `db:"to_lang_region" json:"to_lang_region"`
	FromWord       string      `db:"from_word" json:"from_word"`
	FromLang       string      `db:"from_lang" json:"from_lang"`
	FromLangRegion pgtype.Text `db:"from_lang_region" json:"from_lang_region"`
}

func (q *Queries) SelectTranslations(ctx context.Context, arg SelectTranslationsParams) ([]Word, error) {
	rows, err := q.db.Query(ctx, selectTranslations,
		arg.ToLang,
		arg.ToLangRegion,
		arg.FromWord,
		arg.FromLang,
		arg.FromLangRegion,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Word{}
	for rows.Next() {
		var i Word
		if err := rows.Scan(
			&i.ID,
			&i.Value,
			&i.Lang,
			&i.LangRegion,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectWordIdByWordLangRegion = `-- name: SelectWordIdByWordLangRegion :one
SELECT id FROM words WHERE "value" = $1 AND lang = $2 AND lang_region = $3
`

type SelectWordIdByWordLangRegionParams struct {
	Value      string      `db:"value" json:"value"`
	Lang       string      `db:"lang" json:"lang"`
	LangRegion pgtype.Text `db:"lang_region" json:"lang_region"`
}

func (q *Queries) SelectWordIdByWordLangRegion(ctx context.Context, arg SelectWordIdByWordLangRegionParams) (string, error) {
	row := q.db.QueryRow(ctx, selectWordIdByWordLangRegion, arg.Value, arg.Lang, arg.LangRegion)
	var id string
	err := row.Scan(&id)
	return id, err
}
