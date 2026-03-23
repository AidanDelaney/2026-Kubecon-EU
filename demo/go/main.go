package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	xdraw "golang.org/x/image/draw"
)

func overlayLogo(qrImg image.Image, logoPath string) (image.Image, error) {
	f, err := os.Open(logoPath)
	if err != nil {
		return nil, fmt.Errorf("opening logo: %w", err)
	}
	defer f.Close()

	logo, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decoding logo: %w", err)
	}

	qrSize := qrImg.Bounds().Dx()
	logoSize := qrSize / 5

	// Scale logo
	scaled := image.NewRGBA(image.Rect(0, 0, logoSize, logoSize))
	xdraw.BiLinear.Scale(scaled, scaled.Bounds(), logo, logo.Bounds(), xdraw.Over, nil)

	// Composite onto QR code
	out := image.NewRGBA(qrImg.Bounds())
	draw.Draw(out, out.Bounds(), qrImg, image.Point{}, draw.Src)

	offset := (qrSize - logoSize) / 2
	logoRect := image.Rect(offset, offset, offset+logoSize, offset+logoSize)
	draw.Draw(out, logoRect, scaled, image.Point{}, draw.Over)

	return out, nil
}

func main() {
	var logoPath string

	var rootCmd = &cobra.Command{
		Use:   "qr <url>",
		Short: "Generate a QR code for a URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			filename := uuid.New().String() + ".png"

			qr, err := qrcode.New(url, qrcode.High)
			if err != nil {
				return fmt.Errorf("creating QR code: %w", err)
			}

			img := qr.Image(512)

			if logoPath != "" {
				img, err = overlayLogo(img, logoPath)
				if err != nil {
					return err
				}
			}

			out, err := os.Create(filename)
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer out.Close()

			if err := png.Encode(out, img); err != nil {
				return fmt.Errorf("encoding PNG: %w", err)
			}

			fmt.Println("QR code generated and saved as", filename)
			return nil
		},
	}

	rootCmd.Flags().StringVar(&logoPath, "logo", "", "path to logo image to overlay in the center")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
