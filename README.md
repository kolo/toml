## Overview

toml.go is implementation of parser for configuration files created by using
[TOML](https://github.com/mojombo/toml). It supports all types of data
except of [array of tables](https://github.com/mojombo/toml#array-of-tables).

## Installation

To install toml.go package run `go get github.com/kolo/toml.go`. To use
it in your application add `github.com/kolo/toml.go` string to `import`
statement.

## Usage

    // config.toml
    // [package]
    //   name = "toml.go"
    //   authors = ["Dmitry Maksimov"]

    conf, err := toml.Parse("config.toml")
    if err != nil {
        // handle error
    }

    fmt.Println(conf.String("package.name"))
    authors := conf.Slice("package.authors")
    for _, author := range authors {
        fmt.Println(author)
    }

Keep in mind that getters return default values of requested type if none was
found.

## Contribution

Feel free to fork the project, submit pull requests, ask questions.

## Authors

Dmitry Maksimov (dmtmax@gmail.com)
