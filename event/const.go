package event

const (
	// SystemStartup event is when the system has started up
	SystemStartup = "System startup"

	// FileUploaded event is when the system has found a file to process
	FileUploaded = "File uploaded"

	// FilePlayed event is when the system has successfully played a file
	FilePlayed = "File played"

	// FileDeleted event is when the system has successfully deleted a file
	FileDeleted = "File deleted"

	// FileError event is when there was an error processing a file
	FileError = "File error"

	// SystemShutdown event is when the system is shutting down
	SystemShutdown = "System Shutdown"
)
