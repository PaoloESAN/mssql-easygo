package mssql

import (
	"fmt"
)

func (c *Conexion) CrearBD(nombre string) error {
	if c.conn == nil {
		return fmt.Errorf("no hay conexión activa")
	}

	query := fmt.Sprintf("CREATE DATABASE %s", nombre)
	_, err := c.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("error creando base de datos: %v", err)
	}

	fmt.Printf("Base de datos '%s' creada exitosamente\n", nombre)
	return nil
}

func (c *Conexion) CambiarBD(nombreBD string) error {
	if c.conn == nil {
		return fmt.Errorf("no hay conexión activa")
	}

	c.Database = nombreBD
	c.conn.Close()

	err := c.Connect()
	if err != nil {
		return err
	}

	return nil
}
