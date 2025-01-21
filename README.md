# monke

Simple [jut.su](https://jut.su) parser + API + proxy + cli client

## api and proxy endpoints
search animes by their tags and name

query params:
 * `tags` -- list of space separated tags. minus prefixed are treated as not wanted;
 * `name` -- original or translated anime name. 

`/api/search?tags=<tags>&name=<name>`

`/proxy/<anime_id>/<season>/<episode>`

## ./cmd/cli

simple command-line search and link generation utility.
usage examples:
```bash
monke-cli search "jojo"
monke-cli search --tags "-2025 -ongoing drama comedy"

# pass links of episodes to mpv
mpv $(monke-cli watch --name "one piece")
mpv $(monke-cli watch --link kakeguruui s1e5:s2e2)
```

## deploying

TODO explain and add more info 

```bash
docker compose up -d
docker compose exec -it app /build/monke rescan 40
```
