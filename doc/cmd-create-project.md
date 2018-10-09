# Crear proyecto de aplicación

## Comando

El comando `create-project` creará un proyecto estándar con la estructura adecuada según el tipo de proyecto indicado.

### Parámetros

* **type**: Tipo de proyecto. Requerido.
* **system-id**: Identificador del sistema al cual pertenece la aplicación. Requerido.
* **application-id**: Identificador de la aplicación. Requerido.
* **tracker-id**: Identificador del issue tracker para el cual se crearán releases SQL eventuales.
                  Requerido para proyectos de tipo `oracle-sql-eventual`.

### Ejemplos

#### Crear proyecto SQL evolutivo

El siguiente comando:

```sh
buildr create-project --type oracle-sql-evolutional --system-id ve --application-id ve-sql-database
```

Crea un subirectorio llamado `ve-sql-database` en el directorio actual con la estructura del nuevo proyecto:

```tree
src/
  sql/
    inc/
      readme.md
    rep/
      readme.md
buildr.yaml
README.md
```

Teniendo en `buildr.yaml`:

```yaml
system-id: "ve"
application-id: "ve-sql-database"
type: "oracle-sql-evolutional"
```

Además se inicializará un repositorio git local en el directorio creado y se creará en el mismo el branch `master` con un commit inicial que contiene todos los archivos creados.

#### Crear proyecto SQL diferido

El siguiente comando:

```sh
buildr create-project --type oracle-sql-deferred --system-id ve --application-id ve-sql-process
```

Crea un subirectorio llamado `ve-sql-process` en el directorio actual con la estructura del nuevo proyecto:

```tree
src/
  sql/
    readme.md
buildr.yaml
README.md
```

Teniendo en `buildr.yaml`:

```yaml
system-id: "ve"
application-id: "ve-sql-process"
type: "oracle-sql-deferred"
```

Además se inicializará un repositorio git local en el directorio creado y se creará en el mismo el branch `master` con un commit inicial que contiene todos los archivos creados.

#### Crear proyecto SQL eventual

El siguiente comando:

```sh
buildr create-project --type oracle-sql-eventual --system-id ve --application-id ve-sql-eventual --tracker-id redmine-dieccs
```

Crea un subirectorio llamado `ve-sql-eventual` en el directorio actual con la estructura del nuevo proyecto:

```tree
src/
  sql/
    readme.md
buildr.yaml
README.md
```

Teniendo en `buildr.yaml`:

```yaml
system-id: "ve"
application-id: "ve-sql-eventual"
type: "oracle-sql-evolutional"
tracker-id: "redmine-dieccs"
```

Además se inicializará un repositorio git local en el directorio creado y se creará en el mismo el branch `base` con un commit inicial que contiene todos los archivos creados.
