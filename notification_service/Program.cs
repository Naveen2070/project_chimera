using Consul;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Microsoft.Extensions.Diagnostics.HealthChecks;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();

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
