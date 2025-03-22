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
package com.naveen_r_sam.auth_service.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.info.License;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.List;

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
        ).servers(
               List.of(
                       new io.swagger.v3.oas.models.servers.Server()
                               .url("http://localhost:8081")
                               .description("Local Server"),
                       new io.swagger.v3.oas.models.servers.Server()
                               .url("http://localhost:8080/auth")
                               .description("Gateway Server")
               )
        );
    }
}

