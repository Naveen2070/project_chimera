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
package com.naveen_r_sam.gateway_service.config;

import jakarta.annotation.PostConstruct;
import org.springdoc.core.properties.SwaggerUiConfigProperties;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.cloud.client.ServiceInstance;
import org.springframework.cloud.client.discovery.event.HeartbeatEvent;
import org.springframework.context.ApplicationListener;
import org.springframework.context.annotation.Configuration;

import java.util.*;
import java.util.concurrent.locks.ReentrantLock;

@Configuration
public class SwaggerUiConfig implements ApplicationListener<HeartbeatEvent> {

    private final DiscoveryClient discoveryClient;
    private final SwaggerUiConfigProperties swaggerUiConfigProperties;
    private static final String GATEWAY_URL = "http://localhost:8080"; // Use the gateway
    private final ReentrantLock lock = new ReentrantLock();
    private Set<String> lastKnownServices = new HashSet<>();

    public SwaggerUiConfig(DiscoveryClient discoveryClient,
                           SwaggerUiConfigProperties swaggerUiConfigProperties) {
        this.discoveryClient = discoveryClient;
        this.swaggerUiConfigProperties = swaggerUiConfigProperties;
    }

    @PostConstruct
    public void configureSwaggerUi() {
        updateSwaggerUrls(true);
    }

    @Override
    public void onApplicationEvent(HeartbeatEvent event) {
        updateSwaggerUrls(false);
    }

    private void updateSwaggerUrls(boolean forceUpdate) {
        lock.lock();
        try {
            List<String> services = discoveryClient.getServices();
            Set<String> currentServices = new HashSet<>(services);

            // If services haven't changed, skip the update to avoid unnecessary processing
            if (!forceUpdate && currentServices.equals(lastKnownServices)) {
                return;
            }

            Set<SwaggerUiConfigProperties.SwaggerUrl> urls = new LinkedHashSet<>();

            for (String serviceName : services) {
                if (serviceName.equalsIgnoreCase("consul")|| serviceName.equalsIgnoreCase("gateway-service")) {
                    continue;
                }

                List<ServiceInstance> instances = discoveryClient.getInstances(serviceName);
                if (instances.isEmpty()) {
                    continue;
                }

                String cleanedServiceName = removeServiceSuffix(serviceName);
                String apiDocsUrl = resolveApiDocsUrl(cleanedServiceName);

                SwaggerUiConfigProperties.SwaggerUrl swaggerUrl = new SwaggerUiConfigProperties.SwaggerUrl();
                swaggerUrl.setName(serviceName + " (via Gateway)");
                swaggerUrl.setUrl(apiDocsUrl);

                urls.add(swaggerUrl);
            }

            // Manually add the API Gateway service docs
            SwaggerUiConfigProperties.SwaggerUrl gatewayUrl = new SwaggerUiConfigProperties.SwaggerUrl();
            gatewayUrl.setName("gateway-service");
            gatewayUrl.setUrl(GATEWAY_URL + "/v3/api-docs");
            urls.add(gatewayUrl);

            swaggerUiConfigProperties.setUrls(urls);
            lastKnownServices = currentServices;
        } finally {
            lock.unlock();
        }
    }

    private String resolveApiDocsUrl(String serviceName) {
        System.out.println("Resolving API Docs URL for service: " + serviceName);
        return switch (serviceName.toLowerCase()) {
            case "user" -> GATEWAY_URL + "/user/swagger/v1/swagger.json"; // .NET service
            case "flora-upstream" -> GATEWAY_URL + "/flora-upstream/swagger/v1/swagger.json"; // NestJS service
            case "gene-bank" -> GATEWAY_URL + "/gene-bank/swagger/v1/swagger.json"; // Golang service
            case "flora-downstream" -> GATEWAY_URL + "/flora-downstream/swagger/v1//openapi.json"; // Spring service
            case "notification" -> GATEWAY_URL + "/notifications/swagger/v1/swagger.json"; // .NET service
            case "error-handler" -> GATEWAY_URL + "/error-handler/swagger/v1/swagger.json"; // Golang service
            default -> GATEWAY_URL + "/" + serviceName + "/v3/api-docs"; // Spring services
        };
    }

    private String removeServiceSuffix(String serviceName) {
        if (serviceName.toLowerCase().endsWith("-service")) {
            return serviceName.substring(0, serviceName.length() - 8); // Remove "-service"
        }
        return serviceName;
    }
}
