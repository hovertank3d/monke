package server

import (
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/hovertank3d/monke"
	"github.com/hovertank3d/monke/env"
	"github.com/hovertank3d/monke/parser"
	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo
	db *monke.DB
	p  *parser.Parser
}

//go:embed player.html
var playerPage []byte

func New(db *monke.DB) *Server {
	s := &Server{
		Echo: echo.New(),
		db:   db,
		p:    parser.New(),
	}

	s.db.AutoMigrate(&monke.Anime{})
	s.db.AutoMigrate(&monke.Season{})

	s.GET("/proxy/:id/:season/:episode", s.proxyVideo)

	s.GET("/api/:id", s.apiAnime)
	s.GET("/api/search", s.apiSearch)

	if env.MonkePlayer {
		s.GET("/", func(c echo.Context) error {
			return c.Stream(200, "text/html", bytes.NewBuffer(playerPage))
		})
	}

	return s
}

func (s *Server) Start(addr string) error {
	return s.Echo.Start(addr)
}

func (s *Server) proxyVideo(c echo.Context) error {
	fmt.Println(c.Path())

	id := strings.Trim(c.Param("id"), "/")
	seasonStr := strings.Trim(c.Param("season"), "/")
	episodeStr := strings.Trim(c.Param("episode"), "/")

	episode, err := strconv.Atoi(strings.TrimSuffix(episodeStr, ".mp4"))
	if err != nil || episode == 0 {
		episode = 1
	}
	season, err := strconv.Atoi(seasonStr)
	if err != nil || season == 0 {
		season = 1
	}

	var anime monke.Anime
	s.db.Model(&monke.Anime{}).
		Where("id = $1", id).
		Preload("Seasons").
		First(&anime)

	uri, err := s.p.VideoSource(anime, season, episode)
	if err != nil {
		fmt.Println(err)
		return err
	}

	reader, err := s.p.GetVid(uri)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if reader.Header != nil {
		for k, v := range reader.Header {
			if len(v) == 0 {
				continue
			}
			c.Response().Header().Set(k, v[0])
		}
	}

	return c.Stream(reader.StatusCode, "video/mp4", reader.Body)
}
