# `go-copy`

The tool for copying files which reports
progress and doesn't lie about it (too much)!

## Usage

```shell
go-copy --from source/path --to destination/path
```

For now

- the destination path has to include
    the filename and extension. It will never
    be inferred from the source path.
- the source path must be a file, not a directory

## Motivation

- I hate how `cp` doesn't report progress
- `rsync` has a lot of options to provide to get progress
- Due to some clever buffering, progress was never reported correctly on my Ubuntu laptop
- I don't think `cp`/`rsync` etc... give control over buffering/flushing to disk

## Goals

- Good enough copying performance that it's not annoying
- Reasonably accurate reporting of progress 
- Progress reflects what has actually been written to disk

## Non-Goals

- Be a feature complete replacement for `cp` or `rsync`
    - Keep this tool small, focussed and easy to maintain
- Highest possible copying performance
    - The focus is on realistically providing progress updates
      for transfers that will definitely take a long time
