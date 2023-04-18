package hosting

// service that can be run in the background, started and stopped
type BackgroundService interface {
	Start()
	Stop()
	GetName() string
}
