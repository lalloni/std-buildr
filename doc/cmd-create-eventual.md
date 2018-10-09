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

Crea un branch llamado "redmine-dieccs-1234" basado en el branch "base" (si éste no existe, lo crea conteniendo un commit con el archivo de configuración) y hace un checkout de dicho branch.

Crea archivos `.sql` para cada script especificado con `--dml`, `--ddl` y `--dcl` en el orden especificado en el comando. Para el ejemplo anterior, sería:

```tree
src/
  sql/
    001-ddl-tabla-temporal.sql
    002-dcl-grants-tabla-temporal.sql
    003-dml-consulta-x.sql
```

Esto deja el proyecto listo para que el desarrollador modifique dichos scripts sql o agregue (manualmente) más scripts sql.
