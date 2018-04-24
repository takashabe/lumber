# lumber

*This project is currently Alpha version*

Lumber is the blog management system. Lumber has an API server to return blog content and CLI tool to manage blog content.

## Installation

```
go get -u github.com/takashabe/lumber
```

## CLI

### Post entry

- single file

```
client post -addr="YOUR_LUMBER_SERVER_ADDR" -token="YOUR_LUMBER_SERVER_TOKEN" -file=path/to/file.md
```

- multi files in directory

```
client post-dir -addr="YOUR_LUMBER_SERVER_ADDR" -token="YOUR_LUMBER_SERVER_TOKEN" -dir=path/to/dir
```

## REST API

REST API to backend of the `lumber-web` frontend and lumber CLI tool.

### Entry

| Method                | URL                                  | Behavior                                                |
| ------                | ------                               | -----                                                   |
| Get entry             | GET:    `/api/entry/:id`             | Get detail a the entry                                  |
| Get list entry ids    | GET:    `/api/entries`               | Get all the entry ids                                   |
| Get list entry titles | GET:    `/api/titles/:start/:length` | Get the ":length" numbers entry titles from ":start" id |
| Post entry            | POST:    `/api/entry`                | Post the entry                                          |
| Edit entry            | PUT:    `/api/entry/:id`             | Edit the entry                                          |
| Delete entry          | DELETE:   `/api/entry/:id`           | Delete the entry                                        |

