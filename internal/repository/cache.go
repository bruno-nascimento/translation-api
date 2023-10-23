package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/bruno-nascimento/translation-api/internal/config"
)

const (
	SimilarWordsKeyPrefix              = "similar_words_"
	SelectTranslationsKeyPrefix        = "select_translations_"
	SelectWordIdByWordLangRegionPrefix = "select_word_id_by_word_lang_region_"
)

type Cache struct {
	cfg         *config.Config
	redisClient *redis.Client
	repo        Querier
	hash        *XXHashPool
}

func NewCachedRepository(cfg *config.Config, redisClient *redis.Client, repo Querier) Querier {
	return &Cache{cfg: cfg, redisClient: redisClient, repo: repo, hash: NewXXHashPool()}
}

func (c Cache) InsertTranslation(ctx context.Context, arg InsertTranslationParams) error {
	return c.repo.InsertTranslation(ctx, arg) // TODO invalidate cache
}

func (c Cache) InsertWord(ctx context.Context, arg InsertWordParams) error {
	err := c.repo.InsertWord(ctx, arg)
	if err != nil {
		return err
	}
	key, err := c.getInsertingWordKey(arg)
	if err != nil {
		log.Error().Err(err).Interface("args", arg).Msg("error getting key on InsertWord at cached repo")
		return nil
	}
	err = c.redisClient.Del(ctx, key).Err()
	if err != nil {
		log.Error().Err(err).Interface("args", arg).Msg("error deleting key on InsertWord at cached repo")
	}
	return nil
}

func (c Cache) SelectSimilarWords(ctx context.Context, arg SelectSimilarWordsParams) ([]SelectSimilarWordsRow, error) {
	key, err := c.getSimilarWordsKey(arg)
	if err != nil {
		return c.repo.SelectSimilarWords(ctx, arg)
	}

	val, err := c.redisClient.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		log.Error().Err(err).Str("key", key).Interface("arg", arg).Msg("error fetching key from redis on SelectSimilarWords at cached repo")
		return c.repo.SelectSimilarWords(ctx, arg)
	}

	var res []SelectSimilarWordsRow
	if val == "" {
		res, err = c.repo.SelectSimilarWords(ctx, arg)
		if err != nil {
			return res, err
		}

		bytes, err := json.Marshal(res)
		if err != nil {
			return res, err
		}

		err = c.redisClient.Set(ctx, key, bytes, c.cfg.Cache.TTL/2).Err() // we set the TTL to half of the original value to avoid stale data since we can't invalidate the cache on insert
		if err != nil {
			log.Error().Err(err).Interface("key", arg).RawJSON("value", bytes).Msg("error setting key on SelectSimilarWords at cached repo")
		}
		return res, err
	}
	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return c.repo.SelectSimilarWords(ctx, arg)
	}
	return res, nil

}

func (c Cache) SelectTranslations(ctx context.Context, arg SelectTranslationsParams) ([]Word, error) {
	key, err := c.getSelectTranslationsKey(arg)
	if err != nil {
		log.Error().Err(err).Interface("arg", arg).Msg("error getting key on SelectTranslations at cached repo")
		return c.repo.SelectTranslations(ctx, arg)
	}

	val, err := c.redisClient.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		log.Error().Err(err).Str("key", key).Interface("arg", arg).Msg("error fetching key from redis on SelectTranslations at cached repo")
		return c.repo.SelectTranslations(ctx, arg)
	}

	var res []Word
	if val == "" {
		res, err = c.repo.SelectTranslations(ctx, arg)
		if err != nil {
			return res, err
		}

		bytes, err := json.Marshal(res)
		if err != nil {
			return res, err
		}

		err = c.redisClient.Set(ctx, key, bytes, c.cfg.Cache.TTL).Err()
		if err != nil {
			log.Error().Err(err).Interface("key", arg).RawJSON("value", bytes).Msg("error setting key on SelectTranslations at cached repo")
		}
		return res, err
	}
	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return c.repo.SelectTranslations(ctx, arg)
	}
	return res, nil
}

func (c Cache) SelectWordIdByWordLangRegion(ctx context.Context, arg SelectWordIdByWordLangRegionParams) (string, error) {
	key, err := c.getSelectWordIdByWordLangRegionKey(arg)
	if err != nil {
		return c.repo.SelectWordIdByWordLangRegion(ctx, arg)
	}

	val, err := c.redisClient.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return c.repo.SelectWordIdByWordLangRegion(ctx, arg)
	}

	if val == "" {
		var res string
		res, err = c.repo.SelectWordIdByWordLangRegion(ctx, arg)
		if err != nil {
			return res, err
		}

		err = c.redisClient.Set(ctx, key, res, c.cfg.Cache.TTL).Err()
		if err != nil {
			log.Error().Err(err).Interface("key", arg).Str("value", res).Msg("error setting key on SelectWordIdByWordLangRegion at cached repo")
		}
		return res, err
	}

	return val, nil
}

func (c Cache) getKey(prefix string, args ...any) (string, error) {
	hash := c.hash.Acquire()
	defer c.hash.Release(hash)
	for _, arg := range args {
		_, err := hash.WriteString(fmt.Sprint(arg))
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s%d", prefix, hash.Sum64()), nil
}

func (c Cache) getSimilarWordsKey(arg SelectSimilarWordsParams) (string, error) {
	return c.getKey(SimilarWordsKeyPrefix, []any{arg.StrictWordSimilarity, arg.Lang, arg.LangRegion, arg.Limit}...)
}

func (c Cache) getSelectTranslationsKey(arg SelectTranslationsParams) (string, error) {
	return c.getKey(SelectTranslationsKeyPrefix, []any{arg.FromWord, arg.FromLang, arg.FromLangRegion, arg.ToLang, arg.ToLangRegion}...)
}

func (c Cache) getSelectWordIdByWordLangRegionKey(arg SelectWordIdByWordLangRegionParams) (string, error) {
	return c.getKey(SelectWordIdByWordLangRegionPrefix, []any{arg.Lang, arg.LangRegion, arg.Value}...)
}

func (c Cache) getInsertingWordKey(arg InsertWordParams) (string, error) {
	return c.getKey(SelectWordIdByWordLangRegionPrefix, []any{arg.Lang, arg.LangRegion, arg.Value}...)
}

type XXHashPool struct {
	pool *sync.Pool
}

func NewXXHashPool() *XXHashPool {
	return &XXHashPool{
		pool: &sync.Pool{
			New: func() any {
				return xxhash.New()
			},
		},
	}
}

func (p *XXHashPool) Acquire() *xxhash.Digest {
	return p.pool.Get().(*xxhash.Digest)
}

func (p *XXHashPool) Release(d *xxhash.Digest) {
	d.Reset()
	p.pool.Put(d)
}
