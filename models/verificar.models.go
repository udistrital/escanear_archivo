package models

type VerificarRequest struct {
	PdfBase64 string `json:"pdf_base64"`
	Firma     string `json:"firma"`
	UrlFileUp string `json:"urlFileUp"`
}