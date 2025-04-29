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
using Microsoft.EntityFrameworkCore;
using Microsoft.OpenApi.Models;
using Steeltoe.Discovery.Client;
using Steeltoe.Management.Endpoint;
using user_service.Database;
using user_service.Services.Interfaces;
using user_service.Services;
using user_service.Utils;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddAutoMapper(typeof(Program));
builder.Services.AddSingleton(new PasswordEncoder(12));
builder.Services.AddScoped<IUserService, UserService>();
builder.Services.AddSingleton<IRMQPublisherService>(sp =>
{
    var configuration = sp.GetRequiredService<IConfiguration>();
    var service = RMQPublisherService.CreateAsync(configuration).GetAwaiter().GetResult();
    return service;
});
builder.Services.AddControllers();

// Add CORS configuration
builder.Services.AddCors(options =>
{
    options.AddPolicy(name:"AllowGateway", policy =>
    {
        policy.WithOrigins("http://localhost:8080")
              .AllowAnyHeader()
              .AllowAnyMethod()
              .AllowCredentials()  
              .WithExposedHeaders("Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods");
    });

    options.AddPolicy(name:"AllowAll", policy =>
    {
        policy.AllowAnyOrigin()
              .AllowAnyHeader()
              .AllowAnyMethod()
              .WithExposedHeaders("Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods");
    });
});


// Add Database Context
builder.Services.AddDbContext<DatabaseContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection"))
);

// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();

builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new OpenApiInfo
    {
        Title = "Chimera User Service", 
        Version = "1.0.0",
        Description = "This API provides user management capabilities for the Chimera IAM system.",
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
        Url = "http://localhost:8080/user", // Adjust to match your API Gateway route
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

// Steeltoe Consul
builder.Services.AddDiscoveryClient();

// Steeltoe Management Actuators
builder.Services.AddAllActuators();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/swagger/v1/swagger.json", "User Service API v1");
    });
}

app.UseHttpsRedirection();

app.UseCors("AllowGateway");

app.UseAuthorization();

app.MapAllActuators();

app.MapControllers();

app.Run();
