# codekeeper : (ckp)
CLI that helps you store and reuse your history and one liner scripts from anywhere, better than gists.

## Overview

If you ever found yourself using a bunch of complex scripts or useful bash oneliners and you find it hard to manually add them to a file, send them to a server and then fetch this scripts to that new machine you have recently acquired or ssh-ed into, this tool is for you. Store and fetch your scripts, your terminal history and your notes from anywhere.

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

#### How to set your text editor

Vim is the default text editor to use a different code editor you might need to create a `~/.ckp/config.yaml` file,
then open the file and set the `editor` field as follows.

```yaml
editor: nano
```


#### How to `Add` your scripts and solutions

The `add code` command will store your script as a code entry in ckp.

```sh
$> ckp add code 'echo say hi!' --alias="sayHi" --comment="a script that says hi"
```

The `add solution` command will store your script as a solution entry in ckp.

```sh
$> ckp add solution 'https://career-ladders.dev/engineering/' --comment="carreer ladders"
```

#### How to add scripts from my `bash_history` or `zh_history`

The `add history` command will read scripts from your history files and store them in ckp.
the `--skip-secrets` flag will force ckp to skip scripts that potentially contains secrets.

```sh
$> ckp add history --skip-secrets
```

#### How to `Push` your scripts to your remote storage repository

The `push` command will be commited and pushed to your remote repoitory.

```sh
$> ckp push
```

#### How to `Pull` your scripts from your remote storage repository

The `pull` command will pull changes from your remote storage repository.

```sh
$> ckp pull
```

#### How to `Find` a script or solution

The `find` command will prompt a search and selection UI, that can be used to find.

```sh
$> ckp find
```

To find a script in your history.

```sh
$> ckp find --from-history
```


#### How to `Run` a script or solution

The `run` command will prompt a search and selection UI, that can be used to find and run a specific script.

```sh
$> ckp run
```

To run a script from your history.

```sh
$> ckp run --from-history
```


#### How to `Remove` a script or solution

The `rm` command will prompt a search and selection UI, that can be used to find and run a specific script.

```sh
$> ckp rm
```

To remove a script from your history.

```sh
$> ckp rm --from-history
```


## License

MIT.
