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

func (c *Conexion) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
