package event

const (
	// SystemStartup event is when the system has started up
	SystemStartup = "System startup"

	// FileUploaded event is when the system has found a file to process
	FileUploaded = "File uploaded"

	// FileDeleted event is when the system has successfully deleted a file
	FileDeleted = "File deleted"

	// RequestPlay event is when the system processes a request to play audio
	RequestPlay = "Request play"

	// RequestStop event is when the system processes a request to stop audio
	RequestStop = "Request stop"

	// RequestStopAll event is when the system processes a request to stop all audio
	RequestStopAll = "Request stop all"

	// FileError event is when there was an error processing a file
	FileError = "File error"

	// SystemShutdown event is when the system is shutting down
	SystemShutdown = "System Shutdown"
)
