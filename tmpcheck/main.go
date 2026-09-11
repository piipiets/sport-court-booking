package main

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/piipiets/sport-court-booking/configs"
	"github.com/piipiets/sport-court-booking/storage"
)

func main() {
	configs.Initiator()

	s := storage.NewStorageClient(
		configs.SupabaseStorageURL(),
		configs.SupabaseSecretKey(),
		configs.SupabaseStorageBucket(),
	)

	png, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==")

	paths := []string{"courts/999/.jpg", "courts/999/court.jpg"}
	for _, p := range paths {
		url, err := s.Upload(p, bytes.NewReader(png), "image/jpeg")
		if err != nil {
			fmt.Printf("UPLOAD %-20s => FAIL: %v\n", p, err)
			continue
		}
		fmt.Printf("UPLOAD %-20s => OK: %s\n", p, url)

		data, ct, err := s.Download(url)
		if err != nil {
			fmt.Printf("DL     %-20s => FAIL: %v\n", p, err)
		} else {
			fmt.Printf("DL     %-20s => OK: %d bytes, %s\n", p, len(data), ct)
		}

		if err := s.Delete(url); err != nil {
			fmt.Printf("DELETE %-20s => FAIL: %v\n", p, err)
		} else {
			fmt.Printf("DELETE %-20s => OK\n", p)
		}
	}
}