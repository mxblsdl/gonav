# GoNav

A simple CLI tool for navigating project folders

This project was originally [build in python](https://github.com/mxblsdl/pynav) but I thought it would be fun to remake in Go.

The tool runs off of a config file that specifies folders to search. The main command is `go <arg>` which will search for partial matches on the `arg` then present a user with a menu if more than one match is found. The tool then opens up a folder window for the selected project, or opens the project in VS Code if the `-c` flag is supplied.

While there are probably other tools that perform this same function I really wanted to build something for myself to solve my exact need. I hate having to navigate through a folder system and remember exactly what a project is called.

## Future improvements

- make `go` recursive with a max depth parameter specified
- exclude venv, node_modules, others by default
- make binary executable somehow
  - tarball?
