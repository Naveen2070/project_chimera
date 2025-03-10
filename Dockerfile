# Use the official HashiCorp Consul image
FROM hashicorp/consul:latest

# Expose the default port for Consul UI
EXPOSE 8500

# Start Consul in development mode (default)
CMD ["consul", "agent", "-dev", "-client=0.0.0.0"]