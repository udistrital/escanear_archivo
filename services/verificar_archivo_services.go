package services

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"github.com/udistrital/escanear_archivo/models"
)

const (
	MaxFileSize     = 6 * 1024 * 1024 // 6 MB
	TempScanDir     = "./files"
	ClamAVScanCmd   = "clamdscan"     
)

func VerificarArchivo(pdfBase64 string) (*models.RequestResponse, error) {

	if _, err := exec.LookPath(ClamAVScanCmd); err != nil {
		return &models.RequestResponse{
			Status:    "error",
			RawOutput: "'clamdscan' no disponible en el entorno",
		}, nil
	}

	pdfBytes, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		return &models.RequestResponse{
			Status:    "error",
			RawOutput: "Base64 inválido",
		}, nil
	}

	if len(pdfBytes) > MaxFileSize {
		return &models.RequestResponse{
			Status:    "error",
			RawOutput: "Archivo demasiado grande (máx 6MB)",
		}, nil
	}

	if err := os.MkdirAll(TempScanDir, 0755); err != nil {
		return &models.RequestResponse{
			Status:    "error",
			RawOutput: "Error creando directorio temporal",
		}, nil
	}

	timestamp := time.Now().UnixNano()
	tempFilePath := filepath.Join(TempScanDir, fmt.Sprintf("scan_%d.pdf", timestamp))

	if err := os.WriteFile(tempFilePath, pdfBytes, 0644); err != nil {

		return &models.RequestResponse{
			Status:    "error",
			RawOutput: "Error escribiendo archivo temporal",
		}, nil
	}

	defer func() {
		if err := os.Remove(tempFilePath); err != nil {
			time.Sleep(60 * time.Millisecond)
			if errRetry := os.Remove(tempFilePath); errRetry != nil {
				log.Printf("❌ Segundo intento fallido para eliminar archivo temporal: %v", errRetry)
			} 
		} 
	}()
	

	cmd := exec.Command(ClamAVScanCmd, "--no-summary", tempFilePath)
	output, err := cmd.CombinedOutput()

	exitCode := cmd.ProcessState.ExitCode()

	status := "error"
	switch exitCode {
	case 0:
		status = "clean"
	case 1:
		status = "infected"
	default:
		status = "error"
	}

	return &models.RequestResponse{
		Status:    status,
		RawOutput: string(output),
	}, nil
}
