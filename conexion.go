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

func NewConexion(server, user, password string, optional ...interface{}) *Conexion {
	db := "master"
	p := 1433

	if len(optional) > 0 {
		if dbVal, ok := optional[0].(string); ok && dbVal != "" {
			db = dbVal
		}
	}

	if len(optional) > 1 {
		if portVal, ok := optional[1].(int); ok {
			p = portVal
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

func (c *Conexion) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
