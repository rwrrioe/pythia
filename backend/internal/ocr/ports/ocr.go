package ports

import "context"

type OCRProvider interface {
	ProcessImage(
		ctx context.Context,
		imageData []byte,
		lang string,
	) ([]string, error)
}
