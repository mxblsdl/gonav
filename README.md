# GoNav

A simple CLI tool for navigating project folders

This project was originally build in python but I thought it would be fun to remake in Go.

The tool runs off of a config file that specifies folders to search. The main command is `go <arg>` which will search for partial matches on the `arg` then present a user with a menu if more than one match is found. The tool then opens up a folder window for the selected project, or opens the project in VS Code if the `-c` flag is supplied.

## TODOs

- make `go` recursive
- exclude venv folders by default
- speed tests?
