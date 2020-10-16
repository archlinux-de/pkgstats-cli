#!/usr/bin/env bats

function setup() {
	pushd $BATS_TEST_DIRNAME
	export PKGSTATS_URL=http://localhost:8888
	PKGSTATS="run ../pkgstats"
}

@test "show help" {
	$PKGSTATS -h
	echo "${lines[0]}" | grep -q 'usage:'
}

@test "show version" {
	$PKGSTATS -v
	echo "${lines[0]}" | grep -q 'version'
}

@test "show information to be sent" {
	$PKGSTATS -s
	echo "${lines[0]}" | grep -q 'Collecting data'
	echo "${output}" | grep -q 'packages='
	echo "${output}" | grep -q 'pacman'
	echo "${output}" | grep -q 'mirror=https://'
	echo "${output}" | grep -q 'quiet=false'
}

@test "set quiet mode" {
	$PKGSTATS -s -q
	echo "${output}" | grep -q 'packages='
	echo "${output}" | grep -q 'pacman'
	echo "${output}" | grep -q 'quiet=true'
}

@test "send informaition" {
	$PKGSTATS
	echo "${lines[0]}" | grep -q 'Collecting data'
	echo "${lines[1]}" | grep -q 'Submitting data'
	echo "${output}" | grep -q 'TEST OK'
}
