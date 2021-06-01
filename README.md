# codekeeper : (ckp)
CLI that helps you store and reuse your common scripts and solutions from anywhere

## Overview

If you ever found yourself using a bunch of complex scripts or useful bash oneliner and you find it hard to manually add them to a file, send them to a server and then fetch them to this new machine you have recently acquired or ssh-ed into, this tool is for you.

![ckp_demo](https://user-images.githubusercontent.com/5704817/120272338-39377a80-c2ad-11eb-9058-a16f98745bb1.gif)

## Prerequisite
`ckp` uses several dependencies such as:
1. `git` version >= 2.24.3 you can follow this [steps](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) to install git
2. only `bash` compatible commands can be run using `ckp`, you can use the tool without `bash` but you won't be able to run your commands using the CLI

## Install

#### Using the install script:

Run
```sh
$> curl https://raw.githubusercontent.com/elhmn/ckp/master/install.sh | bash
```
It will create a `./bin/ckp` binary on your machine
In order to run the command add it to your `/usr/local/bin`
```sh
$> cp ./bin/ckp /usr/local/bin
```

#### Using homebrew:

Run
```sh
 $> brew tap elhmn/ckp https://github.com/elhmn/ckp
 $> brew install ckp
```

#### Download

Download the lastest version [here](https://github.com/elhmn/ckp/releases)
Then copy the binary to your system binary `/usr/local/bin` folder


## Usage

#### How to `Init`-ialize `ckp`

1. You first need to create an empty git repository that `ckp` will use as a storage. we higly recommend to keep this repository private

2. Once the repository is created you can initialise `ckp` using the init command.
    Copy the ssh or https url and pass it as an argument to the `ckp init` command

```sh
$> ckp init git@github.com:elhmn/store.git
```

This will create a `~/.ckp` folder, and clone the storage repository

#### How to `Add` your scripts and solutions

The `add code` command will store your script as a code entry in ckp

```
ckp add code 'echo say hi!' --alias="sayHi" --comment="a script that says hi"
```

The `add solution` command will store your script as a solution entry in ckp

```
ckp add solution 'https://career-ladders.dev/engineering/' --comment="carreer ladders"
```

#### How to `Push` your scripts to your remote storage repository

The `push` command will be commited and pushed to your remote repoitory

```
ckp push
```

#### How to `Pull` your scripts from your remote storage repository

The `pull` command will pull changes from your remote storage repository

```
ckp pull
```

#### How to `Find` a script or solution

The `find` command will prompt a search and selection UI, that can be used to find

```
ckp find
```

#### How to `Run` a script or solution

The `run` command will prompt a search and selection UI, that can be used to find and run a specific script

```
ckp run
```


## License

MIT.
