package modelos

import (
	"errors"
	"strings"
	"time"
)

type Cliente struct {
	id            string
	nombre        string
	email         string
	telefono      string
	fechaRegistro string
}

func NuevoCliente(nombre string, email string, telefono string) (Cliente, error) {
	return ReconstruirCliente(generarID("CLI"), nombre, email, telefono, time.Now().Format(time.RFC3339))
}

func ReconstruirCliente(id string, nombre string, email string, telefono string, fechaRegistro string) (Cliente, error) {
	id = strings.TrimSpace(id)
	nombre = strings.TrimSpace(nombre)
	email = strings.TrimSpace(strings.ToLower(email))
	telefono = strings.TrimSpace(telefono)
	fechaRegistro = strings.TrimSpace(fechaRegistro)

	if id == "" {
		return Cliente{}, errors.New("el id del cliente es obligatorio")
	}
	if nombre == "" {
		return Cliente{}, errors.New("el nombre del cliente es obligatorio")
	}
	if !ValidarEmail(email) {
		return Cliente{}, errors.New("el email del cliente no tiene un formato válido")
	}
	if fechaRegistro == "" {
		return Cliente{}, errors.New("la fecha de registro es obligatoria")
	}

	return Cliente{
		id:            id,
		nombre:        nombre,
		email:         email,
		telefono:      telefono,
		fechaRegistro: fechaRegistro,
	}, nil
}

func ValidarEmail(email string) bool {
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") && strings.Contains(email, ".") && !strings.Contains(email, " ")
}

func (c Cliente) ID() string            { return c.id }
func (c Cliente) Nombre() string        { return c.nombre }
func (c Cliente) Email() string         { return c.email }
func (c Cliente) Telefono() string      { return c.telefono }
func (c Cliente) FechaRegistro() string { return c.fechaRegistro }
