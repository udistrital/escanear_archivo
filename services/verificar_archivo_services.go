/*
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
*/

/*  clamscan clamscan  */

/*
package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	MaxFileSize   = 6 * 1024 * 1024 // 6 MB
	ScanTimeout   = 120 * time.Second
	TempScanDir   = "./files"
	ClamscanCheck = "clamscan"
)

type LambdaResponse struct {
	Status    string `json:"status"`     // "clean", "infected", o "error"
	RawOutput string `json:"raw_output"` // salida de ClamAV
}

func VerificarArchivo(pdfBase64 string) (*LambdaResponse, error) {
	log.Println("🔍 Verificando si 'clamscan' está disponible...")

	if _, err := exec.LookPath(ClamscanCheck); err != nil {
		log.Printf("❌ 'clamscan' no está instalado o no está en el PATH: %v", err)
		return &LambdaResponse{
			Status:    "error",
			RawOutput: "'clamscan' no disponible en el entorno",
		}, nil
	}

	log.Println("✅ 'clamscan' disponible.")
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

	// Crear carpeta ./files si no existe
	if err := os.MkdirAll(TempScanDir, 0755); err != nil {
		log.Printf("❌ No se pudo crear directorio %s: %v", TempScanDir, err)
		return nil, errors.New("error creando directorio de archivos")
	}

	// Crear archivo temporal en ./files con nombre único
	timestamp := time.Now().UnixNano()
	tempFilePath := filepath.Join(TempScanDir, fmt.Sprintf("scan_%d.pdf", timestamp))
	if err := os.WriteFile(tempFilePath, pdfBytes, 0644); err != nil {
		log.Printf("❌ Error escribiendo archivo %s: %v", tempFilePath, err)
		return nil, errors.New("error escribiendo archivo temporal")
	}
	log.Printf("📁 Archivo temporal guardado en: %s", tempFilePath)

	// Ejecutar clamscan con timeout
	cmd := exec.Command("timeout", fmt.Sprintf("%ds", int(ScanTimeout.Seconds())), "clamscan", "--no-summary", tempFilePath)
	output, err := cmd.CombinedOutput()

	log.Println("📤 Resultado de clamscan:\n", string(output))

	exitCode := cmd.ProcessState.ExitCode()
	log.Printf("📊 Exit code de clamscan: %d", exitCode)

	status := "error"
	switch exitCode {
	case 0:
		status = "clean"
	case 1:
		status = "infected"
	default:
		status = "error"
		if err != nil {
			log.Printf("⚠️ Error al ejecutar clamscan: %v", err)
		}
	}

	log.Printf("✅ Resultado final: %s", status)

	return &LambdaResponse{
		Status:    status,
		RawOutput: string(output),
	}, nil
}
*/

package services

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	MaxFileSize     = 6 * 1024 * 1024 // 6 MB
	TempScanDir     = "./files"
	ClamAVScanCmd   = "clamdscan"     // ahora usamos clamdscan
)

type LambdaResponse struct {
	Status    string `json:"status"`     // "clean", "infected", o "error"
	RawOutput string `json:"raw_output"` // salida de ClamAV
}

func VerificarArchivo(pdfBase64 string) (*LambdaResponse, error) {
	log.Println("🔍 Verificando si 'clamdscan' está disponible...")

	if _, err := exec.LookPath(ClamAVScanCmd); err != nil {
		log.Printf("❌ 'clamdscan' no está instalado o no está en el PATH: %v", err)
		return &LambdaResponse{
			Status:    "error",
			RawOutput: "'clamdscan' no disponible en el entorno",
		}, nil
	}

	log.Println("✅ 'clamdscan' disponible.")
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

	if err := os.MkdirAll(TempScanDir, 0755); err != nil {
		log.Printf("❌ No se pudo crear directorio %s: %v", TempScanDir, err)
		return nil, errors.New("error creando directorio de archivos")
	}

	timestamp := time.Now().UnixNano()
	tempFilePath := filepath.Join(TempScanDir, fmt.Sprintf("scan_%d.pdf", timestamp))

	if err := os.WriteFile(tempFilePath, pdfBytes, 0644); err != nil {
		log.Printf("❌ Error escribiendo archivo %s: %v", tempFilePath, err)
		return nil, errors.New("error escribiendo archivo temporal")
	}
	log.Printf("📁 Archivo temporal guardado en: %s", tempFilePath)

	// Limpieza del archivo después de escanear
	defer func() {
		if err := os.Remove(tempFilePath); err != nil {
			log.Printf("⚠️ No se pudo eliminar archivo temporal: %v", err)
		} else {
			log.Printf("🧹 Archivo temporal eliminado: %s", tempFilePath)
		}
	}()

	// Ejecutar clamdscan
	cmd := exec.Command(ClamAVScanCmd, "--no-summary", tempFilePath)
	output, err := cmd.CombinedOutput()

	log.Println("📤 Resultado de clamdscan:\n", string(output))

	exitCode := cmd.ProcessState.ExitCode()
	log.Printf("📊 Exit code de clamdscan: %d", exitCode)

	status := "error"
	switch exitCode {
	case 0:
		status = "clean"
	case 1:
		status = "infected"
	default:
		status = "error"
		if err != nil {
			log.Printf("⚠️ Error al ejecutar clamdscan: %v", err)
		}
	}

	log.Printf("✅ Resultado final: %s", status)

	return &LambdaResponse{
		Status:    status,
		RawOutput: string(output),
	}, nil
}
