package parser

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/hovertank3d/monke"
	"golang.org/x/text/encoding/charmap"
)

const host = "https://jut.su"

var header = http.Header{
	"User-Agent":      {"Mozilla/5.0 (X11; Linux x86_64; rv:134.0) Gecko/20100101 Firefox/134.0"},
	"Accept-Encoding": {"gzip"},
	"Content-Type":    {"application/x-www-form-urlencoded; charset=UTF-8"},
	"Origin":          {"https://jut.su"},
}

type Parser struct {
	*http.Client
}

func New() *Parser {
	return &Parser{
		Client: &http.Client{},
	}
}

func (p *Parser) get(uri string, params url.Values) (io.Reader, error) {
	req, err := http.NewRequest(http.MethodGet, host+uri+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header = header

	resp, err := p.Do(req)
	if err != nil {
		return nil, err
	}

	reader := io.Reader(resp.Body)

	enc := resp.Header.Get("Content-Encoding")
	if enc == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return charmap.Windows1251.NewDecoder().Reader(reader), nil
}

func (p *Parser) GetVid(uri string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header = header
	req.Header.Set("Content-Type", "video/mp4")
	req.Header.Set("Connection", "keep-alive")

	resp, err := p.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *Parser) post(uri string, params url.Values) (io.Reader, error) {
	body := strings.NewReader(params.Encode())
	req, err := http.NewRequest(http.MethodPost, host+uri, body)
	if err != nil {
		return nil, err
	}
	req.Header = header

	resp, err := p.Do(req)
	if err != nil {
		return nil, err
	}

	reader := io.Reader(resp.Body)

	enc := resp.Header.Get("Content-Encoding")
	if enc == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return charmap.Windows1251.NewDecoder().Reader(reader), nil
}

func (p *Parser) VideoSource(a monke.Anime, season int, episode int) (src string, _ error) {
	link := fmt.Sprintf("/%s/", a.ID)
	if len(a.Seasons) > 1 {
		link += fmt.Sprintf("season-%d/", season)
	}
	link += fmt.Sprintf("episode-%d.html", episode)

	resp, err := p.get(link, nil)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		return "", err
	}

	src = doc.Find("video > source").First().AttrOr("src", "")
	return src, nil
}

func (p *Parser) Page(pg int) ([]string, error) {
	params := url.Values{}

	params.Set("ajax_load", "yes")
	params.Set("start_from_page", strconv.Itoa(pg))
	params.Set("anime_of_user", "")

	resp, err := p.post("/anime/", params)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		return nil, err
	}

	animes := []string{}
	doc.Find("body > div > a").Each(func(i int, s *goquery.Selection) {
		animes = append(animes, s.AttrOr("href", ""))
	})

	return animes, nil
}

func (p *Parser) Anime(link string) (anime monke.Anime, _ error) {
	resp, err := p.get(link, nil)
	if err != nil {
		return anime, err
	}

	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		return anime, err
	}

	anime.Name = doc.Find("meta[itemprop*=\"name\"]").First().AttrOr("content", "")
	anime.OriginalName = doc.Find("meta[itemprop*=\"alternateName\"]").First().AttrOr("content", "")
	anime.ID = strings.Trim(link, "/")

	ssn := monke.Season{
		SeasonNum: 1,
		AnimeID:   anime.ID,
	}
	doc.Find(".short-btn").Each(func(i int, s *goquery.Selection) {
		str := s.Nodes[0].LastChild.Data

		if strings.Contains(str, "фильм") {
			return
		}

		str = strings.TrimSuffix(str, " серия")
		ep, _ := strconv.Atoi(str)
		if ep < ssn.Episodes {
			anime.Seasons = append(anime.Seasons, ssn)
			ssn.Episodes = 0
			ssn.SeasonNum++
		} else {
			ssn.Episodes = ep
		}
	})

	anime.Seasons = append(anime.Seasons, ssn)

	tags := " "
	doc.Find(".under_video_additional > a").Each(func(i int, s *goquery.Selection) {
		tag := s.AttrOr("href", "")
		tag = strings.TrimPrefix(tag, "/anime/")

		if unicode.IsDigit(rune(tag[0])) {
			tags += s.Nodes[0].LastChild.Data + " "
			return
		}

		tags += tag[:len(tag)-1] + " "
	})

	anime.Tags = strings.TrimSpace(tags)
	anime.Description = doc.Find(".under_video").Text()

	return anime, nil
}

func (p *Parser) Animes() func(func(string) bool) {
	return func(yield func(string) bool) {

	}
}
