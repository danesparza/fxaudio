package media

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/danesparza/fxaudio/internal/data"
)

// mockAudioService is a mock implementation of the AudioService interface for testing
type mockAudioService struct {
	checkForPlayerCalled bool
	playAudioCalled      bool
	playAudioCtx         context.Context
	playAudioLoop        bool
	playAudioPath        string
	playAudioError       error
	playAudioDelay       time.Duration // Delay before returning from PlayAudio
}

func (m *mockAudioService) CheckForPlayer() error {
	m.checkForPlayerCalled = true
	return nil
}

func (m *mockAudioService) PlayAudio(ctx context.Context, loop bool, audioPathOrUrl string) error {
	m.playAudioCalled = true
	m.playAudioCtx = ctx
	m.playAudioLoop = loop
	m.playAudioPath = audioPathOrUrl

	// If a delay is set, wait for the specified time or until the context is cancelled
	if m.playAudioDelay > 0 {
		select {
		case <-time.After(m.playAudioDelay):
			// Delay completed
		case <-ctx.Done():
			// Context cancelled
			return ctx.Err()
		}
	}

	return m.playAudioError
}

// mockDataService is a mock implementation of the data.AppDataService interface for testing
type mockDataService struct{}

func (m *mockDataService) AddFile(ctx context.Context, filepath, description string) (data.File, error) {
	return data.File{
		ID:          "mock-id",
		Created:     0,
		FilePath:    filepath,
		Description: description,
	}, nil
}

func (m *mockDataService) GetFile(ctx context.Context, id string) (data.File, error) {
	return data.File{
		ID:          id,
		Created:     0,
		FilePath:    "mock-filepath",
		Description: "mock-description",
	}, nil
}

func (m *mockDataService) GetAllFiles(ctx context.Context) ([]data.File, error) {
	return []data.File{}, nil
}

func (m *mockDataService) GetAllFilesWithTag(ctx context.Context, tag string) ([]data.File, error) {
	return []data.File{}, nil
}

func (m *mockDataService) DeleteFile(ctx context.Context, id string) error {
	return nil
}

func (m *mockDataService) UpdateTags(ctx context.Context, id string, tags []string) error {
	return nil
}

func (m *mockDataService) GetConfig(ctx context.Context) (data.SystemConfig, error) {
	return data.SystemConfig{
		AlsaDevice: "mock-device",
	}, nil
}

func (m *mockDataService) SetConfig(ctx context.Context, config data.SystemConfig) error {
	return nil
}

func TestBackgroundProcess_HandleAndProcess(t *testing.T) {
	// Create a mock audio service
	mockAS := &mockAudioService{
		playAudioDelay: 100 * time.Millisecond, // Short delay to simulate audio playback
	}

	// Create a mock data service
	mockDB := &mockDataService{}

	// Create a BackgroundProcess with the mock services
	bp := &BackgroundProcess{
		DB:           mockDB,
		AS:           mockAS,
		PlayAudio:    make(chan PlayAudioRequest),
		StopAudio:    make(chan string),
		StopAllAudio: make(chan bool),
	}

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the HandleAndProcess method in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		bp.HandleAndProcess(ctx)
	}()

	// Test playing audio
	processID := "test-process-1"
	filePath := "test-file.mp3"
	bp.PlayAudio <- PlayAudioRequest{
		ProcessID: processID,
		FilePath:  filePath,
		Loop:      true,
	}

	// Give some time for the request to be processed
	time.Sleep(50 * time.Millisecond)

	// Check that PlayAudio was called with the correct parameters
	if !mockAS.playAudioCalled {
		t.Error("PlayAudio was not called")
	}
	if mockAS.playAudioPath != filePath {
		t.Errorf("PlayAudio was called with path %q, expected %q", mockAS.playAudioPath, filePath)
	}
	if !mockAS.playAudioLoop {
		t.Error("PlayAudio was called with loop=false, expected true")
	}

	// Test stopping audio
	bp.StopAudio <- processID

	// Give some time for the request to be processed
	time.Sleep(50 * time.Millisecond)

	// Check that the process was removed from the map
	bp.PlayingAudio.rwMutex.RLock()
	_, exists := bp.PlayingAudio.m[processID]
	bp.PlayingAudio.rwMutex.RUnlock()
	if exists {
		t.Error("Process was not removed from the map after StopAudio")
	}

	// Test playing multiple audio files
	processID1 := "test-process-2"
	processID2 := "test-process-3"
	bp.PlayAudio <- PlayAudioRequest{
		ProcessID: processID1,
		FilePath:  filePath,
		Loop:      true,
	}
	bp.PlayAudio <- PlayAudioRequest{
		ProcessID: processID2,
		FilePath:  filePath,
		Loop:      true,
	}

	// Give some time for the requests to be processed
	time.Sleep(50 * time.Millisecond)

	// Test stopping all audio
	bp.StopAllAudio <- true

	// Give some time for the request to be processed
	time.Sleep(50 * time.Millisecond)

	// Check that all processes were removed from the map
	bp.PlayingAudio.rwMutex.RLock()
	mapSize := len(bp.PlayingAudio.m)
	bp.PlayingAudio.rwMutex.RUnlock()
	if mapSize != 0 {
		t.Errorf("Map contains %d processes after StopAllAudio, expected 0", mapSize)
	}

	// Cancel the context to stop the HandleAndProcess method
	cancel()

	// Wait for the HandleAndProcess method to return
	wg.Wait()
}
