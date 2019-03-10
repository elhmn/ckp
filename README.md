# codekeeper : (ckp)
Simple CLI that calls, keeps and fetches your useful scripts from a git repository

## Overview [![GoDoc](https://godoc.org/github.com/elhmn/codekeeper?status.svg)](https://godoc.org/github.com/elhmn/ckp)

The codekeeper (ckp) CLI is a tool that will help you call, keep and fetch your useful scripts a git remote repository.
If you use a bunch of complex shell scripts and you are too lazy to manually add them to a file, send them to a server,
then fetch them from a server and maybe add them to an .*rc file when you need to run your scripts this tool is for you.

## Install

Before installing `ckp` you will need to have the golang package installed follow this [instructions](https://golang.org/dl/)

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
	{save,start,stop,sync,list}
        start		Clone your remote solution repoitory and init ckp
	stop		Remove ckp instance
	save		Save your scripts locally
	send		Send your local scripts to a remote server
	sync		Add your aliased scripts to your local .rc file
	list		List local scripts
```

### Start

Use `start` command to save your scrips locally

```
$ ckp start
usage: ckp start remote

Clone your scripts from the remote git repository and stores
it into $HOME/.ckp/repo folder.

Example : ckp start https://github/username/reponame

Positional arguments:
	[remote]	Is an alias to the remote server you want send your scripts
```


### Stop

Use `stop` command to remove your `.ckp` folder

```
$ ckp stop
usage: ckp stop

Simply remove $HOME/.ckp folder for the moment

```

### List

Use `list` to list your local scripts 

```
$ ckp list
usage: ckp list

list your local scripts

```

### Sync

Use `sync` to add your local scripts to your .zhshrc OR .bashrc OR .shrc 

```
$ ckp sync
usage: ckp sync

Sync will add all your aliases (scripts recorded with the -alias flag) script to your .rc files

Note : Don't forget to reload your .rc files use `source yourcfilepath`
```


### Send

Use `send` to commit and push your scripts to your remote git repository 

```
$ ckp send
usage: ckp send

Sync add, commit and push your local changes to the remote git repository
```

### Save

Use `save` command to save your scrips locally

```
$ ckp save
usage: ckp save [-flags] script

Saves your scripts locally.

Example : ckp save -alias=sayHi -comment="script that says hi" "echo Hi"

Flags:

	-file		(not yet implemented) Specify the file you want your scripts to be copied from
	-alias		Add an alias to your script
			The alias will be used by sync to alias your script in your .rc files 
	-comment	Add a comment to your script
```

## License

MIT.
