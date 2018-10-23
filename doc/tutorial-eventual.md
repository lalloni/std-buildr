# Paso a paso: SQL Eventual

Este documento esta destinado a guiar al desarrollador en el creación, desarrollo y empaquetado y publicación de un proyectos SQL Eventual según los estándares.

## Crear proyecto

Se desea crear un proyecto SQL Eventual para el sistema `factura-blockchain` que corresponde a la aplicacion `factura-blockchain-sql-eventual` y la instancia de redmine en la que se originaran las peticiones para la creación de eventuales es `redmine-dieccs`

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

**Por ultimo se debe ejecutar el siguiente comando para grabar la estructura del branch base, del cual surgiran los nuevos eventuales, en el repositorio remoto de gitlab:**

```bash
git push
```

## Crear eventual

En caso de no ser el creador del proyecto [Crear proyecto](#Crear-proyecto) se debe clonar el proyecto de gitlab con el siguiente comando:

```bash
git clone git@gitlab.cloudint.afip.gob.ar:factura-blockchain/factura-blockchain-sql-eventual.git
```

Se requiere iniciar el desarrollo de un eventual originado por la petición `1234`. Para el cual es necesario dos scripts DML, un DCL y un DDL que deben ser ejecutados en el siguiente orden:

1. obtener-data
2. crear-tabla
3. agregar-permisos
4. agregar-data

Se ejecuta el comando `create-eventual` del `std-buildr`

```bash
std-buildr create-eventual --dml obtener --ddl crear-tabla --dcl agregar-permisos --dml agregar-data -i 1234
```

En caso de no existir el branch `base` inicializara el proyecto para obtener la estructura base y realizara un commit con ella.

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

Y en `buildr.yaml`

```yaml
system-id: factura-blockchain
application-id: factura-blockchain-sql-eventual
type: oracle-sql-eventual
tracker-id: redmine-dieccs
issue-id: 1234
```

En el repositorio local se crea un branch con el siguiente nombre: `redmine-dieccs-1234`

## Desarrollar eventual

Se modifica el contenido de los scripts creados por el comando `create-eventual`.

Se agregan los archivos al indice de git y se realiza un comit ejecutando los siguientes comandos :

```bash
git add src/sql/001-dml-obtener-data.sql src/sql/002-ddl-crear-tabla.sql src/sql/003-dcl-agregar-permisos.sql src/sql/004-dml-agregar-data.sql
git commit
```

## Publicar eventual

Se desea publicar la versión `1`  del eventual `redmine-dieccs-1234`. Para ello es necesario crear un tag de la siguiente manera:

```bash
git tag redmine-dieccs-1234-1 -sa -m "Version 1 del eventual redmine-dieccs-1234"
```

Se modifica el `buildr.yaml` para configurar de empaquetamiento y repositorio de nexus:

```yaml
system-id: factura-blockchain
application-id: factura-blockchain-sql-eventual
type: oracle-sql-eventual
tracker-id: redmine-dieccs
issue-id: 1234
package:
    format: "zip"
nexus:
    url: "https://nexus.cloudint.afip.gob.ar/nexus/repository/"
```

Para finalizar se ejecuta el comando publish del gestor de proyectos std-buildr, el cual empaquetara y publicara el proyecto:

```bash
std-buildr publish
```

El comando package de `std-buildr` creara lo siguiente:

```tree
└── factura-blockchain-sql-eventual-redmine-dieccs-1234-1.zip
```

Cuyo contenido sera:

```tree
└── redmine-dieccs-1234-1-001-dml-obtener-data.sql
└── redmine-dieccs-1234-1-002-ddl-crear-tabla.sql
└── redmine-dieccs-1234-1-003-dcl-agregar-permisos.sql
└── redmine-dieccs-1234-1-004-dml-agregar-data.sql
```


## Próximos pasos

Para la nueva version del eventual `redmine-dieccs-1234` se deben seguir los pasos desde [Desarrollar eventual](#Desarrollar-eventual)

Para crear nuevos eventuales se debe seguir los pasos desde [Crear eventual](#Crear-eventual)