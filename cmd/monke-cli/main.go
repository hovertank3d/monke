package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"

	"github.com/hovertank3d/monke"
	"github.com/hovertank3d/monke/client"
)

const usageFmt = `Usage: %s <COMMAND> [OPTION]...
'search' command format:
	search [--tags[TAGS]] [name]
	search without params will fetch all available titles.

'watch' command format:
	watch --link <anime link relative to https://monke.oparysh.online/anime> [[sNeN]:[sNeN]]
	watch [--tags [TAGS]] [--name [NAME]] [[sNeN]:[sNeN]]

Examples:
	monke search "oshi no ko"
	monke search --tags "ongoing -drama" "one"

	monke watch --link "kakeguruui" :s2e1
	monke watch --name "one piece" s1e200:s1e300
	monke watch --name "one punch man"
`

func usage() {
	fmt.Printf(usageFmt, os.Args[0])
	os.Exit(0)
}

func printAnime(anime monke.Anime) string {
	return fmt.Sprintf("%s(%s): %s\n", anime.Name, anime.OriginalName, anime.ID) +
		fmt.Sprintf("seasons: %d\n", len(anime.Seasons)) +
		fmt.Sprintf("description: %s\n\n", anime.Description)
}

func watch(params []string) {
	tags := ""
	name := ""
	link := ""
	selector := ""

	var anime monke.Anime

	l := len(params)
	for i := 0; i < l; i++ {
		if params[i] == "--tags" {
			if l <= i+1 {
				usage()
			}

			i += 1
			tags = params[i]
		} else if params[i] == "--name" {
			if l <= i+1 {
				usage()
			}

			i += 1
			name = params[i]
		} else if params[i] == "--link" {
			if l <= i+1 {
				usage()
			}

			i += 1
			link = params[i]
		} else {
			selector = params[i]
		}
	}

	if link == "" {
		if name == "" && tags == "" {
			usage()
		}

		animes, err := client.DefaultClient.Search(tags, name)
		if err != nil {
			log.Fatal(err)
		}

		if len(animes) == 0 {
			fmt.Println("nothing found.")
			os.Exit(0)
		}

		anime = animes[0]
	} else {
		var err error
		anime, err = client.DefaultClient.Anime(link)
		if err != nil {
			log.Fatal(err)
		}
	}

	var (
		season  = 1
		episode = 1

		lastSeason  = len(anime.Seasons) - 1
		lastEpisode = anime.Seasons[lastSeason].Episodes

		right = false
	)

	for i, r := range selector {
		if r == ':' && right {
			usage()
		}
		if r == ':' {
			right = true
			continue
		}

		if r == 's' {
			if i+1 >= len(selector) {
				usage()
			}

			idx := strings.IndexFunc(selector[i:], unicode.IsDigit)
			if idx <= 0 {
				usage()
			}

			n, _ := strconv.Atoi(selector[i+1 : i+idx+1])
			if right {
				lastSeason = n
			} else {
				season = n
			}
		}

		if r == 'e' {
			if i+1 >= len(selector) {
				usage()
			}

			idx := strings.IndexFunc(selector[i:], unicode.IsDigit)
			if idx <= 0 {
				usage()
			}

			n, _ := strconv.Atoi(selector[i+1 : i+1+idx])
			if right {
				lastEpisode = n
			} else {
				episode = n
			}
		}
	}

	if lastSeason >= len(anime.Seasons) || season >= len(anime.Seasons) {
		fmt.Println("non existent season")
		os.Exit(0)
	}

	if season > lastSeason {
		usage()
	}

	for ; season < lastSeason; season++ {
		for ; episode < anime.Seasons[season].Episodes; episode++ {
			fmt.Print(client.DefaultClient.VideoLink(anime.ID, season, episode), " ")
		}
		episode = 1
	}
	for ; episode <= lastEpisode; episode++ {
		fmt.Print(client.DefaultClient.VideoLink(anime.ID, season, episode), " ")
	}
}

func search(params []string) {
	tags := ""
	name := ""

	l := len(params)
	for i := 0; i < l; i++ {
		if params[i] == "--tags" {
			if l <= i+1 {
				usage()
			}

			i += 1
			tags = params[i]
		} else {
			name = params[0]
		}
	}

	out := ""

	animes, _ := client.DefaultClient.Search(tags, name)
	for _, a := range animes {
		out += printAnime(a)
	}

	pager := "/usr/bin/less"
	if os.Getenv("PAGER") != "" {
		pager = os.Getenv("PAGER")
	}

	cmd := exec.Command(pager)
	cmd.Stdin = strings.NewReader(out)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "watch":
		watch(os.Args[2:])
	case "search":
		search(os.Args[2:])
	default:
		usage()
	}
}
