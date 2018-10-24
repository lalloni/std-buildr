# Crear nuevo eventual

## Comando

El comando `create-eventual` prepara el proyecto para comenzar a trabajar en un nuevo SQL eventual.

Sólo se puede ejecutar dentro de un proyecto de tipo `oracle-sql-eventual`.

### Parámetros

* **issue-id**: Identificador de la petición que origina el eventual. Requerido.
* **ddl**: Nombre de script DDL a crear. Opcional. Repetible.
* **dml**: Nombre de script DML a crear. Opcional. Repetible.
* **dcl**: Nombre de script DCL a crear. Opcional. Repetible.

### Ejemplos

Crear nuevo SQL eventual originado por el issue 1234 e inicializa 3 archivos para desarrollar scripts DDL, DCL y DML ordenados según línea de comando:

```sh
buildr create-eventual --issue-id 1234 --ddl tabla-temporal --dcl grants-tabla-temporal --dml consulta-x
```

Crea un branch llamado "redmine-dieccs-1234" que contiene un commit con la estructura básica del nuevo eventual y hace un checkout de dicho branch.

Si en el repositorio local se encontrara un branch local llamado `base` se utilizará como padre del nuevo branch. Si no se encontrara dicho branch pero se encontrara un branch remoto llamado `origin/base` se creará el branch local `base` apuntando al mismo commit que el remoto y se utilizará como padre del nuevo branch. Este branch `base` puede ser utilizado para contener archivos que se desean reutilizar en todos los eventuales generados a partir de su existencia.

Crea archivos `.sql` para cada script especificado con `--dml`, `--ddl` y `--dcl` en el orden especificado en el comando. Para el ejemplo anterior, sería:

```tree
src/
  sql/
    001-ddl-tabla-temporal.sql
    002-dcl-grants-tabla-temporal.sql
    003-dml-consulta-x.sql
```

Esto deja el proyecto listo para que el desarrollador modifique dichos scripts sql o agregue (manualmente) más scripts sql.
