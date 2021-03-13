<?php

declare(strict_types=1);
error_reporting(E_ALL);

$userAgent = $_SERVER['HTTP_USER_AGENT'];
$requestUri = $_SERVER['REQUEST_URI'];
$requestPath = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);

error_log('Got request from ' . $userAgent . ' on ' . $requestUri);
if (!preg_match('#^pkgstats/[\w.-]+$#', $userAgent)) {
	header('HTTP/1.0 400');
	error_log('Invalid user agent');
	exit(1);
}

switch ($requestPath) {
	case '/api/submit':
		if (!isset($_GET['redirect'])) {
			error_log('Testing redirect');
			header('HTTP/1.0 308');
			header('Location: /api/submit?redirect=1');
			exit();
		}

		$request = json_decode(file_get_contents('php://input'), true);

		if (
			$request['version'] === '3'
			&& $request['os']['architecture'] === php_uname('m')
			&& in_array($request['system']['architecture'], ['x86_64', 'x86_64_v2', 'x86_64_v3', 'x86_64_v4'])
			&& preg_match('#^https?://.+$#', $request['pacman']['mirror'])
			&& count($request['pacman']['packages']) > 1
			&& in_array('pacman-mirrorlist', $request['pacman']['packages'])
		) {
			error_log('Request was vaild');
			header('HTTP/1.0 204');
		} else {
			error_log('Request was invaild');
			error_log(print_r($request, true));
			header('HTTP/1.0 400');
			echo 'TEST FAILED';
		}
		break;

	case '/api/packages':
		$response = json_encode([
			'total' => 42,
			'count' => 2,
			'packagePopularities' => [
				[
					'name' => 'php',
					'popularity' => 56.78
				],
				[
					'name' => 'php-fpm',
					'popularity' => 12.34
				]
			]
		]);
		echo $response;
		break;

	case '/api/packages/pacman':
		$response = json_encode([
			'name' => 'pacman',
			'popularity' => 12.34
		]);
		echo $response;
		break;

	case '/api/packages/php':
		$response = json_encode([
			'name' => 'php',
			'popularity' => 56.78
		]);
		echo $response;
		break;

	default:
		error_log('Unknown request');
		header('HTTP/1.0 400');
}
