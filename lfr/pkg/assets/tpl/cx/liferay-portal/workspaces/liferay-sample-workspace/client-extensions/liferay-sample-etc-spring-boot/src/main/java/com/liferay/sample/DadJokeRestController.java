/**
 * Copyright (c) 2000-present Liferay, Inc. All rights reserved.
 *
 * This library is free software; you can redistribute it and/or modify it under
 * the terms of the GNU Lesser General Public License as published by the Free
 * Software Foundation; either version 2.1 of the License, or (at your option)
 * any later version.
 *
 * This library is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
 * FOR A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more
 * details.
 */

package com.liferay.sample;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.reactive.function.client.WebClient;

/**
 * @author Raymond Augé
 * @author Gregory Amerson
 * @author Brian Wing Shun Chan
 */
@RequestMapping("/dad/joke")
@RestController
public class DadJokeRestController extends BaseRestController {

	@GetMapping
	public ResponseEntity<String> get(@AuthenticationPrincipal Jwt jwt) {
		log(jwt, _log);

		String dadJoke = WebClient.create(
		).get(
		).uri(
			"https://icanhazdadjoke.com"
		).accept(
			MediaType.TEXT_PLAIN
		).retrieve(
		).bodyToMono(
			String.class
		).block();

		if (_log.isInfoEnabled()) {
			_log.info("Dad joke: " + dadJoke);
		}

		return new ResponseEntity<>(dadJoke, HttpStatus.OK);
	}

	private static final Log _log = LogFactory.getLog(
		DadJokeRestController.class);

}