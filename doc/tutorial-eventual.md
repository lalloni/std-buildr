# Paso a paso: SQL Eventual

Este documento esta destinado a guiar al desarrollador en la creación, desarrollo, empaquetado y publicación de un proyectos SQL Eventual según los estándares.

## Preparación del entorno de trabajo

El primer paso necesario es tener una copia local del proyecto de la Aplicación Oracle SQL Eventual en el entorno del desarrollador.

Para esto hay dos alternativas según la fase del proyecto en el que se desea trabajar, si el proyecto ya existe será suficiente con clonar el repositorio existente en GitLab al entorno del desarrollador, mientras que si el proyecto no existe habrá que crear uno nuevo en el entorno del desarrollador y subirlo a un nuevo repositorio en GitLab.

A continuación se detallan ambas alternativas:

### Clonar proyecto

Se desea **clonar** un proyecto SQL Eventual del sistema `factura-blockchain` que corresponde a la aplicacíón `factura-blockchain-sql-eventual`.

Para ello se ejecuta el comando clone de git:

```bash
git clone git@gitlab.cloudint.afip.gob.ar:factura-blockchain/factura-blockchain-sql-eventual.git
cd factura-blockchain-sql-eventual
```

Luego de lo cual podremos empezar a realizar los cambios requeridos en el directorio actual.

### Crear proyecto

Se desea **crear** un proyecto SQL Eventual para el sistema `factura-blockchain` que corresponde a la aplicacion `factura-blockchain-sql-eventual` y la instancia de redmine en la que se originaran las peticiones para la creación de eventuales es `redmine-dieccs`

Se ejecuta el comando create-project del gestor de proyectos std-buildr:

```bash
std-buildr create-project --application-id factura-blockchain-sql-eventual  --system-id factura-blockchain --type oracle-sql-eventual --tracker-id redmine-dieccs
```

Como resultado se obtiene la siguiente estructura de directorio:

```tree
factura-blockchain-sql-eventual
├── buildr.yaml
├── README.md
└── src
    └── sql
        └── README.md
```

Y en `buildr.yaml`

```yaml
system-id: factura-blockchain
application-id: factura-blockchain-sql-eventual
type: oracle-sql-eventual
tracker-id: redmine-dieccs
```

**El directorio creado `factura-blockchain-sql-eventual` es un proyecto git con un `commit` inicial que incluye a todos los archivos generados.  Se deberá agregar un repositorio `remote` y por último grabar la nueva estructura en el repositorio remoto de gitlab.**

```bash
git push -u origin master
```

## Crear eventual

Se requiere iniciar el desarrollo de un eventual originado por la petición `1234`. Para el cual es necesario dos scripts DML, un DCL y un DDL que deben ser ejecutados en el siguiente orden:

1. obtener-data
2. crear-tabla
3. agregar-permisos
4. agregar-data

Ejecutar el comando `create-eventual` de `std-buildr`

```bash
std-buildr create-eventual --dml obtener-data --ddl crear-tabla --dcl agregar-permisos --dml agregar-data -i 1234
```

Como resultado se obtiene la siguiente estructura de directorio:

```tree
factura-blockchain-sql-eventual
├── buildr.yaml
├── README.md
└── src
    └── sql
        └── README.mdd
        └── 001-dml-obtener-data.sql
        └── 002-ddl-crear-tabla.sql
        └── 003-dcl-agregar-permisos.sql
        └── 004-dml-agregar-data.sql
```

Y en `buildr.yaml` quedará este contenido:

```yaml
system-id: factura-blockchain
application-id: factura-blockchain-sql-eventual
type: oracle-sql-eventual
tracker-id: redmine-dieccs
issue-id: 1234
```

Además, en el repositorio local se crea un branch llamado `redmine-dieccs-1234` en el cual se realizará el desarrollo.

## Desarrollar eventual

Modificar el contenido de los scripts creados por el comando `create-eventual` u otros scripts que se desee agregar.

Agregar los archivos al index de git y crear un commit ejecutando los siguientes comandos:

```bash
git add src/sql/001-dml-obtener-data.sql src/sql/002-ddl-crear-tabla.sql src/sql/003-dcl-agregar-permisos.sql src/sql/004-dml-agregar-data.sql
git commit
```

## Publicar eventual

Se desea publicar la versión `1`  del eventual `redmine-dieccs-1234`.

Crear un tag de la siguiente manera:

```bash
git tag redmine-dieccs-1234-1 -a -m "Version 1 del eventual redmine-dieccs-1234"
```

Ejecutar el comando publish del gestor de proyectos std-buildr, el cual empaquetará y publicará el proyecto:

```bash
std-buildr publish
```

A partir de este momento ya no se podrá modificar la versión publicada del eventual y si se requiriera hacer algún cambio será necesario publicar una nueva versioń del mismo como se describe en la siguiente sección.

## Próximos pasos

Para publicar una nueva version del eventual `redmine-dieccs-1234` se deben seguir los pasos desde [Desarrollar eventual](#Desarrollar-eventual).

Para crear nuevos eventuales se debe seguir los pasos desde [Crear eventual](#Crear-eventual).
