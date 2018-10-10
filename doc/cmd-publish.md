# Publicar versión

## Comando

El comando `publish` construye y publica una versión de una aplicación.

### Configuración

* **nexus.url**: URL base del servicio repositorio de artefactos. Opcional. Por defecto apunta a `https://nexus.cloudint.afip.gob.ar/nexus/repository`.

### Parámetros

* **allow-dirty**: Permite construir paquetes que contengan archivos modificados en el directorio de trabajo.
* **allow-untagged**: Permite construir paquetes que contengan commits posteriores al último tag.
* **trust**: Ubicación del archivo en formato PEM que contiene la cadenas de certificados de confianza. Opcional. Por defecto se utilizan los certificados confiables de AFIP de manera tal que se pueda publicar al Nexus de la organización.

### Comportamiento

Este comando realiza los siguientes pasos:

1. Limpieza del directorio de construcción `target`
2. Ejecución del [comando package](cmd-package.md)
3. Publicación de los paquetes construidos en el paso anterior

En el tercer paso el programa publica los artefactos construidos al repositorio correspondiente al sistema y al tipo de aplicación, incluyendo los archivos de verificación correspondientes (digestos MD5 y SHA1).

Para la publicación realizada en el paso 2 se calcula la ubicación de los paquetes en Nexus mediante el patrón:

    /{system-id}-{repo-type}/{system-id}/{application-id}/{version-id}/

Siendo:

* `{system-id}`: Identificador del sistema especificado en configuración
* `{repo-type}`: Tipo de repositorio determinado según el tipo de aplicación configurado
* `{application-id}`: Identificador de la aplicación especificado en configuración
* `{version-id}`: Versión calculada al momento de crear el paquete

### Ejemplos

#### SQL Evolutivo

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
from:
  - "1.0.0"
  - "1.2.0"
```

Se creó un tag denominado `v1.2.3` apuntando al último commit del branch actual.

Entonces, el comando:

```sh
buildr publish
```

Producirá los archivos:

```tree
target/
  factura-blockchain-sql-1.2.3.zip
  factura-blockchain-sql-1.2.3-from-1.0.0.zip
  factura-blockchain-sql-1.2.3-from-1.2.0.zip
```

Estos archivos y sus digestos se publicarán en Nexus en las siguientes URLs:

    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3.zip
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3.zip.md5
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3.zip.sha1
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3-from-1.0.0.zip
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3-from-1.0.0.zip.md5
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3-from-1.0.0.zip.sha1
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3-from-1.2.0.zip
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3-from-1.2.0.zip.md5
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql/1.2.3/factura-blockchain-sql-1.2.3-from-1.2.0.zip.sha1

#### SQL Diferido

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
```

Todo la estructura de fuentes se encuentra versionada en git e incluida en commits del branch actual.

Se creó un tag denominado `v1.2.3` apuntando al último commit del branch actual.

Entonces, el comando:

```sh
buildr publish
```

Producirá el archivo:

```tree
target/
  factura-blockchain-sql-process-1.2.3.zip
```

Este archivo y sus digestos se publicarán en Nexus en las siguientes URLs:

    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql-process/1.2.3/factura-blockchain-sql-process-1.2.3.zip
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql-process/1.2.3/factura-blockchain-sql-process-1.2.3.zip.md5
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql-process/1.2.3/factura-blockchain-sql-process-1.2.3.zip.sha1

#### SQL Eventual

Estructura de fuentes:

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
```

Toda la estructura de fuentes se encuentra versionada en git e incluida en commits del branch actual, denominado `redmine-dieccs-1234`.

Se creó un tag denominado `redmine-dieccs-1234-1` apuntando al último commit del branch actual.

Entonces, el comando:

```sh
buildr publish
```

Producirá el archivo:

```tree
target/
  factura-blockchain-sql-eventual-redmine-dieccs-1234-1.zip
```

Este archivo y sus digestos se publicarán en Nexus en las siguientes URLs:

    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql-eventual/redmine-dieccs-1234-1/factura-blockchain-sql-eventual-redmine-dieccs-1234-1.zip
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql-eventual/redmine-dieccs-1234-1/factura-blockchain-sql-eventual-redmine-dieccs-1234-1.zip.md5
    https://nexus.cloudhomo.afip.gob.ar/nexus/repository/factura-blockchain-raw/factura-blockchain/factura-blockchain-sql-eventual/redmine-dieccs-1234-1/factura-blockchain-sql-eventual-redmine-dieccs-1234-1.zip.sha1
