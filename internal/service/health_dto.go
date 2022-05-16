package service

type HealthStatus string

const (
	HealthStatusUp   HealthStatus = "up"
	HealthStatusDown HealthStatus = "down"
)

type Health struct {
	Status HealthStatus `json:"status"`
}
