package services

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/udistrital/escanear_archivo/models"
)

const (
	MaxFileSize      = 6 * 1024 * 1024 // 6 MB
	defaultClamdAddr = "127.0.0.1:3310"
	clamdChunkSize   = 4096
	clamdTimeout     = 30 * time.Second
)

func clamdAddr() string {
	if addr := os.Getenv("CLAMD_ADDR"); addr != "" {
		return addr
	}
	return defaultClamdAddr
}

func VerificarArchivo(pdfBase64 string) (*models.RequestResponse, error) {

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

	output, err := scanWithClamd(pdfBytes)
	if err != nil {
		return &models.RequestResponse{
			Status:    "error",
			RawOutput: "No se pudo contactar clamd: " + err.Error(),
		}, nil
	}

	status := "error"
	switch {
	case strings.HasSuffix(output, "FOUND"):
		status = "infected"
	case strings.HasSuffix(output, "OK"):
		status = "clean"
	}

	return &models.RequestResponse{
		Status:    status,
		RawOutput: output,
	}, nil
}

// scanWithClamd envía los bytes al daemon clamd usando el protocolo INSTREAM
// sobre TCP y devuelve la línea de respuesta (p.ej. "stream: OK" o
// "stream: Eicar-Test-Signature FOUND").
func scanWithClamd(data []byte) (string, error) {
	conn, err := net.DialTimeout("tcp", clamdAddr(), clamdTimeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(clamdTimeout)); err != nil {
		return "", err
	}

	// Inicia la sesión de streaming.
	if _, err := conn.Write([]byte("zINSTREAM\x00")); err != nil {
		return "", err
	}

	// Envía el archivo en fragmentos con el formato <len uint32 BE><datos>.
	for i := 0; i < len(data); i += clamdChunkSize {
		end := i + clamdChunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]

		var sizeBuf [4]byte
		binary.BigEndian.PutUint32(sizeBuf[:], uint32(len(chunk)))
		if _, err := conn.Write(sizeBuf[:]); err != nil {
			return "", err
		}
		if _, err := conn.Write(chunk); err != nil {
			return "", err
		}
	}

	// Fragmento de longitud cero: marca el fin del stream.
	if _, err := conn.Write([]byte{0, 0, 0, 0}); err != nil {
		return "", err
	}

	// La respuesta termina en NUL.
	resp, err := bufio.NewReader(conn).ReadString(0x00)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(resp, "\x00\n"), nil
}

func PingClamd() error {
	conn, err := net.DialTimeout("tcp", clamdAddr(), clamdTimeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(clamdTimeout)); err != nil {
		return err
	}

	if _, err := conn.Write([]byte("zPING\x00")); err != nil {
		return err
	}

	buf := make([]byte, 16)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}
	if !strings.Contains(string(buf[:n]), "PONG") {
		return fmt.Errorf("respuesta inesperada de clamd: %q", string(buf[:n]))
	}
	return nil
}
