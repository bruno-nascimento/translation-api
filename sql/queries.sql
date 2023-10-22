-- name: InsertLanguage :exec
INSERT INTO languages (id, base, region) VALUES ($1, $2, $3);

-- name: InsertWord :exec
INSERT INTO words (id, language_id, "value") VALUES ($1, $2, $3);

-- name: InsertTranslation :exec
INSERT INTO translations (word_id, word_translation_id) VALUES ($1, $2);