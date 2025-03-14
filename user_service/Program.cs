using Microsoft.EntityFrameworkCore;
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

// Add Database Context
builder.Services.AddDbContext<DatabaseContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection"))
   );

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
