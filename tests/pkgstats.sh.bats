#!/usr/bin/env bats

function setup() {
	pushd $BATS_TEST_DIRNAME
	load mocks
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
#	echo "${output}" | grep -q 'mypackage'
#	echo "${output}" | grep -q 'mirror=https://my.mirror/'
#	echo "${output}" | grep -q 'quiet=false'
}

@test "set quiet mode" {
#	$PKGSTATS -sq
	$PKGSTATS -s -q
	echo "${output}" | grep -q 'packages='
#	echo "${output}" | grep -q 'mypackage'
#	echo "${output}" | grep -q 'quiet=true'
}

@test "send informaition" {
	$PKGSTATS
	echo "${lines[0]}" | grep -q 'Collecting data'
	echo "${lines[1]}" | grep -q 'Submitting data'
#	echo "${output}" | grep -q 'myresponse'
}
