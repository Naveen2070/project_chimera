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
builder.Services.AddControllers();

// Add CORS configuration
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowGateway", policy =>
    {
        policy.WithOrigins("http://localhost:8080")
              .AllowAnyHeader()
              .AllowAnyMethod()
              .AllowCredentials()  
              .WithExposedHeaders("Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods");
    });

    options.AddPolicy("AllowAll", policy =>
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
        Title = "User Service API", 
        Version = "v1",
        Description = "API for managing users within the system.",
        Contact = new OpenApiContact
        {
            Name = "Naveen R",
            Email = "naveenrameshcud@gmail.com",
            Url = new Uri("https://naveen2070.github.io/portfolio")
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

app.Use(async (context, next) =>
{
    Console.WriteLine($"Origin Header: {context.Request.Headers.Origin}");
    await next();
});


// Use CORS policy
app.UseCors("AllowGateway");

app.UseHttpsRedirection();

app.UseAuthorization();

app.MapAllActuators();

app.MapControllers();

app.Run();
