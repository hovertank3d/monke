package monke

import (
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func (db *DB) AnimeByID(id string) (anime Anime, _ error) {
	return anime, db.Model(&Anime{}).
		Where("id = $1", id).
		Preload("Seasons").
		First(&anime).Error
}

func (db *DB) SearchAnime(name string, tags string) ([]Anime, error) {
	whereQuery := `
NOT EXISTS (
    SELECT 1
    FROM unnest($1::text[]) AS accept_tag
    WHERE accept_tag NOT IN (SELECT unnest(string_to_array(TRIM(BOTH FROM tags), ' ')))
)
AND NOT EXISTS (
	SELECT 1
    FROM unnest($2::text[]) AS discard_tag
    WHERE discard_tag IN (SELECT unnest(string_to_array(TRIM(BOTH FROM tags), ' ')))
)
AND (LOWER(name) SIMILAR TO LOWER ($3) OR LOWER(original_name) SIMILAR TO LOWER($3))`

	var (
		animes []Anime
	)

	require, discard := parseTags(tags)
	err := db.Model(&Anime{}).
		Where(whereQuery, require, discard, "%"+name+"%").
		Order("name").
		Preload("Seasons").
		Find(&animes).Error

	return animes, err
}
