# Empaquetar versión de aplicación

## Comando

El comando `package` creará un paquete estándar con la estructura adecuada según el tipo de proyecto indicado.

### Identificador de versión

Este comando determina automáticamente el identificador de la versión del paquete producido en función del último tag que encuentre en el branch actual del proyecto.

**Para poder ejecutarlo se debe haber creado un tag previamente.**

El identificador de la versión calculado tendrá en cuenta tanto el nombre del último tag encontrado, los commits posteriores al mismo y las posibles modificaciones que se encuentren en el directorio de trabajo, adicionando al identificador de versión generado indicadores que reflejen estas cuestiones.

El identificador generado tendrá la forma:

    {base-version}[-{count}-{commit}][-dirty]

Siendo:

* `{base-version}`: Identificador de la versión, calculado a partir del nombre del último tag encontrado. Se calcula y se valida dependiendo del tipo de proyecto. Por ejemplo: "1.0.0" o "redmine-dieccs-1234".
* `{count}`: Cantidad de commits encontrados en el branch actual desde el último tag. Sólo se agrega si se encuentran commits posteriores al último tag.
* `{commit}`: Identificador del último commit encontrado en el branch actual. Sólo se agrega si se encuentran commits posteriores al último tag.
* `-dirty`: Indicador de que el directorio de trabajo contiene modificaciones no capturadas en git. Sólo se agrega si se encontraran dichas modificaciones.

### Configuración

* **package.format**: Formato del archivo empaquetado. Debe ser `tar.xz`, `tar.gz` o `zip`.
                      Opcional. De no estar especificado se utiliza el valor `zip`.
* **from**: Lista de versiones preexistentes para las que se crearán paquetes incrementales. Solo es utilizado en los proyectos de tipo `oracle-sql-evolutional`. Opcional. De no estar especificado, no se crearán paquetes incrementales.

### Parámetros

* **allow-dirty**: Permite construir paquetes que contengan archivos modificados en el directorio de trabajo
* **allow-untagged**: Permite construir paquetes que contengan commits posteriores al último tag

### Comportamiento básico

Este comando realiza los siguientes pasos:

1. Limpieza del directorio de construcción `target`
2. Cálculo de la versión según el nombre del último tag del branch
3. Validaciones de estándares según el tipo de proyecto:
    * Nombre del tag
    * Nombre del branch
    * Estado de código fuente
    * Historia de commits
4. Empaquetamiento

Para el tercer paso se verifica que el código fuente desde el cual se contruyeron los paquetes no incluya modificaciones no capturadas en el último tag del branch. Esta verificación incluye validar la ausencia de:

1. Commits posteriores al tag (se omite si se especifica `--allow-untagged`)
2. Archivos nuevos o modificados en el directorio de trabajo (se omite si se especifica `--allow-dirty`)

En este paso también se verifica que los nombres del tag y del branch cumplan con los requerimientos estándar específicos según el tipo de proyecto.

Si no se cumple cualquier verificación de las anteriores, se cancela la construcción de los paquetes, se informa del problema al usuario y el programa termina con estado de error.

El comportamiento del cuarto paso depende del tipo de proyecto y se describe a continuación:

### Comportamiento según tipo de proyecto

#### SQL Evolutivo

En este tipo de proyectos se crean dos tipos de paquetes:

* El **paquete completo** que contiene todos los scripts de evolución
* Los **paquetes incrementales** que contienen un subconjunto de los scripts para actualizar desde una versión específica hasta la versión empaquetada

Todos los scripts SQL del directorio `src/sql/inc` son incluidos en el paquete completo.

Los scripts SQL del directorio `src/sql/inc` serán incluidos en los paquetes incrementales si fueron introducidos entre una de las versiones listadas en la configuración `from` y la versión que se está empaquetando.

Todos los scripts SQL del directorio `src/sql/inc` son procesados reemplazando las directivas `@@` por el contenido de los archivos que las mismas referencien.

Se valida que el tag cumpla con [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html) y que tenga el caracter 'v' prefijado.

Se valida que los nombres de los scripts cumplan con el estándar de nombres, abortando el proceso en caso de encontrar nombres incorrectos.

Todos los scripts SQL incluidos son renombrados anteponiendo el identificador de la aplicación.
Si los nombres de los scripts fuente poseen el identificador de la aplicación como prefijo, esta acción no se realiza.

#### SQL Diferido

En este tipo de proyecto se crea un paquete que contiene todos los scripts SQL del directorio `src/sql`.

Se valida que el tag cumpla con [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html) y que tenga el caracter 'v' prefijado.

