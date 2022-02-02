#!/bin/bash
set -eo pipefail
die() {
    printf '%b\n' "$1" >&2
    exit "${2-1}"
}

assert() {
    result=$(./escape-any "$1")
    if [[ ! "$result" == "$2" ]]; then
        die "FAIL\nFORMAT:$1\nEXPECTED:$2\nRESULT:$result"
    fi
}

go build

printf "xxx\nbar" | assert "→, x→y,\n→|" "yyy|bar"
printf '\"qwe\"'  | assert '→, \"→"'     '"qwe"'
printf "foo\n"    | assert "→, \n→\\"    "foo\\"
printf "\nbar\m"  | assert '→, \n→\\n'  '\nbar\m'
printf "\nbar\n"  | assert '→, \n→\\n'  '\nbar\n'
printf '\\bar\\'  | assert '→, \\→n'    'nbarn'
