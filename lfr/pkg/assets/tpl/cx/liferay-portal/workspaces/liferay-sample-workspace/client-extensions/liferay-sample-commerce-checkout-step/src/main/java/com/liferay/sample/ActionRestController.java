/**
 * SPDX-FileCopyrightText: (c) 2024 Liferay, Inc. https://liferay.com
 * SPDX-License-Identifier: LGPL-2.1-or-later OR LicenseRef-Liferay-DXP-EULA-2.0.0-2023-06
 */

package com.liferay.sample;

import com.liferay.petra.string.StringBundler;

import org.json.JSONObject;

import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.reactive.function.client.WebClient;

/**
 * @author Andrea Sbarra
 */
@RequestMapping("/action")
@RestController
public class ActionRestController extends BaseRestController {

	@PostMapping
	public ResponseEntity<String> post(
		@AuthenticationPrincipal Jwt jwt, @RequestBody String json) {

		JSONObject jsonObject = new JSONObject(json);

		return new ResponseEntity<>(
			WebClient.create(
				StringBundler.concat(
					lxcDXPServerProtocol, "://", lxcDXPMainDomain,
					"/o/headless-commerce-delivery-cart/v1.0/carts/",
					jsonObject.getLong("commerceOrderId"))
			).patch(
			).accept(
				MediaType.APPLICATION_JSON
			).contentType(
				MediaType.APPLICATION_JSON
			).bodyValue(
				new JSONObject(
				).put(
					"purchaseOrderNumber", jsonObject.getString("pon")
				).toString()
			).header(
				HttpHeaders.AUTHORIZATION, "Bearer " + jwt.getTokenValue()
			).retrieve(
			).bodyToMono(
				String.class
			).block(),
			HttpStatus.OK);
	}

}