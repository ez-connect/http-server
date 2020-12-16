APP_NAME=http-static
BUILD_DIR=build
PLATFORMS=( windows linux darwin )
ARCHS=( 386 amd64 )

if [ -d $BUILD_DIR ]
then
    rm -rf $BUILD_DIR
fi

for p in "${PLATFORMS[@]}"
do
    for a in "${ARCHS[@]}"
    do
        echo $p-$a

        OUTPUT=$BUILD_DIR/$APP_NAME-$p-$a
        OUTPUT_ZIP=$OUTPUT.zip

        if [ "$p" = "windows" ]
        then
            OUTPUT=$OUTPUT.exe
        fi

        if [ "$p" != "darwin" ] || [ "$a" != "386" ]
        then
            GOOS=$p GOARCH=$a go build -ldflags="-s -w" -o $OUTPUT
            zip $OUTPUT_ZIP $OUTPUT
        fi
    done
done
