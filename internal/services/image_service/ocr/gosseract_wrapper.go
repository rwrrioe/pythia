package ocr

import (
	"bytes"
	"context"
	"os/exec"
)

type ImageProcesser interface {
	ProcessImage(ctx context.Context, imagePath string, lang string) error
}

type OCRProcessor struct{}

func NewOCRProcessor() *OCRProcessor {
	return &OCRProcessor{}
}

func (p *OCRProcessor) RecognizeText(imagePath string) (string, error) {
	cmd := exec.Command(`C:\Program Files\Tesseract-OCR\tesseract.exe`, imagePath, "stdout")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
