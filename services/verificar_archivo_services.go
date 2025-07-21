package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	MaxFileSize = 6 * 1024 * 1024 // 6 MB
	ScanTimeout = 120 * time.Second
)

type LambdaResponse struct {
	Status    string `json:"status"`     // "clean", "infected", o "error"
	RawOutput string `json:"raw_output"` // salida de ClamAV
}

func VerificarArchivo(pdfBase64 string) (*LambdaResponse, error) {
	log.Println("🔍 Decodificando base64...")

	pdfBytes, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		log.Printf("❌ Error decodificando base64: %v", err)
		return &LambdaResponse{
			Status:    "error",
			RawOutput: "Base64 inválido",
		}, nil
	}

	log.Printf("📦 Tamaño del PDF: %.2f KB", float64(len(pdfBytes))/1024)
	if len(pdfBytes) > MaxFileSize {
		log.Printf("🚫 Excede el tamaño máximo (%d bytes)", MaxFileSize)
		return &LambdaResponse{
			Status:    "error",
			RawOutput: "Archivo demasiado grande (máx 6MB)",
		}, nil
	}

	tempFile, err := ioutil.TempFile("", "*.pdf")
	if err != nil {
		log.Printf("❌ Error creando archivo temporal: %v", err)
		return nil, errors.New("error creando archivo temporal")
	}
	defer safeRemove(tempFile.Name())

	if _, err := tempFile.Write(pdfBytes); err != nil {
		log.Printf("❌ Error escribiendo archivo temporal: %v", err)
		return nil, errors.New("error escribiendo archivo temporal")
	}
	tempFile.Close()

	// Ejecutar clamscan con timeout
	cmd := exec.Command("timeout", fmt.Sprintf("%ds", int(ScanTimeout.Seconds())), "clamscan", "--no-summary", tempFile.Name())
	output, err := cmd.CombinedOutput()
	log.Println("📤 Resultado de clamscan:\n", string(output))

	status := "error"
	switch cmd.ProcessState.ExitCode() {
	case 0:
		status = "clean"
	case 1:
		status = "infected"
	default:
		status = "error"
	}
	log.Printf("✅ Resultado final: %s", status)

	return &LambdaResponse{
		Status:    status,
		RawOutput: string(output),
	}, nil
}

func safeRemove(path string) {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("⚠️ No se pudo eliminar %s: %v", path, err)
	}
}
