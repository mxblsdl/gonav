#!/bin/bash

# Function to get next version number
get_next_version() {
    local version_type=${1}  # Default to patch if no argument

    if [ "$version_type" != "patch" ] && [ "$version_type" != "minor" ] && [ "$version_type" != "major" ]; then
        >&2 echo "Warning: Invalid version type '$version_type'. Defaulting to 'patch'." 
        >&2 echo "Use one of: 'patch', 'minor', 'major'" 
        version_type="patch"
    fi

    # Get latest version from GitHub releases
    local version=$(gh release list -L 1 --json tagName --jq '.[0].tagName' | sed 's/[^0-9.]//g')
    if [ -z "$version" ]; then
        version="0.0.0"  # Default if no releases exist
    fi
    # echo $version

    # Split the version
    local major=$(echo "$version" | cut -d. -f1)
    local minor=$(echo "$version" | cut -d. -f2)
    local patch=$(echo "$version" | cut -d. -f3)

    # Increment based on version type
    case "$version_type" in
        "patch")
            patch=$((patch + 1))
            ;;
        "minor")
            minor=$((minor + 1))
            patch=0
            ;;
        "major")
            major=$((major + 1))
            minor=0
            patch=0
            ;;
    esac

    echo "v$major.$minor.$patch"
}

# If script is being sourced, only define the function
# If script is being run directly, execute the function with arguments
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    get_next_version "$@"
fi