package fetch

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
)

func UInt16(url string) ([]uint16, error) {
	data, err := getAssetFile(url)
	if err != nil {
		return nil, err
	}

	var parsed []uint16

	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, &parsed)
	if err != nil {
		return nil, fmt.Errorf("failed converting binary to float array: %+v", err)
	}

	return parsed, nil
}

func Float32(url string) ([]float32, error) {
	data, err := getAssetFile(url)
	if err != nil {
		return nil, err
	}

	var parsed []float32

	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, &parsed)
	if err != nil {
		return nil, fmt.Errorf("failed converting binary to float array: %+v", err)
	}

	return parsed, nil
}

func Png(url string) ([]uint8, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset: %+v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("failed to fetch asset: status code %d", resp.StatusCode)
	}

	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image body: %+v", err)
	}

	if img, ok := img.(*image.NRGBA); ok {
		log.Printf("image data loaded: pixels %d (expected %d)\n", len(img.Pix), 256*256*4)
		return img.Pix, nil
	}

	return nil, fmt.Errorf("image not in rgba format: %T", img)
}

func getAssetFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset %s: %+v", url, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("failed to fetch asset %s: status code %d", url, resp.StatusCode)
	}

	var data []byte

	_, err = resp.Body.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset %s: failed to read request body: %+v", url, err)
	}

	return data, nil
}
