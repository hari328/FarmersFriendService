#!/bin/sh

goPath="/usr/local/Cellar/go"
filename="goose"
isGlideInstalled="$(glide -v  >/dev/null 2>&1 || 'glide not installed')"


#check for invalid command line arguments
if [ $# -gt 1 ]; then
    echo "only one parameter accepted"
    echo "usage: setup/setup.sh <go path>"
    exit 1
fi

#check if go path is provided as commandline argument
if [ $# -eq 1 ]; then
    goPath=$1
    echo "go path(installed location) provided is: $1"
    echo "using $1 as go installed location"
fi

echo "using ${goPath} as the go installation folder"

#check if go in installed
if [ ! -d ${goPath} ]; then
    echo "go not found at ${goPath}"
    exit 1
fi

#add goose to go path
if [ -f ${goPath}/${filename} ]; then
    echo "${filename} file found at ${goPath}"
else
    echo "copying ${filename} to ${goPath}"
	cp ./setup/${filename} ${goPath}
	echo "file copied successfully"
fi

#install glide
if [ ! -z "$isGlideInstalled" ]
then
installStatus="$(curl https://glide.sh/get | sh > val)"
echo "$installStatus"
else
    echo "glide installed already"
fi

#remove previous dependencies for this service
if [ -d ./vendor ]; then
    rm -rf ./vendors
    echo "removed vendor folder."
fi

#intall dependencies
glide install

echo ""

numMigations="$(ls ./db/migrations | wc -l)"
echo "number of migrations: ${numMigations}"


goose -env=production -pgschema=farmerApp.db up
goose -env=production -pgschema=farmerApp.db status

