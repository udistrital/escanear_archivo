package models

type VerificarRequest struct {
	PdfBase64 string `json:"pdf_base64"`
	Firma     string `json:"firma"`
	UrlFileUp string `json:"urlFileUp"`
}

type RequestResponse struct {
	Status    string `json:"status"`     // "clean", "infected", o "error"
	RawOutput string `json:"raw_output"` 
}
