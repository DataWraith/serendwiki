# Serendwiki

Serendwiki is a static site generator / wiki compiler written in
[Go](https://golang.org).

It takes a folder full of markdown articles and converts them into HTML files.
Contrary to a normal wiki, where all links between articles are explicit,
serendwiki is auto-linking or serendipitous (hence the name). If you have an
article `foo.md`, then any instance of the word `foo` in any other article will
link to that article.

Currently serendwiki only supports a single flat directory, as I use `:` to
namespace my personal wiki within a single directory.


## Inspiration

Serendwiki was largely inspired by the Vim plugin
[VimBoy](https://morr.cc/keeping-a-personal-wiki/). It is the only plugin I miss
after switching to [Spacemacs](https://spacemacs.org), so I set out to replace
it.


## Quickstart

### Installation

TODO

### Creating a Wiki

- `mkdir $HOME/wiki/`
- Add some markdown files to `~/wiki/`
- `serendwiki $HOME/wiki/ $HOME/wiki-out/`

WARNING: Be careful when specifying the output directory. `serendwiki` will
overwrite files in `wiki-out` if they have the same name as a generated article.
This is useful for regenerating the wiki, but you probably should not give it
e.g. your home directory.


## Options

`serendwiki` takes several options:

| Option             | Default  | Effect                                                |
| ------------------ | :------: | ----------------------------------------------------- |
| `--min-link-len`   | 3        | Don't link short article titles (e.g. single letters) |
| `--reverse-links`  | true     | Do/do not generate a "what links here" section        |

## License

`serendwiki` is MIT-licensed. See the accompanying LICENSE file for details.
