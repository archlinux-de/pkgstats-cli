#!/usr/bin/env bats

function setup() {
	pushd $BATS_TEST_DIRNAME
	PKGSTATS="run ../../pkgstats --base-url http://localhost:8888"
}

@test "show help" {
	$PKGSTATS help
	echo "${output}" | grep -q 'Usage:'
}

@test "show version" {
	$PKGSTATS version
	echo "${lines[0]}" | grep -q 'version'
}

@test "show information to be sent" {
	$PKGSTATS submit --dump-json
	[ $(echo "${output}" | jq -r '.version') -eq 3 ]
	echo "${output}" | jq -r '.system.architecture' | grep -q '^x86_64'
	[ $(echo "${output}" | jq -r '.os.architecture') = 'x86_64' ]
	echo "${output}" | jq -r '.pacman.mirror' | grep -q '^https://'
	echo "${output}" | jq -r '.pacman.packages' | grep -q '"pacman-mirrorlist"'
}

@test "set quiet mode" {
	$PKGSTATS submit --quiet
	[ "$status" -eq 0 ]
	[ "$output" = "" ]
}

@test "send informaition" {
	$PKGSTATS submit
	echo "${lines[0]}" | grep -q 'Collecting data'
	echo "${lines[1]}" | grep -q 'Submitting data'
	echo "${lines[2]}" | grep -q 'Data were successfully sent'
}

@test "search packages" {
	$PKGSTATS search php
	echo "${lines[0]}" | grep -q 'php'
	echo "${lines[1]}" | grep -q 'php-fpm'
}

@test "show packages" {
	$PKGSTATS show php pacman
	echo "${lines[0]}" | grep -q 'php'
	echo "${lines[1]}" | grep -q 'pacman'
}
