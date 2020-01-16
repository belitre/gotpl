# gotpl - CLI tool for Golang templates

Command line tool that compiles Golang
[templates](http://golang.org/pkg/text/template/) with values from YAML and JSON files.

Supports [Sprig functions](https://github.com/Masterminds/sprig).

Inspired by [Helm](https://github.com/kubernetes/helm)

Started as a fork from [github.com/tsg/gotpl](https://github.com/tsg/gotpl) but this project looks abandoned

Added some things from [Helm](https://github.com/kubernetes/helm), like the way they read values, and some functions for the templates.

And some improvements by myself :)

## Install

You can download gotpl binaries for windows, linux and mac from here: https://github.com/belitre/gotpl/releases

## Usage

Say you have a `template.tpl` file like this:

```
    {{.first_name}} {{.last_name}} is {{.age}} years old.
```

and a `user.yml` YAML file like this one:

```
    first_name: Max
    last_name: Mustermann
    age: 30
```

You can compile the template like this:

```
    gotpl template.tpl -f user.yml
```

You can compile multiple templates at the same time like this (__warning: gotpl will generate a single ouput for all the templates!__):

```
    gotpl template.tpl other_template.tpl -f user.yml
```

You can compile templates and directories, __includes subdirectories__ (__warning: gotpl will generate a single ouput for all the templates!__):

```
    gotpl template.tpl /templates/directory -f user.yml
```

You can set values though the command line too:

```
    gotpl template.tpl -f user.yaml --set age=40
```

You can set an output folder:

```
    gotpl template.tpl /templates/directory -f user.yml --set age=40 -o /tplresult
```

__When using `-o`gotpl will generate files with the same names as the templates, so if the template file is `template.tpl` gotpl will generate `$OUPUT_FOLDER/template.tpl`. If multiple files with the same name are provided then gotpl will override these files and will keep only the last one!__

You can get help about how to use the command running:

```
    gotpl -h
```

## Development

* Requires `go 1.13.1``
