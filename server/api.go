package server

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *Server) apiSearch(c echo.Context) error {
	tags := c.QueryParam("tags")
	name := c.QueryParam("name")

	animes, err := s.db.SearchAnime(name, tags)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(200, animes)
}

func (s *Server) apiAnime(c echo.Context) error {
	id := c.Param("id")

	anime, err := s.db.AnimeByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(200, anime)
}
