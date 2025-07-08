import base64
import tempfile
import subprocess
import os
import json
import re
import logging

# Configura logging
logger = logging.getLogger()
logger.setLevel(logging.INFO)

MAX_FILE_SIZE = 6 * 1024 * 1024  # 6 MB

def lambda_handler(event, context):
    try:
        encoded_pdf = event.get("pdf_base64")

        if not encoded_pdf:
            logger.warning("Campo 'pdf_base64' ausente en la petición.")
            return _response(400, {"error": "Falta el campo 'pdf_base64'"})

        try:
            pdf_bytes = base64.b64decode(encoded_pdf, validate=True)
        except Exception as e:
            logger.error(f"Base64 inválido: {str(e)}")
            return _response(400, {"error": "Base64 inválido"})

        if len(pdf_bytes) > MAX_FILE_SIZE:
            logger.warning("Archivo demasiado grande")
            return _response(413, {"error": "Archivo demasiado grande (máx 6MB)"})

        # Guardar temporalmente el archivo con permisos seguros
        with tempfile.NamedTemporaryFile(delete=False, suffix=".pdf", mode='wb') as tmp:
            tmp.write(pdf_bytes)
            file_path = tmp.name

        logger.info(f"Archivo temporal creado: {file_path}")

        try:
            result = subprocess.run(
                ["clamscan", "--no-summary", file_path],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                timeout=20  # Protección contra cuelgues
            )
        except subprocess.TimeoutExpired:
            logger.error("El escaneo tomó demasiado tiempo.")
            _safe_remove(file_path)
            return _response(500, {"error": "El escaneo tomó demasiado tiempo"})

        _safe_remove(file_path)

        output = result.stdout.decode(errors='replace')
        error_output = result.stderr.decode(errors='replace')

        logger.info("Salida clamscan: %s", output.strip())

        # Procesar resultado
        if result.returncode == 1 and "FOUND" in output:
            virus_name = _extract_virus_name(output)
            return _response(200, {
                "status": "infected",
                "virus_name": virus_name,
                "raw_output": output
            })
        elif result.returncode == 0:
            return _response(200, {
                "status": "clean",
                "raw_output": output
            })
        else:
            return _response(500, {
                "status": "error",
                "details": output,
                "stderr": error_output
            })

    except Exception as e:
        logger.exception("Error inesperado:")
        return _response(500, {"error": str(e)})


def _extract_virus_name(scan_output: str) -> str:
    """
    Extrae el nombre del virus desde la salida del comando clamscan.
    """
    match = re.search(r": ([^:]+) FOUND", scan_output)
    if match:
        return match.group(1).strip()
    return "Unknown"


def _safe_remove(path: str):
    """Elimina el archivo de forma segura."""
    try:
        if os.path.exists(path):
            os.remove(path)
    except Exception as e:
        logger.warning(f"No se pudo eliminar {path}: {e}")


def _response(status_code: int, body: dict):
    """Helper para construir respuestas."""
    return {
        "statusCode": status_code,
        "headers": {"Content-Type": "application/json"},
        "body": json.dumps(body)
    }
