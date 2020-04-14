# Kwik

Kwik is a tiny and fast wiki engine for Markdown formatted documents, intended for personal use. It allows searching for strings and documents, and uses the filesystem as storage.

## Motivation

I was using my own Rails wiki [kwik-ruby](https://github.com/dncrht/kwik-ruby) but felt slow. I started to be interested in Golang, therefore decided to rewrite it in Go as a learning exercise. Learning by example is the best!

## Installation

I've prepared a convenient installation script. It is documented step by step:

    ./install

_NOTE: successfully tested on MacOS 10.14 with Go 1.14_

At this point you'll have a kwik server on [port 2005](http://localhost:2005) that starts up with your system.

If you want to uninstall it, just run:

    ./uninstall

## Usage

Point your browser to [http://localhost:2005](http://localhost:2005). Create a page. Preview the content. Save it. Search for it.

All your pages will be stored as files in the `pages` directory.

## Development

Compile and run locally:

    ./build_n_serve

## Contributing

Still learning Go, so I haven't written tests. Shame!

Bug reports and pull requests are welcome on GitHub at https://github.com/dncrht/kwik.

## License

This codebase is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
