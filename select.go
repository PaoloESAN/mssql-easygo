package mssql

import (
	"fmt"
)

func (c *Conexion) Select(table string, print ...bool) ([]map[string]interface{}, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("no hay conexión activa")
	}

	shouldPrint := false
	if len(print) > 0 {
		shouldPrint = print[0]
	}

	query := fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := c.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo columnas: %v", err)
	}

	if shouldPrint {
		for i, col := range columns {
			if i > 0 {
				fmt.Print(" | ")
			}
			fmt.Print(col)
		}
		fmt.Println()
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

		if shouldPrint {
			for i, col := range columns {
				if i > 0 {
					fmt.Print(" | ")
				}
				fmt.Print(entry[col])
			}
			fmt.Println()
		}
	}

	return results, nil
}
