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
						.filters(f -> f.stripPrefix(1))
						.uri("lb://AUTH-SERVICE"))
				.route("user-service", r -> r.path("/user/**")
						.filters(f -> f
								.stripPrefix(1))
						.uri("lb://USER-SERVICE"))
				.route("flora-upstream-service", r -> r.path("/flora-upstream/**")
						.filters(f -> f
								.stripPrefix(1))
						.uri("lb://FLORA-UPSTREAM-SERVICE"))
				.build();
	}

}
