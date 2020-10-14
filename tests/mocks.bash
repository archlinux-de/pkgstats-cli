#!/usr/bin/env bash

function pacman-conf() {
	echo https://my.mirror/extra/os/x86_64
}
export -f pacman-conf

function pacman() {
	echo mypackage
}
export -f pacman

#function curl() {
#	echo myresponse
#}
#export -f curl
