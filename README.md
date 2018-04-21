# gotpl - CLI tool for Golang templates

Command line tool that compiles Golang
[templates](http://golang.org/pkg/text/template/) with values from YAML and JSON files.

Supports [Sprig functions](https://github.com/Masterminds/sprig).

Inspired by [Helm](https://github.com/kubernetes/helm)

Started as a fork from [github.com/tsg/gotpl](https://github.com/tsg/gotpl) but this project looks abandoned

Added some things from [Helm](https://github.com/kubernetes/helm), like the way they read values.

And some improvements by myself :)

## Install

    go get github.com/belitre/gotpl

## Usage

Say you have a `template` file like this:

    {{.first_name}} {{.last_name}} is {{.age}} years old.

and a `user.yml` YAML file like this one:

    first_name: Max
    last_name: Mustermann
    age: 30

You can compile the template like this:

    gotpl template -f user.yml

You can set values though the command line too:

    gotpl template -f user.yam --set age=40

You can get help about how to use the command running:

    gotpl -h

## Development

* Requires `go 1.10.1``
* [godep](https://github.com/tools/godep) for dependency management