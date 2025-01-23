package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/phpdave11/gofpdf"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func generateRandomGarbageText(numBytes int) string {
	var sb strings.Builder
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for sb.Len() < numBytes {
		wl := rand.Intn(8) + 3
		for i := 0; i < wl; i++ {
			sb.WriteRune(letters[rand.Intn(len(letters))])
		}
		sb.WriteByte(' ')
	}
	return sb.String()
}

func main() {
	sizeMB := flag.Int("size", 5, "")
	outFile := flag.String("out", "out.pdf", "")
	flag.Parse()
	if *sizeMB <= 0 {
		log.Fatalf("Invalid -size: %d\n", *sizeMB)
	}
	targetBytes := int64(*sizeMB) * 1024 * 1024
	fmt.Println("Starting PDF generation...")
	log.Printf("Target: %d MB → %d bytes, Output: %s\n", *sizeMB, targetBytes, *outFile)
	rand.Seed(time.Now().UnixNano())
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetCompression(false)
	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()
	var buf bytes.Buffer
	const chunkSize = 2_000_000
	log.Printf("Using chunk size: %d bytes (~2 MB)\n", chunkSize)
	startTime := time.Now()
	chunkCount := 0
	for {
		chunkCount++
		text := generateRandomGarbageText(chunkSize)
		pdf.MultiCell(0, 5, text, "", "L", false)
		pdf.AddPage()
		buf.Reset()
		err := pdf.Output(&buf)
		if err != nil {
			log.Fatalf("Failed to render PDF: %v", err)
		}
		currentSize := int64(buf.Len())
		currentMB := float64(currentSize) / (1024.0 * 1024.0)
		elapsed := time.Since(startTime).Seconds()
		log.Printf("Chunk %d: PDF = %.2f MB (%.2fs elapsed)\n", chunkCount, currentMB, elapsed)
		if currentSize >= targetBytes {
			break
		}
	}
	err := os.WriteFile(*outFile, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Failed to write PDF: %v", err)
	}
	finalMB := float64(buf.Len()) / (1024.0 * 1024.0)
	duration := time.Since(startTime).Seconds()
	log.Printf("Done! Wrote %.2f MB to %s in %.2f seconds.\n", finalMB, *outFile, duration)
}
