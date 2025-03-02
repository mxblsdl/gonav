source ./sh/build.sh
source ./sh/version_number.sh

version=$(get_next_version "$1")

gh release create "$version" \
 --notes-file release_notes.md \
 installer/dist/nav-linux-amd64 installer/dist/nav-darwin-amd64 installer/dist/nav-windows-amd64.exe