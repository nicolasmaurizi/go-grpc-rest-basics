# gRPC + REST + Go + PostgreSQL (BloomRPC Demo)

Proyecto de ejemplo escrito en **Go** que muestra cómo:

- Exponer un **servicio gRPC** para crear usuarios (ideal para probar con **BloomRPC**).
- Guardar los datos en una base **PostgreSQL**.
- Exponer un endpoint **REST `GET /users`** para listar todos los usuarios.
- Separar la configuración de la base de datos en un módulo externo (`config`) y la conexión en otro (`db`).

Es un proyecto perfecto para:
- Practicar **gRPC en Go**.
- Entender cómo combinar **gRPC + REST** en el mismo binario.
- Usar **PostgreSQL** como capa de persistencia.

---

## 🧱 Arquitectura general

### Flujo de funcionamiento

- **gRPC**:
  - Método `CreateUser(name, email)`.
  - Inserta el usuario en la tabla `users` de PostgreSQL.
- **REST**:
  - Endpoint `GET /users`.
  - Consulta todos los usuarios desde PostgreSQL.
  - Devuelve la lista en formato JSON.

### Diagrama (Mermaid)

> Podés ver este diagrama directamente en GitHub si Mermaid está habilitado o copiarlo en https://mermaid.live

```mermaid
flowchart LR
    A[BloomRPC / Cliente gRPC] -->|CreateUser(name,email)| B[gRPC Server en Go]
    B -->|INSERT INTO| C[(PostgreSQL<br/>Tabla users)]

    D[Cliente REST<br/>(Browser / Postman)] -->|GET /users| E[HTTP Server en Go]
    E -->|SELECT * FROM users| C
    C -->|Lista de usuarios| E
    E -->|JSON<br/>200 OK| D
