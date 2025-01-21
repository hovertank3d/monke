package monke

import "strings"

type Anime struct {
	ID           string `gorm:"primarykey"`
	Name         string
	OriginalName string
	Tags         string
	Description  string
	Seasons      []Season `gorm:"foreignKey:AnimeID"`
}

type Season struct {
	ID        uint `gorm:"primarykey"`
	AnimeID   string
	SeasonNum int
	Episodes  int
}

func parseTags(tags string) (require []string, discard []string) {
	tags = strings.TrimSpace(tags)
	tagsComaSep := strings.ReplaceAll(tags, " ", ",")
	tokens := strings.Split(tagsComaSep, ",")

	for _, t := range tokens {
		if len(t) == 0 {
			continue
		}

		if t[0] == '-' {
			discard = append(discard, t[1:])
		} else {
			require = append(require, t)
		}
	}

	return
}
