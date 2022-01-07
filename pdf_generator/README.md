# CV PDF Generator

This dart project will convert CV data into a good looking PDF.

## Development

```sh
dart run bin/pdf_generator.dart --dummy

# Linux
xdg-open example.pdf
# Macos
open example.pdf
# Windows (Powershell)
ii example.pdf
```

### Production

*The exe doesn't mean its only for windows it's just the extension of the file. I don't know why dart thought it would be a good idea to suffix all native binaries with .exe??*

```sh
dart compile exe bin/pdf_generator.dart

./bin/pdf_generator.exe --dummy
```
