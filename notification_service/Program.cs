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
builder.Services.AddSwaggerGen();

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
    app.UseSwaggerUI();
}

// Get the lifetime of the application
var lifetime = app.Lifetime;
// Get the Consul client
var consulClient = app.Services.GetRequiredService<IConsulClient>();
// Generate a unique service ID
var serviceId = Guid.NewGuid().ToString();

// Register the service with Consul
var registration = new AgentServiceRegistration
{
    ID = serviceId,
    Name = "notification_service",
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
