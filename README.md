Web to Markdown
===============

A recursive, parallel webscraper and markdown converter written as a beginner project to get a feeling for the go programming language, it's take on parallelization and testing.

Usage
=====

```
usage: web_to_md <root_url> <storage_directory>
```

recursively scrapes the `root_url` for any link that leads further down, gets it and saves the body as a markdown file in a file hierarchy under `storage_directory` that can be used as a basis for the `content` folder of a [hugo](https://gohugo.io) project.


Notes
=====

* `go list -f '{{ .Imports }}'` vs. `go list -f '{{ .TestImports }}'` gives me confidence that I cant import whatever I want during testing without cluttering the final binary. [This](https://dave.cheney.net/2014/09/14/go-list-your-swiss-army-knife) shows you more nice things you can do with `go list`.

* `Interfaces` is a very interesting and ergonomic approach to abstractions. It was used mainly to ease the task of testing.
