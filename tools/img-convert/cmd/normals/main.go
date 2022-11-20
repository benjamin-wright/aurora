package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := run(os.Args); err != nil {
		logrus.Fatalf("Error: %+v", err)
	}
}

func run(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("expected one argument, got %d", len(args))
	}

	inpath := args[1]
	outpath := args[2]

	data, err := getImageData(inpath)
	if err != nil {
		return fmt.Errorf("failed to read image data: %+v", err)
	}

	normals, err := convertToNormals(data)
	if err != nil {
		return fmt.Errorf("failed to convert normal data: %+v", err)
	}

	err = writeNormals(normals, outpath)
	if err != nil {
		return fmt.Errorf("failed to write normal data: %+v", err)
	}

	return nil
}

type ImageData struct {
	Pixels []uint8
	Width  int
	Height int
}

func getImageData(filepath string) (*ImageData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		logrus.Fatalf("Failed to open normals file: %+v", err)
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse image data: %+v", err)
	}

	if img, ok := img.(*image.NRGBA); ok {
		return &ImageData{
			Pixels: img.Pix,
			Width:  img.Rect.Dx(),
			Height: img.Rect.Dy(),
		}, nil
	}

	return nil, fmt.Errorf("image not in expected NRGBA format: %T", img)
}

func convertToNormals(data *ImageData) ([]float32, error) {
	size := data.Width * data.Height * 2
	output := make([]float32, size)

	stride := data.Width * 4
	outStride := data.Width * 2

	for y := 0; y < data.Height; y++ {
		offset := y * stride
		outOffset := y * outStride

		for x := 0; x < data.Width; x++ {
			output[x*2+outOffset] = (float32(data.Pixels[x*4+offset]) - 127.0) / 127.0
			output[x*2+outOffset+1] = (float32(data.Pixels[x*4+offset+1]) - 127.0) / 127.0
		}
	}

	return output, nil
}

func writeNormals(normals []float32, filepath string) error {
	var data bytes.Buffer
	err := binary.Write(&data, binary.LittleEndian, &normals)
	if err != nil {
		return fmt.Errorf("failed to convert data to binary: %+v", err)
	}

	err = os.WriteFile(filepath, data.Bytes(), 0777)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %+v", err)
	}

	return nil
}
