#!/usr/bin/env bash

# Version of libbpf to fetch headers from
LIBBPF_VERSION=0.8.0

# The headers we want
prefix=libbpf-"$LIBBPF_VERSION"
headers=(
    "$prefix"/LICENSE.BSD-2-Clause
    "$prefix"/src/bpf_endian.h
    "$prefix"/src/bpf_helper_defs.h
    "$prefix"/src/bpf_helpers.h
    "$prefix"/src/bpf_tracing.h
)

# check if header files are already downloaded on local filesystem
# loop through headers array and check if each header file exists
for header in "${headers[@]}"; do
    # split header name into directory and file name
    dir=$(dirname "$header")
    file=$(basename "$header")
    if [ ! -f "$file" ] || [ ! -s "$file" ]; then
        echo "file $header not found. Fetching from github..."
        curl -sL "https://github.com/libbpf/libbpf/archive/refs/tags/v${LIBBPF_VERSION}.tar.gz" |
            tar -xz --xform='s#.*/##' "$header"
    fi
done
