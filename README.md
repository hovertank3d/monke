# monke

Simple [jut.su](https://jut.su) parser + API + proxy + cli client

## api and proxy endpoints
search animes by their tags and name

query params:
 * `tags` -- list of space separated tags. minus prefixed are treated as now wanted;
 * `name` -- original or translated anime name. 

`/api/search?tags=<tags>&name=<name>`

## deploying

TODO explain and add more info 

```bash
docker compose up -d
docker compose exec -it app /build/monke rescan 40
```
