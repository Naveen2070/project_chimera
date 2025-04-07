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

using Consul;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Microsoft.Extensions.Diagnostics.HealthChecks;
using notification_service.Services.Interfaces;
using notification_service.Services;
using Microsoft.OpenApi.Models;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();

// Register RMQConsumerService as a singleton
builder.Services.AddSingleton<IRMQConsumerService>(provider =>
{
    var configuration = provider.GetRequiredService<IConfiguration>();
    return RMQConsumerService.CreateAsync(configuration).GetAwaiter().GetResult();
});

builder.Services.AddScoped<IFloraNotificationService, FloraNotificationService>();

// Add health checks
builder.Services.AddHealthChecks()
    .AddCheck("Self", () => HealthCheckResult.Healthy(), tags: new[] { "ready", "live" });

// Add Consul client
builder.Services.AddSingleton<IConsulClient>(p => new ConsulClient(config =>
{
    var consulHost = builder.Configuration["Consul:Host"];
    var consulPort = builder.Configuration["Consul:Port"];
    config.Address = new Uri($"{consulHost}:{consulPort}");
}));

// Add Swagger
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new OpenApiInfo
    {
        Title = "Chimera User Service",
        Version = "1.0.0",
        Description = "This API provides notification service for Chimera Application",
        Contact = new OpenApiContact
        {
            Name = "Naveen R",
            Email = "naveenrameshcud@gmail.com",
            Url = new Uri("https://naveen2070.github.io/portfolio")
        },
        License = new OpenApiLicense
        {
            Name = "Apache 2.0",
            Url = new Uri("https://www.apache.org/licenses/LICENSE-2.0")
        }
    });

    c.AddServer(new OpenApiServer
    {
        Url = "http://localhost:5035", // Replace with actual service URL
        Description = "Local User Service"
    });

    c.AddServer(new OpenApiServer
    {
        Url = "http://localhost:8080/notifications", // Adjust to match your API Gateway route
        Description = "API Gateway"
    });

    // Add JWT Authentication to Swagger 
    c.AddSecurityDefinition("Bearer", new OpenApiSecurityScheme
    {
        Name = "Authorization",
        Type = SecuritySchemeType.Http,
        Scheme = "Bearer",
        BearerFormat = "JWT",
        In = ParameterLocation.Header,
        Description = "Enter 'Bearer {token}' without quotes"
    });

    c.AddSecurityRequirement(new OpenApiSecurityRequirement
    {
        {
            new OpenApiSecurityScheme
            {
                Reference = new OpenApiReference
                {
                    Type = ReferenceType.SecurityScheme,
                    Id = "Bearer"
                }
            },
            new string[] {}
        }
    });

    // Optional: Enable annotations if you use [SwaggerOperation], etc.
    // c.EnableAnnotations();
});

// Add CORS policy
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowAll", builder =>
    {
        builder.AllowAnyOrigin()
               .AllowAnyMethod()
               .AllowAnyHeader();
    });
});


// Build the application
var app = builder.Build();

// Configure Swagger for development
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/swagger/v1/swagger.json", "User Service API v1");
    });
}

// Get the lifetime of the application
IHostApplicationLifetime lifetime = app.Lifetime;

// Get the Consul client
IConsulClient consulClient = app.Services.GetRequiredService<IConsulClient>();

// Generate a unique service ID
string guid = Guid.NewGuid().ToString();
string serviceId = "notification-service" + "-"  + guid.Split('-')[0]+ "-" + guid.Split('-')[2];

// Register the service with Consul
var registration = new AgentServiceRegistration
{
    ID = serviceId,
    Name = "notification-service",
    Address = "localhost",
    Port = 5249,
    Check = new AgentServiceCheck
    {
        HTTP = "http://host.docker.internal:5249/health/ready", 
        Interval = TimeSpan.FromSeconds(10),
        Timeout = TimeSpan.FromSeconds(5),
        DeregisterCriticalServiceAfter = TimeSpan.FromMinutes(1) 
    }
};

await consulClient.Agent.ServiceRegister(registration);

// Deregister the service when the application is shutting down
lifetime.ApplicationStopping.Register(async () =>
{
    await consulClient.Agent.ServiceDeregister(serviceId);
});

app.UseHttpsRedirection();

// Use CORS
app.UseCors("AllowAll");

app.UseAuthorization();
app.MapControllers();

// Enable health checks endpoint
app.MapHealthChecks("/health/ready", new HealthCheckOptions
{
    Predicate = check => check.Tags.Contains("ready"), 
});

app.MapHealthChecks("/health/live", new HealthCheckOptions
{
    Predicate = _ => true, 
});

app.Run();
