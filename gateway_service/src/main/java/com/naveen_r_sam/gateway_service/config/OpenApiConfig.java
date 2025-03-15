package com.naveen_r_sam.gateway_service.config;

import io.swagger.v3.oas.models.ExternalDocumentation;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.info.License;
import io.swagger.v3.oas.models.OpenAPI;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI gatewayOpenAPI() {
        return new OpenAPI()
                .info(new Info()
                        .title("Chimera Gateway Service")
                        .description("This is the documentation for the Chimera Gateway Service")
                        .version("1.0.0")
//                        .termsOfService("https://naveenr.com/terms")
                        .contact(new Contact()
                                .name("Naveen R")
                                .url("https://naveen2070.github.io/portfolio")
                                .email("naveenrameshcud@gmail.com"))
                        .license(new License()
                                .name("Apache 2.0")
                                .url("https://www.apache.org/licenses/LICENSE-2.0.html"))
                );
//                .externalDocs(new ExternalDocumentation()
//                        .description("Chimera Gateway Wiki Documentation")
//                        .url("https://chimera-docs.com"));
    }
}

