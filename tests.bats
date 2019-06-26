@test "show help" {
	run ./pkgstats.sh -h
	echo "${lines[0]}" | grep -q 'usage:'
}

@test "show version" {
	run ./pkgstats.sh -v
	echo "${lines[0]}" | grep -q 'version'
}

@test "show information to be sent" {
	run ./pkgstats.sh -s
	echo "${lines[0]}" | grep -q 'Collecting data'
	echo "${output}" | grep -q 'packages='
	echo "${output}" | grep -q 'quiet=false'
}

@test "set quiet mode" {
	run ./pkgstats.sh -sq
	echo "${output}" | grep -q 'packages='
	echo "${output}" | grep -q 'quiet=true'
}
