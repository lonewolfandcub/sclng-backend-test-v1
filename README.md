# Canvas for Backend Technical Test at Scalingo

## Execution

```
docker compose up
```

Application will be then running on port `5000`

## Test web-application is alive

```
$ curl localhost:5000/ping
{ "status": "pong" }
```

## Usage

### List last 100 public repositories
```
$ curl localhost:5000/repos
[
    {
        "created_at": "2023-12-03T13:10:50Z",
        "full_name": "fajarhidayatt/foodmarket",
        "id": 726806148,
        "language": "Java",
        "license": {
            "name": ""
        },
        "url": "https://api.github.com/repos/fajarhidayatt/foodmarket"
    },
    {
        "created_at": "2023-12-04T00:39:15Z",
        "full_name": "leofigue/ReactContext",
        "id": 726995192,
        "language": "JavaScript",
        "license": {
            "name": ""
        },
        "url": "https://api.github.com/repos/leofigue/ReactContext"
    },
    ...
]
```

### List last 100 public repositories with all languages
```
$ curl localhost:5000/stats
[
    {
        "full_name": "kobili/go-migration-tool",
        "id": 683900069,
        "languages": [
            {
                "bytes": 1139,
                "name": "Go"
            },
            {
                "bytes": 159,
                "name": "Makefile"
            }
        ],
        "license": {
            "name": ""
        },
        "url": "https://api.github.com/repos/kobili/go-migration-tool"
    },
    {
        "full_name": "Carcharodon1503/html-portfolio",
        "id": 724350283,
        "languages": [
            {
                "bytes": 4023,
                "name": "HTML"
            }
        ],
        "license": {
            "name": ""
        },
        "url": "https://api.github.com/repos/Carcharodon1503/html-portfolio"
    },
    ...
]
```

### Filters

Both API accept `language` and `license` parameters as filters.

```
$ curl localhost:5000/repos?language=go
$ curl localhost:5000/stats?license=gpl
$ curl 'localhost:5000/repos?language=go,cpp,rust&license=mit'
```

## Architecture

The wweb-application is composed of two packages:
* `main` contains the `main.go` file that launches the webserver and handles the URL calls from the users.
* `github` contains the Github client that handles the calls to github HTTP API, and the JSON models used to decode Github API reponses and for the web-application own responses.

![General diagram](/doc/sclng-diagram.png)


### Files

* main.go
* gihub/client.go
* gihub/models.go



### Github client  (client.go)

The main struture is `github.Client` and its two main public methods:
1. `ListLatestRepositories` to list the latest repositories,
2. `GatherLatestRepositoriesStats` to list the latest repositories and agregate the languages informations.


