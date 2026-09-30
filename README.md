# GoNav

A simple CLI tool for navigating project folders

This project was originally [build in python](https://github.com/mxblsdl/pynav) but I thought it would be fun to remake in Go

## Installation

```bash
go install github.com/mxblsdl/gonav@latest
```

To build the binary from the current checkout:

```bash
go build -o nav .
```

To install the current checkout into your Go bin directory:

```bash
go install .
```

To build platform binaries and installers:

```bash
./sh/build.sh
```

The generated files are placed in `installer/binary/` and `installer/dist/`.

The tool runs off of a config file that specifies folders to search. Once installed run `nav` to pull up the help menu.

While there are probably other tools that perform this same function I really wanted to build something for myself to solve my exact need. I hate having to navigate through a folder system and remember exactly what a project is called.

## Future improvements

- exclude venv, node_modules, others by default