package mssql

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

type Conexion struct {
	Server   string
	User     string
	Password string
	Database string
	Port     int
	conn     *sql.DB
}

func NewConexion(server, user, password string, port ...interface{}) *Conexion {
	p := 1433
	db := "master"

	if len(port) > 0 {
		if portVal, ok := port[0].(int); ok {
			p = portVal
		}
	}

	if len(port) > 1 {
		if dbVal, ok := port[1].(string); ok && dbVal != "" {
			db = dbVal
		}
	}

	c := &Conexion{
		Server:   server,
		User:     user,
		Password: password,
		Database: db,
		Port:     p,
	}

	err := c.Connect()
	if err != nil {
		fmt.Printf("Error en la conexión: %v\n", err)
	}

	return c
}

func (c *Conexion) Connect() error {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		c.Server, c.User, c.Password, c.Port, c.Database)

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		return fmt.Errorf("error abriendo base de datos: %v", err)
	}

	err = conn.Ping()
	if err != nil {
		return fmt.Errorf("error conectando a SQL Server: %v", err)
	}

	c.conn = conn
	fmt.Printf("Conexión a SQL Server exitosa - Base de datos: %s\n", c.Database)
	return nil
}

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

func (c *Conexion) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
