#!/bin/bash
# Windows resource generator script
# Usage: ./gen_winres.sh <version> [output_prefix]

set -e

VERSION="$1"
OUTPUT_PREFIX="${2:-openlist}"

if [ -z "$VERSION" ]; then
    echo "Error: Version parameter is required"
    echo "Usage: $0 <version> [output_prefix]"
    exit 1
fi

VERSION_CLEAN=$(echo "$VERSION" | sed 's/^v//')


VERSION_CSV=$(echo "$VERSION_CLEAN" | sed 's/\./, /g')

echo "Generating Windows resource file for version $VERSION_CLEAN"


sed -e "s/{{VERSION}}/$VERSION_CLEAN/g" \
    -e "s/{{VERSION_CSV}}/$VERSION_CSV/g" \
    openlist.rc > "${OUTPUT_PREFIX}.rc"

echo "Generated: ${OUTPUT_PREFIX}.rc"

if command -v x86_64-w64-mingw32-windres >/dev/null 2>&1; then
    x86_64-w64-mingw32-windres "${OUTPUT_PREFIX}.rc" -o "${OUTPUT_PREFIX}.syso"
    echo "Compiled: ${OUTPUT_PREFIX}.syso"
elif [ -n "$CC" ] && [[ "$CC" == *mingw* ]] && command -v "${CC%-gcc}-windres" >/dev/null 2>&1; then
    "${CC%-gcc}-windres" "${OUTPUT_PREFIX}.rc" -o "${OUTPUT_PREFIX}.syso"
    echo "Compiled: ${OUTPUT_PREFIX}.syso (using $CC toolchain)"
else
    echo "Note: Windows resource compiler not found, .syso file not generated"
    echo "The .rc file will be used during Go build process"
fi