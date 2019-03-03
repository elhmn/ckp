# codekeeper : (ckp)
Simple CLI that calls, keeps and fetches your useful scripts from any remote

## Overview [![GoDoc](https://godoc.org/github.com/elhmn/codekeeper?status.svg)](https://godoc.org/github.com/elhmn/codekeeper)

The codekeeper (ckp) CLI is a tool that will help you call, keep and fetch your useful scripts on any remote http server.
If you use a bunch of complex shell scripts and you are too lazy to manually add them to a file, send them to a server,
then fetch them from a server and maybe add them to an .*rc  file when you need to make call to your scripts this tool is for you.

## Install

Before installing you will need to have the golang package installed follow this [instructions](https://golang.org/dl/)

```
go get github.com/elhmn/codekeeper/ckp
```

## Commands

`ckp` provides several commands to manage your scripts.

```
$ ckp help
usage: ckp help commands

A tool to manage your scripts.

positional arguments:
	{save,fetch,run,start,setup,open,debug}
                        Commands
	save		Save your scripts locally
	send		Send your local scripts to a remote server
	fetch           Fetch your scripts from a remote server
	sync		Add your aliased scripts to your local .rc file
	run			Run your local scripts
	edit		Edit your local scripts
	list		List local scripts.
	remote		Display a list of recorded remotes.

```

### Save

Use `save` command to save your scrips locally

```
$ ckp save
usage: ckp save [script] [remote | url]

Saves your scripts locally.

Positional arguments:
	{scripts, remote, url}
	[scripts]	Is the string of the script you want to save locally
	[remote]	Is an alias to the remote server you want send your scripts
	[url]		Is the http url of the remote server you want your file to be stored

	Note: it only works with git or bitbucket bare repositories remotes and with shell scripts only at the moment

Flags:

	-file		Specify the file you want your scripts to be copied from
```

## License

MIT.
