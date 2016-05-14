# Serendwiki

Serendwiki is a static site generator / wiki compiler written in
[Go](https://golang.org).

It takes a folder full of markdown articles (or symlinks to them) and converts
them into HTML files. Contrary to a normal wiki, where all links between
articles are explicit, serendwiki is auto-linking or serendipitous (hence the
name). If you have an article file `foo` (no extension), then any instance of
the word `foo` in any other article will link to that article.

Currently serendwiki only supports a single flat directory, as I use `:` to
namespace my personal wiki within a single directory. Non-hidden files and
directories that serendwiki does not recognize as wiki files will be copied to
the output directory verbatim, so that you can, for example, link image files
within your wiki.


## Inspiration

Serendwiki was inspired by the Vim plugin
[VimBoy](https://morr.cc/keeping-a-personal-wiki/).


## Quickstart

### Installation

- Install [golang](https://golang.org)
- `go get -u -v github.com/DataWraith/serendwiki`

### Creating a Wiki

- `mkdir $HOME/wiki/`
- Add some markdown files to `~/wiki/` -- the files should have no extension
- `serendwiki $HOME/wiki/ $HOME/wiki-out/`

NOTE: `serendwiki` will refuse to run if the `wiki-out` directory already exists
in order to avoid overwriting existing files.


## License

`serendwiki` is MIT-licensed. See the accompanying LICENSE file for details.
