package modelos

import (
	"errors"
	"strings"
	"unicode"
)

func ValidarNombrePersona(nombre string) bool {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return false
	}

	tieneLetra := false
	for _, r := range nombre {
		switch {
		case unicode.IsDigit(r):
			return false
		case unicode.IsLetter(r):
			tieneLetra = true
		case unicode.IsSpace(r), r == '\'', r == '´', r == '`', r == '-', r == '.':
			continue
		default:
			return false
		}
	}

	return tieneLetra
}

func ValidarTelefono(telefono string) bool {
	telefono = strings.TrimSpace(telefono)
	if telefono == "" {
		return true
	}
	if len(telefono) > 10 {
		return false
	}
	for _, r := range telefono {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func ValidarIDCliente(id string) bool {
	return tienePrefijoID(id, "CLI-")
}

func ValidarIDProducto(id string) bool {
	return tienePrefijoID(id, "PROD-")
}

func ValidarIDPedido(id string) bool {
	return tienePrefijoID(id, "PED-")
}

func tienePrefijoID(id string, prefijo string) bool {
	id = strings.TrimSpace(id)
	return strings.HasPrefix(id, prefijo) && len(id) > len(prefijo)
}

func ErrorIDClienteInvalido() error {
	return errors.New("el id del cliente debe iniciar con CLI-")
}

func ErrorIDProductoInvalido() error {
	return errors.New("el id del producto debe iniciar con PROD-")
}

func ErrorIDPedidoInvalido() error {
	return errors.New("el id del pedido debe iniciar con PED-")
}
