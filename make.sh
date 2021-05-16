appName=http-static
buildDir=build
platforms=( windows linux darwin )
archs=( amd64 )

if [ -d $buildDir ]
then
	rm -rf $buildDir
fi

for p in "${platforms[@]}"
do
	for a in "${archs[@]}"
	do
		echo $p-$a

		binaryFile=$buildDir/$appName-$p-$a
		zipFile=$binaryFile.zip

		if [ "$p" = "windows" ]
		then
			binaryFile=$binaryFile.exe
		fi

		if [ "$p" != "darwin" ] || [ "$a" != "386" ]
		then
			GOOS=$p GOARCH=$a go build -ldflags="-s -w" -o $binaryFile
			# zip $zipFile $binaryFile
		fi
	done
done
