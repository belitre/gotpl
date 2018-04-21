# gotpl - CLI tool for Golang templates

Forked from github.com/tsg/gotpl 

Added the content of pull requests:
 * https://github.com/tsg/gotpl/pull/5
 * https://github.com/tsg/gotpl/pull/7

And some improvements by myself :)

Command line tool that compiles Golang
[templates](http://golang.org/pkg/text/template/) with values from YAML and JSON files.

Supports [Sprig functions](https://github.com/Masterminds/sprig).

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

You can get help about how to use the command running:

    gotpl -h