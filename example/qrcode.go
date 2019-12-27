package example

import (
	"github.com/skip2/go-qrcode"
	"image/color"
)

func NewQRCode(content, filepath string) error {
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return err
	}

	qr.ForegroundColor = color.White
	qr.BackgroundColor = color.RGBA{50, 205, 50, 255}
	return qr.WriteFile(256, filepath)
}
