#!/bin/sh

goPath="/usr/local/Cellar/go"
filename="goosee.go"

if [ -f ${goPath}/${filename} ]; then
    echo "file ${filename} found at at ${goPath}"
    exit 0
fi

echo "copying ${filename} to ${goPath}"
if [ -d ${goPath} ]
then
	cp ./goosee.go ${goPath}
	echo "file copied successfully"
else
    echo "go not found at ${goPath} "
fi
