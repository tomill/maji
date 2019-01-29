# maji

maji is a simple command line tool to watch file changes using polling.

## Usage

```
maji [OPTIONS] [<dir>...] -- <command>

<dir>      Directories to watch. (default: .)
<command>  Command to run when an event fires.

Options:
      --dir=     Directories to watch. (default: .)
  -x, --exclude= Directory/File to ignore.
  -h, --help     Show this help message
```

Examples:

```bash
$ maji -- make test

$ maji src -x build -- make build

$ maji -- "curl localhost:8080 | jq ."
```

## Install

```
go get -u github.com/tomill/maji
```

## Description

This command was written as a go version of tokuhirom's perl5 [App-watcher](https://metacpan.org/pod/distribution/App-watcher/script/watcher). Unlike [gomon](https://github.com/c9s/gomon), this uses simple polling (by [radovskyb/watcher](https://github.com/radovskyb/watcher)) instead of inotify. That is it works with NFS file system.

## FAQ

What does "maji" mean?

"Maji" measn "seriously" in Japanese and "Maji maji" means "watching carefully 👀". So `maji` watches files pretty seriously.

## License

MIT