Se valida que los nombres de los scripts cumplan con el estándar de nombres, abortando el proceso en caso de encontrar nombres incorrectos.

Todos los scripts SQL incluidos son renombrados anteponiendo el identificador de la aplicación.
Si los nombres de los scripts fuente ya poseen el identificador de la aplicación como prefijo, esta acción no se realiza.

#### SQL Eventual

En este tipo de proyecto se crea un paquete que contiene todos los scripts SQL del directorio `src/sql`.

Se valida que el nombre del branch comience con el patrón:

    {tracker-id}-{issue-id}{custom}

Siendo:

* `{tracker-id}`: El identificador de la instancia del issue tracker
* `{issue-id}`: El identificador del issue que da origen al eventual
* `{custom}`: Una cadena opcional que se puede utilizar para organizar múltiples branches del mismo eventual o para agregar una descripción corta del motivo del eventual. Si se especifica **debe** comenzar con el caracter '-' o el caracter '/'.

Se valida que el tag cumpla con el patrón:

    {tracker-id}-{issue-id}-{version}

Siendo:

* `{tracker-id}`: El identificador de la instancia del issue tracker.
* `{issue-id}`: El identificador del issue que da origen al eventual.
* `{version}`: El identificador de la versión del eventual. Debe ser un número entero basado en 1 que se incrementa para cada versión.

Se valida que los nombres de los scripts cumplan con el estándar de nombres, abortando el proceso en caso de encontrar nombres incorrectos.

Todos los scripts SQL incluidos son renombrados anteponiendo el identificador del issue tracker, del issue y la versión del eventual.
Si los nombres de los scripts fuente ya poseen el identificador del issue tracker, del issue y la versión del eventual como prefijo, esta acción no se realiza.

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
from:
  - "1.0.0"
  - "1.2.0"
```

Y en `src/sql/incremental/000001-ddl.sql`:

```sql
@@../replaceable/procedure-foo.sql
```

Todo la estructura de fuentes se encuentra versionada en git e incluida en commits del branch actual.

Se creó un tag denominado `v1.2.3` apuntando al último commit del branch actual.

Entonces, el comando:

```sh
buildr package
```

Producirá los archivos:

```tree
target/
  factura-blockchain-sql-1.2.3.zip
  factura-blockchain-sql-1.2.3-from-1.0.0.zip
  factura-blockchain-sql-1.2.3-from-1.2.0.zip
```

La estructura del archivo `factura-blockchain-sql-1.2.3.zip` será:

```tree
factura-blockchain-sql-000001-dml.sql
factura-blockchain-sql-000002-ddl.sql
factura-blockchain-sql-000003-dcl.sql
```

La estructura del archivo `factura-blockchain-sql-1.2.3-from-1.2.0.zip` será:

```tree
factura-blockchain-sql-000003-dcl.sql
```

Todos los script incrementales tendrán reemplazadas las directivas `@@{file}` por el contenido del archivo `{file}`.

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
```

Todo la estructura de fuentes se encuentra versionada en git e incluida en commits del branch actual.

Se creó un tag denominado `v1.2.3` apuntando al último commit del branch actual.

Entonces, el comando:

```sh
buildr package
```

Producirá el archivo:

```tree
target/
  factura-blockchain-sql-process-1.2.3.zip
```

Y la estructura de dicho archivo `factura-blockchain-sql-process-1.2.3.zip` será:

```tree
factura-blockchain-sql-process-otra-tarea.sql
factura-blockchain-sql-process-una-tarea.sql
factura-blockchain-sql-process-y-una-tarea-mas.sql
```

#### Eventual

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

Todo la estructura de fuentes se encuentra versionada en git e incluida en commits del branch actual.

Se creó un tag denominado `redmine-dieccs-1234-1`, ya que es la versión "1" del eventual originado en la petición "1234" de la instancia de Redmine llamada "redmine-dieccs", apuntando al último commit del branch actual.

Entonces, el comando:

```sh
buildr package
```

Producirá el archivo:

```tree
target/
  factura-blockchain-sql-eventual-redmine-dieccs-1234-1.zip
```

Y la estructura de dicho archivo `factura-blockchain-sql-eventual-redmine-dieccs-1234-1.zip` será:

```tree
redmine-dieccs-1234-1-001-ddl-create-tabla-temporal.sql
redmine-dieccs-1234-1-002-dcl-grants-tabla-temporal.sql
redmine-dieccs-1234-1-003-dml-extraccion-x.sql
redmine-dieccs-1234-1-004-ddl-drop-tabla-temporal.sql
```
