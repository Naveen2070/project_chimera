using Microsoft.Extensions.Options;
using Steeltoe.Discovery.Client;
using Steeltoe.Management.Endpoint;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.

builder.Services.AddControllers();
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// Steeltoe Consul
builder.Services.AddDiscoveryClient();

// Steeltoe Management Actuators
builder.Services.AddAllActuators();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

app.UseAuthorization();

app.MapAllActuators();

app.MapControllers();

app.Run();
