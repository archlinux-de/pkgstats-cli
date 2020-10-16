<?php

$userAgent = $_SERVER['HTTP_USER_AGENT'];
$arch = $_POST['arch'];
$cpuarch = $_POST['cpuarch'];
$mirror = $_POST['mirror'];
$packages = explode("\n", $_POST['packages']);
$quiet = $_POST['quiet'];

error_log('Got request from ' . $userAgent);

if (
	strpos($userAgent, 'pkgstats/' . exec('git describe --tags')) === 0
	&& $arch === 'x86_64'
	&& $cpuarch === 'x86_64'
	&& strpos($mirror, 'http') === 0
	&& count($packages) > 100
	&& $quiet === 'false'
) {
	echo 'TEST OK';
	error_log('Request was vaild');
} else {
	echo 'TEST FAILED';
	error_log('Request was invaild');
	error_log(print_r($_POST, true));
}
