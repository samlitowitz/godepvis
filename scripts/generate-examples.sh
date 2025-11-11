#!/usr/bin/env bash

BASE_DIR="$(pwd)"
ASSETS=$BASE_DIR/assets
EXAMPLE_ASSETS=$ASSETS/examples
EXAMPLES_DIR=$BASE_DIR/examples
BUILD_FOR_MODULE_TEST_DIR=$BASE_DIR/internal/primitives/testdata/build-for-module

echo "Remove existing example outputs"
rm -rf $EXAMPLE_ASSETS/*

echo "Generate example outputs"

for d in $EXAMPLES_DIR/*/ ; do
    [ -L "${d%/}" ] && continue

    outputDir="$EXAMPLE_ASSETS/$(basename $d)"

    if [ ! -d "$outputDir" ]; then
      mkdir -p "$outputDir"
    fi

    palette=""
    if [ -f "${d}palette.yaml" ]; then
      palette="--palette ${d}palette.yaml"
    fi

    echo "Processing $d"

    echo "File Resolution"
    godepvis $palette --path $d --resolution file --dot $outputDir/file.dot
    dot -Tpng -o $outputDir/file.png $outputDir/file.dot

    echo "Package Resolution"
    godepvis $palette --path $d --resolution package --dot $outputDir/package.dot
    dot -Tpng -o $outputDir/package.png $outputDir/package.dot
done

for d in $BUILD_FOR_MODULE_TEST_DIR/*/ ; do
    [ -L "${d%/}" ] && continue

    outputDir="$EXAMPLE_ASSETS/$(basename $d)"

    if [ ! -d "$outputDir" ]; then
      mkdir -p "$outputDir"
    fi

    palette=""
    if [ -f "${d}palette.yaml" ]; then
      palette="--palette ${d}palette.yaml"
    fi

    echo "Processing $d"

    echo "File Resolution"
    godepvis $palette --path $d --resolution file --dot $outputDir/file.dot
    dot -Tpng -o $outputDir/file.png $outputDir/file.dot

    echo "Package Resolution"
    godepvis $palette --path $d --resolution package --dot $outputDir/package.dot
    dot -Tpng -o $outputDir/package.png $outputDir/package.dot
done
