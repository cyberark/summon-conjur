#!/bin/bash

##############
# Finds installs and runs go mod tidy for the matching go.mod version
# and utilizes veondoring to create a vendor directory for future
# dependency reference without the need to authenticate and download
# dependencies repeatedly.
#
### LOCAL SYSTEM ALTERATIONS ###
#
# This script will utilize golang's built in system to install multiple
# go versions. The installed version will match the go.mod version (this
# won't alter the system default golang version).
#
# It will alter the GOPRIVATE go env to allow authentication
# to github.cyberng.com for private dependencies.
#
# This will modify git config --global to allow authentication to
# github.cyberng.com for private dependencies via ssh.
#
### LOCAL SYSTEM ALTERATIONS ###
#
# -compat is assumed to match the go.mod version otherwise
# if the Go project requires backwards compatibility
# the -compat version must be explicitly set.
#
# This is to prevent go mod tidy default params
# By default, tidy acts as if the -compat flag were set to the
# version prior to the one indicated by the 'go' directive in the go.mod
# file.
#
##############

function main {

    while [[ $# -gt 0 ]]; do
        case "$1" in
        -g | --gomod)
            local go_mod="$2"
            shift
            ;;
        -c | --compat)
            local compat="$2"
            shift
            ;;
        -h | --help)
            print_help
            exit 1
            ;;
        *)
            echo "Unknown option: $1"
            print_help
            exit 1
            ;;
        esac
        shift
    done

    echo -n "checking for go.mod file..."
    if [[ ! -f "${go_mod}" ]]; then
        echo "go.mod file doesn't exist."
        print_help
        exit 1
    fi
    echo "found"

    echo -n "checking if file is a go.mod..."
    local go_mod_match
    go_mod_match='go\.mod$'
    if ! [[ "${go_mod}" =~ ${go_mod_match} ]]; then
        echo "file must be a go.mod."
        print_help
        exit 1
    fi
    echo "true"

    echo -n "setting script GOPATH..."
    # Check if GOPATH is set, otherwise set it to $HOME/go/bin
    if [ -n "${GOPATH}" ]; then
        GO_PATH=$GOPATH
    else
        GO_PATH="$HOME/go/bin/"
    fi
    echo "set"

    # Extract golang version from go.mod
    echo -n "extracting go.mod version..."
    GO_MOD_VERSION=$(grep -E "^go\ [0-9]+\.[0-9]+" "${go_mod}" | grep -oE "[0-9]+\.[0-9]+")
    echo "done"

    echo -n "adding private go.mod replace statements..."
    # Add the private replace repo statements to go.mod
    go_mod_replace
    echo "done"

    echo -n "setting up git for private dependencies..."
    # Run on Jenkins agent only
    if [ -n "${BRANCH_NAME}" ]; then
        if ! jenkins_git_setup; then
            echo "Error: Failed to setup git on Jenkins for private dependencies"
            exit 1
        fi
    else
        if ! local_git_setup; then
            echo "Error: Failed to setup git for private dependencies"
            exit 1
        fi
    fi
    echo "done"

    echo -n "comparing go.mod version to Go 1.16..."
    local compare_go_mod_version
    compare_go_mod_version=$(compare_versions)
    echo "done"

    echo -n "constructing go mod tidy command based on version..."
    # Construct the go mod tidy command
    declare -a go_mod_tidy
    read -ra go_mod_tidy < <(go_mod_command "${compat}" "${compare_go_mod_version}")
    echo "done"

    echo "installing go version..."
    # Install the go version specified in go.mod
    if ! go get golang.org/dl/go"${GO_MOD_VERSION}"; then
        echo "Error: Failed to get go version ${GO_MOD_VERSION}"
        exit 1
    fi

    if ! go install golang.org/dl/go"${GO_MOD_VERSION}@latest"; then
        echo "Error: Failed to install go version ${GO_MOD_VERSION}"
        exit 1
    fi 

    if ! "${GO_PATH}go${GO_MOD_VERSION}" download; then
        echo "Error: Failed to download go version ${GO_MOD_VERSION}"
        exit 1
    fi

    echo "go version ${GO_MOD_VERSION} installed"

    echo "running go mod tidy..."
    # Run go mod tidy
    if ! "${go_mod_tidy[@]}"; then
        echo "Error: Failed to run go mod tidy"
        exit 1
    fi
    echo "go mod tidy finished"

    echo -n "creating vendor directory..."
    # Create a vendor directory for usage locally and/or on infrapool
    # This prevents the need to go mod download dependencies for future
    # go build or go run calls
    if ! go mod vendor; then
        echo "Error: Failed to create vendor directory"
        exit 1
    fi
    echo "vendor created successfully"

}

