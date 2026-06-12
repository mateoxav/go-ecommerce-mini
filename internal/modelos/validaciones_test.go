package modelos

import "testing"

func TestValidarNombrePersona(t *testing.T) {
	casosValidos := []string{"Ana María", "José Pérez", "Luis-Angel"}
	for _, caso := range casosValidos {
		if !ValidarNombrePersona(caso) {
			t.Fatalf("se esperaba nombre válido: %s", caso)
		}
	}

	casosInvalidos := []string{"", "Carlos2", "12345", "Ana_"}
	for _, caso := range casosInvalidos {
		if ValidarNombrePersona(caso) {
			t.Fatalf("se esperaba nombre inválido: %s", caso)
		}
	}
}

func TestValidarTelefono(t *testing.T) {
	casosValidos := []string{"", "0991234567", "123"}
	for _, caso := range casosValidos {
		if !ValidarTelefono(caso) {
			t.Fatalf("se esperaba teléfono válido: %s", caso)
		}
	}

	casosInvalidos := []string{"09912345678", "099-123", "telefono"}
	for _, caso := range casosInvalidos {
		if ValidarTelefono(caso) {
			t.Fatalf("se esperaba teléfono inválido: %s", caso)
		}
	}
}

func TestValidarPrefijosID(t *testing.T) {
	if !ValidarIDCliente("CLI-123") || !ValidarIDProducto("PROD-123") || !ValidarIDPedido("PED-123") {
		t.Fatal("se esperaban ids válidos con sus prefijos")
	}
	if ValidarIDPedido("PROD-123") || ValidarIDProducto("PED-123") || ValidarIDCliente("123") {
		t.Fatal("no se deben aceptar ids con prefijo incorrecto")
	}
}
