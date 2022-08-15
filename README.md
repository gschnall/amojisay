# amojisay

(˵╹◡╹)━☆ ･ﾟﾟ Terminal cli for one line ascii emojis

## Getting Started

Place the executable in your `/usr/local/bin` directory

### Linux

```shell
wget -O amojisay.tar.gz https://github.com/gschnall/amojisay/releases/download/v0.2.0/amojisay_0.2.0_Linux_x86_64.tar.gz && tar xvf amojisay.tar.gz && sudo mv amojisay /usr/local/bin
```

### Mac & Windows

1. Download the zip file for your sysytem
   https://github.com/gschnall/amojisay/releases  
   for macOs use `Darwin_x86_64.tar.gz`
2. Unzip the contents
3. move the `amojisay` executable to your `/usr/local/bin` directory
4. once amojisay has been started, use the --help arg to bring up a help menu

## Features

- Generate one line ascii emojis on the command line
- Over 250 to choose from

## Args

- `-a` | select ascii emoji
- `-l` | list all emojis
- `-p` | prepend text
- `-s` | use string substition

## Examples

```sh
> amojisay -a orly "you don't say..."
```

```
(눈_눈) you don't say...
```

---

```sh
> amojisay -p -a run 'Run away!'
```

```
Run away! (╯°□°)╯
```

---

```sh
> amojisay -s '%{tada} %{sparkle} Great Job! %{star} %{sparkle2}'
```

```
⊂(o‿o)つ *✧ ･ﾟ Great Job! ★  ･ﾟ✧ *:･ﾟ`
```

---

```sh
> amojisay -a hello 'Greetings human' | cowsay
```

```
________________________
< (ʘ‿ʘ)╯ Greetings human >
 ------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
```

## Issues

Find a bug or want a new feature? Feel free to create an issue

## Contributions

Create a new branch and submit a Pull request with screenshots and a description of changes.

## Licensing

MIT - see LICENSE
