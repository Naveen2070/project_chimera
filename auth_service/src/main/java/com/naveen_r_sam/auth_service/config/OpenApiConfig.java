package com.naveen_r_sam.auth_service.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.info.License;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI customOpenAPI() {
        return new OpenAPI().info(
                new Info()
                        .title("Chimera Auth Service")
                        .version("1.0.0")
                        .description("This API provides authentication and authorization tokens for Chimera Services")
//                        .termsOfService("https://naveenr.com/terms")
                        .contact(new Contact()
                                .name("Naveen R")
                                .url("https://naveen2070.github.io/portfolio")
                                .email("naveenrameshcud@gmail.com"))
                        .license(new License()
                                .name("Apache 2.0")
                                .url("https://www.apache.org/licenses/LICENSE-2.0"))
        );
    }
}

