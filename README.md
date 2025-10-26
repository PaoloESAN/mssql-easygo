# MSSQL EasyGo

Librería Go para conectar y gestionar bases de datos SQL Server de forma sencilla.

## Instalación

```bash
go get github.com/PaoloESAN/mssql-easygo
```

## Importación

Para usar la librería, debes importarla en tu código:

```go
import "github.com/PaoloESAN/mssql-easygo"
```

Luego puedes usar la librería:

```go
conn := mssql.NewConexion("localhost", "sa", "password")
conn := mssql.NewConexion("localhost", "sa", "password", "BaseDeDatos")
conn := mssql.NewConexion("localhost", "sa", "password", "BaseDeDatos", 1433)
```

---

## Funciones

### `NewConexion(server, user, password, database(opcional), port(opcional))`

Crea una nueva conexión a SQL Server. La conexión se realiza automáticamente.

**Parámetros:**
- `server` (string, requerido): Servidor SQL Server (ej: "localhost")
- `user` (string, requerido): Usuario (ej: "sa")
- `password` (string, requerido): Contraseña
- `database` (string, opcional): Nombre de la base de datos (por defecto "master")
- `port` (int, opcional): Puerto (por defecto 1433)

**Retorna:** `*Conexion`

**Ejemplo:**
```go
c := mssql.NewConexion("localhost", "sa", "password")
c := mssql.NewConexion("localhost", "sa", "password", "midb")
c := mssql.NewConexion("localhost", "sa", "password", "midb", 1433)
```

---

### `Select(table, print(opcional))`

Ejecuta un SELECT * en una tabla.

**Parámetros:**
- `table` (string, requerido): Nombre de la tabla
- `print` (bool, opcional): Si es `true`, imprime los resultados en consola (por defecto `false`)

**Retorna:** `[]map[string]interface{}, error`

**Ejemplo:**
```go
results, err := c.Select("usuarios")

results, err := c.Select("usuarios", true)
```

**Salida con `print=true`:**
```
id | nombre | email
1 | Juan | juan@email.com
2 | Maria | maria@email.com
```

---

### `Insert(table, valores...)`

Inserta valores en una tabla. Los valores deben ir en el mismo orden que las columnas de la tabla.

**Parámetros:**
- `table` (string, requerido): Nombre de la tabla donde insertar
- `valores...` (variadic interface{}, requerido): Valores a insertar en orden de columnas

**Retorna:** `error`

**Ejemplo:**
```go
// Insertar en tabla con 3 columnas: id, nombre, email
err := c.Insert("usuarios", 1, "Juan", "juan@email.com")
if err != nil {
    fmt.Printf("Error: %v\n", err)
}

// Insertar diferentes tipos de datos
err = c.Insert("productos", 1, "Laptop", 999.99, 10)
if err != nil {
    fmt.Printf("Error: %v\n", err)
}
```

**Características:**
- ✅ Soporta cualquier tipo de dato
- ✅ Protección contra SQL Injection (usa parámetros)
- ✅ Valida la cantidad de valores
- ✅ Sintaxis simple y directa

**Nota:** Los valores deben coincidir exactamente con la cantidad de columnas de la tabla y en el mismo orden.

---

### `CrearBD(nombre)`

Crea una nueva base de datos.

**Parámetros:**
- `nombre` (string, requerido): Nombre de la base de datos a crear

**Retorna:** `error`

**Ejemplo:**
```go
err := c.CrearBD("nuevadb")
if err != nil {
    fmt.Printf("Error: %v\n", err)
}
```

---

### `CambiarBD(nombreBD)`

Cambia la base de datos activa y reconecta.

**Parámetros:**
- `nombreBD` (string, requerido): Nombre de la base de datos a la que cambiar

**Retorna:** `error`

**Ejemplo:**
```go
err := c.CambiarBD("otradb")
if err != nil {
    fmt.Printf("Error: %v\n", err)
}
```

---

### `SqlSelect(query, args...)`

Ejecuta una consulta SQL personalizada con parámetros. Retorna los resultados como un slice de mapas.

**Parámetros:**
- `query` (string, requerido): Consulta SQL con placeholders (@p1, @p2, etc.)
- `args...` (variadic interface{}, opcional): Valores para reemplazar los placeholders

**Retorna:** `[]map[string]interface{}, error`

**Ejemplo:**
```go
// SELECT con parámetros
nombre := "Juan"
edad := 25

query := "SELECT * FROM usuarios WHERE nombre = @p1 AND edad > @p2"
results, err := c.SqlSelect(query, nombre, edad)
if err != nil {
    fmt.Printf("Error: %v\n", err)
}

// SELECT con LIKE
email := "gmail.com"
query = "SELECT * FROM usuarios WHERE email LIKE @p1"
results, err = c.SqlSelect(query, "%"+email+"%")

// Sin parámetros
results, err = c.SqlSelect("SELECT * FROM productos WHERE precio > 100")
```

