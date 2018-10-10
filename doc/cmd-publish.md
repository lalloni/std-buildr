# Empaquetar versión de aplicación

## Comando

El comando `publish` creará un paquete estándar con la estructura adecuada según el tipo de proyecto indicado (segun como se especifica en el comando `package`) y lo publicara en nexus.

### Configuración

* **repository**: Ruta del repositorio en nexus donde se publicara el artefacto.

### Parámetros

* **force**: Si es especificado el paquete sera publicado incluso si contiene commits posteriores al ultimo tag o modificaciones que se encuentren en el directorio de trabajo.
* **trust**: Certificado de confianza en formato pem.  Opcional.

### Comportamiento

Este comando valida que el paquete a plublicar (creado al ejecutarse automaticamente el comando package) solo contenga al ultimo tag. 

Si se encuentran commits posteriores a dicho tag o modificaciones en el directorio de trabajo el paquete solo se publicara si fue especificado el parametro `force`.

El paquete se publicara en la siguiente ruta:

    {repository}/{system-id}/{application-id}/{version-id}/{package}

Siendo:

* **version-id**: Versión calculada al momento de crear el paquete.
* **package**: Nombre y extención del paquete a ser publicado.

### Ejemplos

#### Evolutivo

Estructura de fuentes:

```tree
src/
  sql/
    inc/
      000001-ddl.sql
      000002-dml.sql
      000003-dcl.sql
    rep/
      procedure-foo.sql
      view-bar.sql
      package-baz.sql
buildr.yaml
```

Teniendo en `buildr.yaml`:

```yaml
system-id: "factura-blockchain"
application-id: "factura-blockchain-sql"
type: "oracle-sql-evolutional"
package:
  format: "zip"
from:
  - "1.0.0"
  - "1.2.0"
repository: "https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/"
```

Se creó un tag denominado `v1.2.3` apuntando al último commit del branch actual.

Entonces, el comando:

```sh
buildr publish -trust ~/certs/homo.pem
```

Producirá los archivos:

```tree
target/
  factura-blockchain-sql-1.2.3.zip
  factura-blockchain-sql-1.2.3-from-1.0.0.zip
  factura-blockchain-sql-1.2.3-from-1.2.0.zip
```

Los cuales seran publicados en las siguientes URLs:

    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3.zip
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3-from-1.0.0.zip
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3-from-1.2.0.zip

#### Diferido

Estructura de fuentes:

```tree
src/
  sql/
    otra-tarea.sql
    una-tarea.sql
    y-una-tarea-mas.sql
buildr.yaml
```

Teniendo en `buildr.yaml`:

```yaml
system-id: "factura-blockchain"
application-id: "factura-blockchain-sql-process"
type: "oracle-sql-deferred"
package:
  format: "zip"
```

Todo la estructura de fuentes se encuentra versionada en git e incluida en commits del branch actual.

Se creó un tag denominado `v1.2.3` apuntando al último commit del branch actual.

Entonces, el comando:

```sh
buildr publish -trust ~/certs/homo.pem
```

Producirá el archivo:

```tree
target/
  factura-blockchain-sql-process-1.2.3.zip
```

El cual sera publicado en la siguiente URL:

    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql-process/1.2.3/factura-blockchain-sql-process-1.2.3.zip

#### Eventual

```tree
src/
  sql/
    001-ddl-create-tabla-temporal.sql
    002-dcl-grants-tabla-temporal.sql
    003-dml-extraccion-x.sql
    004-ddl-drop-tabla-temporal.sql
buildr.yaml
```

Teniendo en `buildr.yaml`:

```yaml
system-id: "factura-blockchain"
application-id: "factura-blockchain-sql-eventual"
type: "oracle-sql-eventual"
tracker-id: "redmine-dieccs"
issue-id: "1234"
package:
  format: "zip"
```

Todo la estructura de fuentes se encuentra versionada en git e incluida en commits del branch actual.

Se creó un tag denominado `redmine-dieccs-1234-1` apuntando al último commit del branch actual.

Entonces, el comando:

```sh
buildr publish -trust ~/certs/homo.pem
```

Producirá el archivo:

```tree
target/
  factura-blockchain-sql-eventual-redmine-dieccs-1234-1.zip
```

El cual sera publicado en la siguiente URL:

    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql-eventual/redmine-dieccs-1234-1/factura-blockchain-sql-eventual-redmine-dieccs-1234-1.zip