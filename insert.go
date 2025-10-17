package mssql

import (
	"fmt"
)

func (c *Conexion) Insert(table string, values ...interface{}) error {
	if c.conn == nil {
		return fmt.Errorf("no hay conexión activa")
	}

	if len(values) == 0 {
		return fmt.Errorf("no hay valores para insertar")
	}

	query := fmt.Sprintf("SELECT TOP 0 * FROM %s", table)
	rows, err := c.conn.Query(query)
	if err != nil {
		return fmt.Errorf("error obteniendo columnas: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("error leyendo columnas: %v", err)
	}

	if len(values) != len(columns) {
		return fmt.Errorf("se esperaban %d valores pero se recibieron %d", len(columns), len(values))
	}

	var placeholders []string
	for i := 1; i <= len(values); i++ {
		placeholders = append(placeholders, fmt.Sprintf("@p%d", i))
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		joinStrings(columns, ", "),
		joinStrings(placeholders, ", "))

	_, err = c.conn.Exec(insertQuery, values...)
	if err != nil {
		return fmt.Errorf("error insertando datos: %v", err)
	}

	fmt.Printf("Datos insertados en '%s' exitosamente\n", table)
	return nil
}
