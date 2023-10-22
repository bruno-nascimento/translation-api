-- migrate:up
CREATE EXTENSION btree_gist;

create table languages (
    id varchar(26) not null primary key,
    base char(2) not null,
    region char(2) not null,
    constraint languages_base_region_unique unique (base, region)
);

create index languages_base_index on languages (base);

create table words (
    id varchar(26) not null primary key,
    language_id varchar(26) not null constraint words_language_id_foreign references languages,
    value varchar(255) not null,
    constraint words_language_id_value_unique unique (language_id, value)
);

create index words_value_index on words using GIST(value);

create table translations (
    word_id varchar(26) not null constraint translations_word_id_foreign references words,
    word_translation_id varchar(26) not null constraint translations_word_translation_id_foreign references words,
    primary key (word_id, word_translation_id)
);

-- migrate:down
drop table if exists translations;
drop table if exists words;
drop table if exists languages;
