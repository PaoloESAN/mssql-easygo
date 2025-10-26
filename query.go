package mssql

import "fmt"

// SqlSelect ejecuta consultas que retornan datos (SELECT)
// Los parámetros se pasan como @p1, @p2, @p3, etc. en la query para evitar inyecciones SQL
func (c *Conexion) SqlSelect(query string, args ...interface{}) ([]map[string]interface{}, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("no hay conexión activa")
	}

	rows, err := c.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo columnas: %v", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, fmt.Errorf("error escaneando fila: %v", err)
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}

		results = append(results, entry)
	}

	return results, nil
}

// SqlExec ejecuta consultas que no retornan datos (INSERT, UPDATE, DELETE)
// Los parámetros se pasan como @p1, @p2, @p3, etc. en la query para evitar inyecciones SQL
func (c *Conexion) SqlExec(query string, args ...interface{}) (int64, error) {
	if c.conn == nil {
		return 0, fmt.Errorf("no hay conexión activa")
	}

	result, err := c.conn.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("error ejecutando query: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error obteniendo filas afectadas: %v", err)
	}

	fmt.Printf("Filas afectadas: %d\n", rowsAffected)
	return rowsAffected, nil
}
