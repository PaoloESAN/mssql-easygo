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
conn := mssql.NewConexion("localhost", "sa", "password", 1433, "midb")
```

---

## Funciones

### `NewConexion(server, user, password, port(opcional), database(opcional))`

Crea una nueva conexión a SQL Server. La conexión se realiza automáticamente.

**Parámetros:**
- `server` (string, requerido): Servidor SQL Server (ej: "localhost")
- `user` (string, requerido): Usuario (ej: "sa")
- `password` (string, requerido): Contraseña
- `port` (int, opcional): Puerto (por defecto 1433)
- `database` (string, opcional): Nombre de la base de datos (por defecto "master")

**Retorna:** `*Conexion`

**Ejemplo:**
```go
c := mssql.NewConexion("localhost", "sa", "password")
c := mssql.NewConexion("localhost", "sa", "password", 1433)
c := mssql.NewConexion("localhost", "sa", "password", 1433, "midb")
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
	"github.com/PaoloESAN/mssql-easygo"
)

func main() {
	c := mssql.NewConexion("localhost", "sa", "password")
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

	// Consultar datos
	results, err := c.Select("usuarios", true)
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
