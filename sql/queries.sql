-- name: InsertWord :exec
INSERT INTO words (id, "value", lang, lang_region) VALUES ($1, $2, $3, $4);

-- name: SelectWordIdByWordLangRegion :one
SELECT id FROM words WHERE "value" = $1 AND lang = $2 AND lang_region = $3;

-- name: InsertTranslation :exec
INSERT INTO translations (word_id, word_translation_id) VALUES ($1, $2);

-- name: SelectTranslations :many
select distinct word_tr.* from words
join words word_tr on word_tr.lang = @to_lang and word_tr.lang_region = @to_lang_region
join translations on (words.id = translations.word_id AND word_tr.id = translations.word_translation_id) OR (words.id = translations.word_translation_id AND word_tr.id = translations.word_id)
where words.value = @from_word and words.lang = @from_lang and words.lang_region = @from_lang_region;

-- name: SelectSimilarWords :many
select value, strict_word_similarity($1, value) from words where lang = $2 and lang_region = $3 order by 2 desc limit $4;