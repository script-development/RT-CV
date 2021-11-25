# Contributing to RT-CV

## Requirements

If you are using VSCode you might need to install the linter used in this project:

```
go install github.com/mgechev/revive@latest
```

## Guidelines

- Try to code as if you leave tomorrow and some else picks up the code
- Avoid using interface with methods as they make it harder to click tough the code, if you need an interface with methods for some reason make sure to document it well see the `db` package for an example
- Add tests for your code to make sure the code is correct and for basically free documentation
- Avoid adding packages of minimal extra value, every extra package adds external documentation that a maintainer has to look at to understand. When things are only in this repo you don't have to leave your code editor to understand new things.
- Try to create meaning full commits
  - The commit message should represent the change
  - A commit message can have multiple lines, please use those to write down extra information and toughs about the change so other can later look back at your change and know why you did something
  - Create multiple commits if your changes are for 2 or more things
    - `git add --patch` _(or `-p`)_ is a great way to select only specific changes for a commit

## New to MongoDB / NoSQL?

This video explains well what a NoSQL database is (MongoDB is a NoSQL database): [youtu.be/v_hR4K4auoQ](https://youtu.be/v_hR4K4auoQ) _(What is a NoSQL Database? How is Cloud Firestore structured? | Get to know Cloud Firestore #1)_
_Note that some information is cloud firestore specific in that video but it should give you an overall idea on what NoSQL is_

## Debug performance issues

A profile of the program can be created by starting the program using

```sh
go run . -profile
```

After that you can inspect the profile in the browser using

```
go tool pprof -http localhost:3333 cpu.profile
```
