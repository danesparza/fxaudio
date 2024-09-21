package media

import "context"

type AudioService interface {
	PlayAudio(ctx context.Context, loop bool, audioPathOrUrl string) error
}
