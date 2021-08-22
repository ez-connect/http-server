appName=http-server
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
		binaryName=$appName-$p-$a
		binaryFile=$buildDir/$appName-$p-$a

		echo $binaryFile

		if [ "$p" = "windows" ]
		then
			binaryFile=$binaryFile.exe
		fi

		if [ "$p" != "darwin" ] || [ "$a" != "386" ]
		then
			GOOS=$p GOARCH=$a go build -ldflags="-s -w" -o $binaryFile

			pushd $buildDir
			tar -zcvf $binaryName.tar.gz $binaryName
			popd
		fi
	done
done