---

### `SqlExec(query, args...)`

Ejecuta consultas SQL que no retornan datos (INSERT, UPDATE, DELETE). Retorna el número de filas afectadas.

**Parámetros:**
- `query` (string, requerido): Consulta SQL con placeholders (@p1, @p2, etc.)
- `args...` (variadic interface{}, opcional): Valores para reemplazar los placeholders

**Retorna:** `int64, error` (número de filas afectadas)

**Ejemplo:**
```go
// INSERT
rows, err := c.SqlExec("INSERT INTO usuarios (nombre, edad) VALUES (@p1, @p2)", "Maria", 30)
if err != nil {
    fmt.Printf("Error: %v\n", err)
}
fmt.Printf("Filas insertadas: %d\n", rows)

// UPDATE
rows, err = c.SqlExec("UPDATE usuarios SET edad = @p1 WHERE nombre = @p2", 26, "Juan")

// DELETE
rows, err = c.SqlExec("DELETE FROM usuarios WHERE edad < @p1", 18)

// Sin parámetros
rows, err = c.SqlExec("TRUNCATE TABLE logs")
```

**Características:**
- ✅ Protección contra SQL Injection usando parámetros
- ✅ Soporta cualquier tipo de consulta SQL
- ✅ Sintaxis flexible con placeholders @p1, @p2, @p3...

---

### `Close()`

Cierra la conexión con la base de datos.

**Parámetros:** Ninguno

**Retorna:** `error`

**Ejemplo:**
```go
defer c.Close()
```

---

## Ejemplo Completo

```go
package main

import (
	"fmt"
	mssql "github.com/PaoloESAN/mssql-easygo"
)

func main() {
	c := mssql.NewConexion("localhost", "sa", "password", "testdb")
	defer c.Close()

	// Crear base de datos
	err := c.CrearBD("testdb")
	if err != nil {
		fmt.Printf("Error creando BD: %v\n", err)
		return
	}

	// Cambiar a la base de datos creada
	err = c.CambiarBD("testdb")
	if err != nil {
		fmt.Printf("Error cambiando BD: %v\n", err)
		return
	}

	// Insertar datos
	err = c.Insert("usuarios", 1, "Juan", "juan@email.com")
	if err != nil {
		fmt.Printf("Error insertando: %v\n", err)
		return
	}

	err = c.Insert("usuarios", 2, "Maria", "maria@email.com")
	if err != nil {
		fmt.Printf("Error insertando: %v\n", err)
		return
	}

	// Consultar con SQL personalizado
	nombre := "Juan"
	query := "SELECT * FROM usuarios WHERE nombre = @p1"
	results, err := c.SqlSelect(query, nombre)
	if err != nil {
		fmt.Printf("Error en SQL: %v\n", err)
		return
	}

	fmt.Printf("\nUsuarios encontrados: %d\n", len(results))
	for _, row := range results {
		fmt.Println(row)
	}

	// Actualizar con SqlExec
	rows, err := c.SqlExec("UPDATE usuarios SET email = @p1 WHERE nombre = @p2", "juan_nuevo@email.com", "Juan")
	if err != nil {
		fmt.Printf("Error actualizando: %v\n", err)
		return
	}

	// Consultar todos los datos
	results, err = c.Select("usuarios", true)
	if err != nil {
		fmt.Printf("Error en SELECT: %v\n", err)
		return
	}

	fmt.Printf("\nTotal de registros: %d\n", len(results))

	for _, row := range results {
		fmt.Println(row)
	}
}
```

**Salida:**
```
Conexión a SQL Server exitosa - Base de datos: master
Base de datos 'testdb' creada exitosamente
Conexión a SQL Server exitosa - Base de datos: testdb
Datos insertados en 'usuarios' exitosamente
Datos insertados en 'usuarios' exitosamente
id | nombre | email
1 | Juan | juan@email.com
2 | Maria | maria@email.com

Total de registros: 2
map[email:juan@email.com id:1 nombre:Juan]
map[email:maria@email.com id:2 nombre:Maria]
```

---

## Manejo de Errores

Todos los métodos retornan `error`. Siempre verifica los errores:

```go
c := mssql.NewConexion("localhost", "sa", "password")

results, err := c.Select("usuarios")
if err != nil {
	fmt.Printf("Error: %v\n", err)
	return
}
```

---

## Notas Importantes

- Siempre cierra la conexión con `defer c.Close()` para evitar fugas de memoria
- Si no especificas base de datos, se conecta a **"master"** por defecto
- Si no especificas puerto, usa **1433** por defecto

---

## Requisitos

- Go 1.23 o superior
- Acceso a un servidor SQL Server
