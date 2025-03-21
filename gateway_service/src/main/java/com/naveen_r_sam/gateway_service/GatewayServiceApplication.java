//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.

package com.naveen_r_sam.gateway_service;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.gateway.route.RouteLocator;
import org.springframework.cloud.gateway.route.builder.RouteLocatorBuilder;
import org.springframework.context.annotation.Bean;

@SpringBootApplication
public class GatewayServiceApplication {

	public static void main(String[] args)
	{
		SpringApplication.run(GatewayServiceApplication.class, args);
	}
	@Bean
	public RouteLocator customRouteLocator(RouteLocatorBuilder builder) {
		return builder.routes()
				.route("auth-service", r -> r.path("/auth/**")
						.filters(f -> f
								.stripPrefix(1))
						.uri("lb://AUTH-SERVICE"))
				.route("user-service", r -> r.path("/user/**")
						.filters(f -> f
								.stripPrefix(1))
						.uri("lb://USER-SERVICE"))
				.route("gene-bank-service", r -> r.path("/gene-bank/**")
						.filters(f -> f
								.stripPrefix(1))
						.uri("lb://GENE-BANK-SERVICE"))
				.route("flora-upstream-service", r -> r.path("/flora-upstream/**")
						.filters(f -> f
								.stripPrefix(1))
						.uri("lb://FLORA-UPSTREAM-SERVICE"))
				.build();
	}

}
