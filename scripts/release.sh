#!/bin/bash

set -e

DIR=$(cd "$(dirname "${0}")"/.. && pwd)
cd "${DIR}"
pwd

# cleanup
ARCHIVES="archives"
if [ -d ${ARCHIVES} ];then
    rm -rf "${ARCHIVES:?}"
fi

# make archives dir
DIST="dist"
mkdir -p "${ARCHIVES}/${DIST}"

# cross compile
XC_OS=${XC_OS:-linux windows}
XC_ARCH=${XC_ARCH:-386 amd64}
gox -output="${ARCHIVES}/{{.OS}}_{{.Arch}}/{{.Dir}}" \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    ./cmd/radigo
# for darwin
XC_OS=${XC_OS:-darwin}
XC_ARCH=${XC_ARCH:-amd64 arm64}
gox -output="${ARCHIVES}/{{.OS}}_{{.Arch}}/{{.Dir}}" \
    -os="${XC_OS}" \
    -arch="${XC_ARCH}" \
    ./cmd/radigo

# set the version number
VERSION=$(grep "const version " radigo.go | sed -E 's/.*"(.+)"$/\1/')
if [ "x${VERSION}" == "x" ]; then
     echo "missing version number" && exit 1
fi

# zip archives
for p in $(find "${ARCHIVES}" -mindepth 1 -maxdepth 1 -type d); do
    p_name=$(basename "${p}")
    if [ "${p_name}" == "${DIST}" ]; then
        continue
    fi

    pushd "${p}"
    archive_name=$(basename "${DIR}")_${VERSION}_${p_name}
    zip "${DIR}/${ARCHIVES}/${DIST}/${archive_name}.zip" ./*
    popd
done

# upload to GitHub release
ghr "${VERSION}" "${ARCHIVES}/${DIST}"
echo "success => https://github.com/yyoshiki41/radigo/releases"
