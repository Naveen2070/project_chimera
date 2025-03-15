package com.naveen_r_sam.gateway_service.config;

import jakarta.annotation.PostConstruct;
import org.springdoc.core.properties.SwaggerUiConfigProperties;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.cloud.client.ServiceInstance;
import org.springframework.context.annotation.Configuration;

import java.util.LinkedHashSet;
import java.util.List;
import java.util.Set;

@Configuration
public class SwaggerUiConfig {

    private final DiscoveryClient discoveryClient;
    private final SwaggerUiConfigProperties swaggerUiConfigProperties;

    public SwaggerUiConfig(DiscoveryClient discoveryClient,
                           SwaggerUiConfigProperties swaggerUiConfigProperties) {
        this.discoveryClient = discoveryClient;
        this.swaggerUiConfigProperties = swaggerUiConfigProperties;
    }

    @PostConstruct
    public void configureSwaggerUi() {
        Set<SwaggerUiConfigProperties.SwaggerUrl> urls = new LinkedHashSet<>();

        List<String> services = discoveryClient.getServices();

        for (String serviceName : services) {
            System.out.println("Service name: " + serviceName);
            if (serviceName.equalsIgnoreCase("consul")) {
                continue;
            }

            List<ServiceInstance> instances = discoveryClient.getInstances(serviceName);

            for (ServiceInstance instance : instances) {
                String apiDocsUrl = resolveApiDocsUrl(serviceName, instance);

                SwaggerUiConfigProperties.SwaggerUrl swaggerUrl = new SwaggerUiConfigProperties.SwaggerUrl();
                swaggerUrl.setName(serviceName + " - " + instance.getHost() + ":" + instance.getPort());
                swaggerUrl.setUrl(apiDocsUrl);

                urls.add(swaggerUrl);
            }
        }

        // ✅ Manually add the API Gateway service docs
        SwaggerUiConfigProperties.SwaggerUrl gatewayUrl = new SwaggerUiConfigProperties.SwaggerUrl();
        gatewayUrl.setName("gateway-service - localhost:8080");
        gatewayUrl.setUrl("http://localhost:8080/v3/api-docs");
        urls.add(gatewayUrl);

        swaggerUiConfigProperties.setUrls(urls);
    }

    private String resolveApiDocsUrl(String serviceName, ServiceInstance instance) {
        String baseUri = instance.getUri().toString();

        return switch (serviceName.toLowerCase()) {
            case "user-service" -> baseUri + "/swagger/v1/swagger.json"; // C# service path
            default -> baseUri + "/v3/api-docs"; // Spring services
        };
    }
}
