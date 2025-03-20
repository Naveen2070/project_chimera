package actuator

// ActuatorService defines the interface for actuator-related operations
type ActuatorService interface {
	GetHealthMessage() string
}

// actuatorService is the concrete implementation of ActuatorService
type actuatorService struct{}

// NewActuatorService creates a new instance of actuatorService
func NewActuatorService() ActuatorService {
	return &actuatorService{}
}

// GetHealthMessage returns the health status
func (s *actuatorService) GetHealthMessage() string {
	// Return the health status
	return "UP"
}
