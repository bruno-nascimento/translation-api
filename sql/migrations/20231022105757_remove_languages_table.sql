-- migrate:up
ALTER TABLE words ADD COLUMN lang varchar(4);
ALTER TABLE words ADD COLUMN lang_region varchar(4);

UPDATE words SET lang = langs.base, lang_region = langs.region
FROM (SELECT id, base, region FROM languages) AS langs
WHERE words.language_id=langs.id;

ALTER TABLE words DROP COLUMN language_id;
ALTER TABLE words ALTER COLUMN lang SET NOT NULL;

ALTER TABLE words ADD CONSTRAINT word_lang_region_idx_unique UNIQUE (value, lang, lang_region);

drop table if exists languages;

-- migrate:down
ALTER TABLE words DROP COLUMN lang;
ALTER TABLE words DROP COLUMN lang_region;

create table if not exists languages (
    id varchar(26) not null primary key,
    base char(2) not null,
    region char(2) not null,
    constraint languages_base_region_unique unique (base, region)
);

ALTER TABLE words ADD COLUMN language_id varchar(26) CONSTRAINT languages_fk REFERENCES languages (id) -- NOT NULL;
