#!/bin/bash

# Build script for terraform-provider-uptimekuma

set -e

# Default values
VERSION=${VERSION:-"1.0.0"}
PLATFORMS=${PLATFORMS:-"linux/amd64 darwin/amd64 windows/amd64"}

echo "Building terraform-provider-uptimekuma v${VERSION}"

# Clean previous builds
rm -rf dist/
mkdir -p dist/

# Build for each platform
for platform in $PLATFORMS; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    
    output_name="terraform-provider-uptimekuma_v${VERSION}"
    if [ $GOOS = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    
    env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build \
        -ldflags="-w -s -X main.version=${VERSION}" \
        -o dist/${output_name} \
        .
    
    # Create zip archive
    cd dist/
    if [ $GOOS = "windows" ]; then
        zip terraform-provider-uptimekuma_${VERSION}_${GOOS}_${GOARCH}.zip ${output_name}
    else
        tar -czf terraform-provider-uptimekuma_${VERSION}_${GOOS}_${GOARCH}.tar.gz ${output_name}
    fi
    rm ${output_name}
    cd ..
done

echo "Build complete! Artifacts are in dist/"
ls -la dist/
