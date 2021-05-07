# codekeeper : (ckp)
Simple CLI that calls, keeps and fetches your useful scripts from a git repository

## Overview [![GoDoc](https://godoc.org/github.com/elhmn/codekeeper?status.svg)](https://godoc.org/github.com/elhmn/ckp)

The codekeeper (ckp) CLI is a tool that will help you call, keep and fetch your useful scripts a git remote repository.
If you use a bunch of complex shell scripts and you are too lazy to manually add them to a file, send them to a server,
then fetch them from a server and maybe add them to an .*rc file when you need to run your scripts this tool is for you.

## Install

Before installing `ckp` you will need to have the golang package installed follow this [instructions](https://golang.org/dl/)

### Download

Download the lastest version [here](https://github.com/elhmn/ckp/releases)

Than copy the binary to your system binary `/bin` folder

### Setup git a repos

* Create an empty git repository to store your scripts and solutions we higly recommaned to keep this repository private

## Usage

#### How to `Init`-ialize `ckp`

This will create a `~/.ckp` folder, and clone the repository your scripts will be stored

```
ckp init [your-git-repo]
```

#### How to `Store` your scripts and solutions

The `store` command will record your script in ckp

```
ckp store -alias="sayHi" -comment="a script that says hi" 'echo sayHi'
```

#### How to `Push` your scripts to your remote solution repository

The `push` command will be commited and pushed to your remote repoitory

```
ckp push
```

#### How to `Pull` your scripts from your remote solution repository

The `pull` command will be commited and pushed to your remote repoitory

```
ckp push
```

## Commands

TODO...

## License

MIT.
