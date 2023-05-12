![donuts-are-good's followers](https://img.shields.io/github/followers/donuts-are-good?&color=555&style=for-the-badge&label=followers) ![donuts-are-good's stars](https://img.shields.io/github/stars/donuts-are-good?affiliations=OWNER%2CCOLLABORATOR&color=555&style=for-the-badge) ![donuts-are-good's visitors](https://komarev.com/ghpvc/?username=donuts-are-good&color=555555&style=for-the-badge&label=visitors)

# copy

copy is like `cp` but with a progress indicator and visual output

## usage

here's how to use `copy`:

```
copy -r /path/to/source /path/to/destination
```
`source` is the file or directory you want to copy, and `destination` is where you want to copy it to.

the `-r` flag means you want to copy recursively.

if the `source` is a directory and the `-r` flag is used, all its contents including subdirectories will be copied to the `destination` directory.

## examples

### copying a single file

to copy a single file, use the path to the file as the source and the path to the destination directory (or the full path including new file name) as the destination.

```
copy /home/user/documents/file.txt /home/user/desktop
```
this will copy `file.txt` from the `documents` directory to the `desktop` directory.

### copying a directory
to copy an entire directory, use the path to the directory as the `source` and the path to the parent of the destination directory as the `destination`, along with the `-r` flag.

```
copy -r /home/user/documents /home/user/desktop
```
this will copy the `documents` directory and all of its contents to the `desktop` directory.

## license

MIT License 2023 donuts-are-good, for more info see license.md
