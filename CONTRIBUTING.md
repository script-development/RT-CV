# Contributing to RT-CV

## Requirements

If you are using VSCode you might need to install the linter used in this project:

```
go install github.com/mgechev/revive@latest
```

## Code guidelines

- Try to code as if you leave tomorrow and some else picks up the code
- Avoid using interface with methods as they make it harder to click tough the code, if you need an interface with methods for some reason make sure to document it well see the `db` package for an example
- Add tests for your code to make sure the code is correct and to add essentially extra documentation
- Avoid adding packages of minimal extra value, every extra package adds external documentation that a maintainer has to look at to understand. When things are only in this repo you don't have to leave your code editor to understand new things.

## New to MongoDB / NoSQL?

This video explains well what a NoSQL database is (MongoDB is a NoSQL database): [https://youtu.be/v_hR4K4auoQ](https://youtu.be/v_hR4K4auoQ) _(What is a NoSQL Database? How is Cloud Firestore structured? | Get to know Cloud Firestore #1)_
_Note that some information is cloud firestore specific in that video but it should give you an overall idea on what NoSQL is_
