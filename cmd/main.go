package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	MaxFileSize = 6 * 1024 * 1024 // 6 MB
	ScanTimeout = 120 * time.Second
)

type RequestBody struct {
	PDFBase64 string `json:"pdf_base64"`
}

type LambdaResponse struct {
	Status    string `json:"status"`     // "clean", "infected", o "error"
	RawOutput string `json:"raw_output"` // salida de ClamAV
}

func main() {
	lambda.Start(handler)
}

/*func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var rawBody string
	
		// 1. Decodificar si el body está en base64
	if event.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(event.Body)
		if err != nil {
			return buildResponse(400, LambdaResponse{
				Status:    "error",
				RawOutput: "Cuerpo base64 inválido",
			}), nil
		}
		rawBody = string(decoded)
	} else {
		rawBody = event.Body
	}
	
	var reqBody RequestBody
	if err := json.Unmarshal([]byte(rawBody), &reqBody); err != nil {
		return buildResponse(400, LambdaResponse{
			Status:    "error",
			RawOutput: "Cuerpo JSON inválido",
		}), nil
	}
	
	if reqBody.PDFBase64 == "" {
		return buildResponse(400, LambdaResponse{
			Status:    "error",
			RawOutput: "Falta el campo 'pdf_base64'",
		}), nil
	}

	pdfBytes, err := base64.StdEncoding.DecodeString(req.PDFBase64)
	if err != nil {
		return buildResponse(400, LambdaResponse{
			Status:    "error",
			RawOutput: "Base64 inválido",
		}), nil
	}

	if len(pdfBytes) > MaxFileSize {
		return buildResponse(413, LambdaResponse{
			Status:    "error",
			RawOutput: "Archivo demasiado grande (máx 6MB)",
		}), nil
	}

	tempFile, err := ioutil.TempFile("", "*.pdf")
	if err != nil {
		log.Printf("Error creando archivo temporal: %v", err)
		return buildResponse(500, LambdaResponse{
			Status:    "error",
			RawOutput: "Error creando archivo temporal",
		}), nil
	}
	defer safeRemove(tempFile.Name())

	if _, err := tempFile.Write(pdfBytes); err != nil {
		log.Printf("Error escribiendo archivo: %v", err)
		return buildResponse(500, LambdaResponse{
			Status:    "error",
			RawOutput: "Error al guardar archivo temporal",
		}), nil
	}
	tempFile.Close()

	ctxScan, cancel := context.WithTimeout(ctx, ScanTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctxScan, "clamscan", "--no-summary", tempFile.Name())
	output, err := cmd.CombinedOutput()

	if ctxScan.Err() == context.DeadlineExceeded {
		log.Println("Timeout del escaneo")
		return buildResponse(500, LambdaResponse{
			Status:    "error",
			RawOutput: "El escaneo tomó demasiado tiempo",
		}), nil
	}

	rawOutput := string(output)
	log.Println("Salida clamscan:", rawOutput)

	status := "error"
	switch cmd.ProcessState.ExitCode() {
	case 0:
		status = "clean"
	case 1:
		status = "infected"
	default:
		status = "error"
	}

	return buildResponse(200, LambdaResponse{
		Status:    status,
		RawOutput: rawOutput,
	}), nil
}*/
func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("event: ", event)
	log.Println("📥 Iniciando función Lambda para escaneo de PDF")

	var rawBody string

	if event.IsBase64Encoded {
		log.Println("🧩 Decodificando body base64 del evento")
		decoded, err := base64.StdEncoding.DecodeString(event.Body)
		if err != nil {
			log.Printf("❌ Error decodificando base64 del body: %v", err)
			return buildResponse(400, LambdaResponse{
				Status:    "error",
				RawOutput: "Cuerpo base64 inválido",
			}), nil
		}
		rawBody = string(decoded)
	} else {
		rawBody = event.Body
	}

	log.Println("📄 Parseando JSON del body")
	var reqBody RequestBody
	if err := json.Unmarshal([]byte(rawBody), &reqBody); err != nil {
		log.Printf("❌ Error parseando JSON: %v", err)
		return buildResponse(400, LambdaResponse{
			Status:    "error",
			RawOutput: "Cuerpo JSON inválido",
		}), nil
	}

	if reqBody.PDFBase64 == "" {
		log.Println("⚠️ Falta el campo 'pdf_base64'")
		return buildResponse(400, LambdaResponse{
			Status:    "error",
			RawOutput: "Falta el campo 'pdf_base64'",
		}), nil
	}

	log.Println("🔍 Decodificando el PDF desde base64")
	pdfBytes, err := base64.StdEncoding.DecodeString(reqBody.PDFBase64)
	if err != nil {
		log.Printf("❌ Error decodificando PDF base64: %v", err)
		return buildResponse(400, LambdaResponse{
			Status:    "error",
			RawOutput: "Base64 inválido",
		}), nil
	}

	log.Printf("📦 Tamaño del archivo PDF: %.2f KB", float64(len(pdfBytes))/1024)
	if len(pdfBytes) > MaxFileSize {
		log.Printf("🚫 Archivo excede el tamaño máximo permitido (%d bytes)", MaxFileSize)
		return buildResponse(413, LambdaResponse{
			Status:    "error",
			RawOutput: "Archivo demasiado grande (máx 6MB)",
		}), nil
	}

	tempFile, err := ioutil.TempFile("", "*.pdf")
	if err != nil {
		log.Printf("❌ Error creando archivo temporal: %v", err)
		return buildResponse(500, LambdaResponse{
			Status:    "error",
			RawOutput: "Error creando archivo temporal",
		}), nil
	}
	defer safeRemove(tempFile.Name())

	log.Printf("📁 Archivo temporal creado: %s", tempFile.Name())
	if _, err := tempFile.Write(pdfBytes); err != nil {
		log.Printf("❌ Error escribiendo archivo temporal: %v", err)
		return buildResponse(500, LambdaResponse{
			Status:    "error",
			RawOutput: "Error al guardar archivo temporal",
		}), nil
	}
	tempFile.Close()

	log.Println("🛡️ Ejecutando clamscan sobre el archivo temporal")
	ctxScan, cancel := context.WithTimeout(ctx, ScanTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctxScan, "clamscan", "--no-summary", tempFile.Name())
	output, err := cmd.CombinedOutput()

	if ctxScan.Err() == context.DeadlineExceeded {
		log.Println("⏰ Tiempo de espera agotado para clamscan")
		return buildResponse(500, LambdaResponse{
			Status:    "error",
			RawOutput: "El escaneo tomó demasiado tiempo",
		}), nil
	}

	rawOutput := string(output)
	log.Printf("📤 Resultado de clamscan:\n%s", rawOutput)

	status := "error"
	switch cmd.ProcessState.ExitCode() {
	case 0:
		status = "clean"
	case 1:
		status = "infected"
	default:
		status = "error"
	}
	log.Printf("✅ Resultado final del escaneo: %s", status)

	return buildResponse(200, LambdaResponse{
		Status:    status,
		RawOutput: rawOutput,
	}), nil
}


func buildResponse(statusCode int, body interface{}) events.APIGatewayProxyResponse {
	jsonBody, _ := json.Marshal(body)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonBody),
	}
}

func safeRemove(path string) {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("No se pudo eliminar %s: %v", path, err)
	}
}