function print_help {
    echo "Usage: $0 [-g|--gomod] [-c|--compat]"
    echo "  -g|--gomod: Path to go.mod file"
    echo "  -c|--compat: Optional version to use for backwards compatibility"
    echo "  -h|--help: Print this help message"

}

function compare_versions() {
    # Compare go.mod version against Go 1.16
    #
    # Prior to Go 1.17 `go mod tidy` doesn't
    # support the -compat argument.
    #
    # Return an integer depending on if the go.mod version is
    # < > or = to 1.16
    local unsupported_go_version="1.16"

    if [[ "$GO_MOD_VERSION" == "$unsupported_go_version" ]]; then
        echo 0 && return 0
    fi

    # Sort the go.mod version and unsupported_go_version
    # to determine which is greater
    if [[ $(printf "%s\n%s" "$GO_MOD_VERSION" "$unsupported_go_version" | sort -V | head -1) == "$GO_MOD_VERSION" ]]; then
        echo 1 && return 1
    else
        echo 2 && return 2
    fi

}

function go_mod_command() {
    # Constructs the go mod tidy command

    local compat="$1"
    local compare_go_mod_version="$2"
    local go_mod_tidy_command

    # If the go.mod version is less than or equal to 1.16
    # `go mod tidy` is called
    if [[ ${compare_go_mod_version} -eq 0 || ${compare_go_mod_version} -eq 1 ]]; then
        go_mod_tidy_command=("${GO_PATH}go${GO_MOD_VERSION}" "mod" "tidy")
        echo "${go_mod_tidy_command[@]}"
        return 0
    elif [[ "${compat}" != null ]]; then
        go_mod_tidy_command=("${GO_PATH}go${GO_MOD_VERSION}" "mod" "tidy" "--compat=${compat}")
        echo "${go_mod_tidy_command[@]}"
        return 0
    else
        go_mod_tidy_command=("${GO_PATH}go${GO_MOD_VERSION}" "mod" "tidy" "--compat=${GO_MOD_VERSION}")
        echo "${go_mod_tidy_command[@]}"
        return 0
    fi

}

function jenkins_git_setup() {
    # Setup github.cyberng.com for a successful go mod tidy/vendor invocation on the Jenkins agent

    echo https://${GITHUB_USER}:${TOKEN}@github.cyberng.com > "${HOME}/.git-credentials"
    git config --global credential.helper "store --file ${HOME}/.git-credentials" || return 1
    go env -w GOPRIVATE=github.cyberng.com/conjur-enterprise/* || return 1
    return 0

}

function local_git_setup() {
    # Setup github.cyberng.com for a successful go mod tidy/vendor invocation on local machine
    git config --global url."git@github.cyberng.com:".insteadOf "https://github.cyberng.com/" || return 1
    go env -w GOPRIVATE=github.cyberng.com/conjur-enterprise/* || return 1
    return 0

}

function go_mod_replace() {
    # Append replace github.com/cyberark dependencies with github.cyberng.com inside go.mod

    local github_cyberark_regex
    github_cyberark_regex="github\.com\/cyberark\/.*"

    local go_mod_replace_statements
    go_mod_replace_statements=$(grep -E "^replace ${github_cyberark_regex}$" go.mod)
    if [ -n "$go_mod_replace_statements" ]; then
        echo -n "Existing replace statements found, overwriting with private repo replace statements..." >&2
        # Remove existing replace statements
        sed -i -E "/^replace ${github_cyberark_regex}$/d" go.mod
    else
        echo -n "No replace statements found..." >&2
        return 0
    fi

    # Get all go.mod require statements that match github.com/cyberark and sed remove special characters and require
    local go_mod_require_statements
    go_mod_require_statements=$(grep -E "^\W+${github_cyberark_regex}$|require ${github_cyberark_regex}$" go.mod | sed -E 's/^\s+|^\t+|require //')

    # Remove comments from statements
    go_mod_require_statements=$(echo "${go_mod_require_statements}" | sed -E 's/\/\/.*$//')

    # For each require statement found, replace github.com/cyberark with github.cyberng.com/conjur-enterprise
    while IFS= read -r line; do
        github_org='github.com/cyberark'
        # Remove github.com/cyberark from line
        repo_ver=${line/$github_org/}
        # Add a new line to go.mod
        echo "" >> go.mod

        if [ -n "$go_mod_replace_statements" ]; then
            version=' v[0-9]*'
            # Remove version from repo_ver
            repo=${repo_ver/$version/}
            # Remove version from line
            line=${line/$version/}
            echo "replace ${line} => github.cyberng.com/conjur-enterprise${repo} latest" >> go.mod
        else
            echo "replace ${line} => github.cyberng.com/conjur-enterprise${repo_ver}" >> go.mod
        fi
    done < <(printf '%s\n' "$go_mod_require_statements")

}

main "$@"
