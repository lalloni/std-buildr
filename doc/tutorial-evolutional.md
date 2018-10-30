# Paso a paso: SQL Evolutivo

Este documento esta destinado a guiar al desarrollador en el creación, desarrollo, empaquetado y publicación de un proyectos SQL Evolutivo según los estándares.

## Preparación del entorno de trabajo

El primer paso necesario es tener una copia local del proyecto de la Aplicación Oracle SQL Diferida en el entorno del desarrollador.

Para esto hay dos alternativas según la fase del proyecto en el que se desea trabajar, si el proyecto ya existe será suficiente con clonar el repositorio existente en GitLab al entorno del desarrollador, mientras que si el proyecto no existe habrá que crear uno nuevo en el entorno del desarrollador y subirlo a un nuevo repositorio en GitLab.

A continuación se detallan ambas alternativas:

### Clonar proyecto

Se desea **clonar** un proyecto SQL Diferido del sistema `factura-blockchain` que corresponde a la aplicacíón  `factura-blockchain-sql`.

Para ello se debe ejecutar el comando clone de git:

```bash
git clone git@gitlab.cloudint.afip.gob.ar:factura-blockchain/factura-blockchain-sql.git
cd factura-blockchain-sql
```

Luego de lo cual podremos empezar a realizar los cambios requeridos en el directorio actual.

### Crear proyecto

Se desea crear un proyecto SQL Evolutivo para el sistema `factura-blockchain` que corresponde a la aplicacion `factura-blockchain-sql`

Ejecutar el comando create-project del gestor de proyectos std-buildr:

```bash
std-buildr create-project --application-id factura-blockchain-sql  --system-id factura-blockchain --type oracle-sql-evolutional
```

Como resultado se obtiene la siguiente estructura de directorio:

```tree
factura-blockchain-sql
├── buildr.yaml
├── README.md
└── src
    └── sql
        └── inc
            └── README.md
        └── rep
            └── README.md
```

Y en `buildr.yaml`

```yaml
system-id: factura-blockchain
application-id: factura-blockchain-sql
type: oracle-sql-evolutional
```

**El directorio creado, `factura-blockchain-sql` es un proyecto git con un `commit` inicial que incluye a todos los archivos generados.  Se deberá agregar un repositorio `remote` y por último grabar la nueva estructura en el repositorio remoto de gitlab.**

```bash
git push -u origin master
```

## Desarrollar

Crear uno o más scripts en `src/sql`:

```tree
factura-blockchain-sql
├── buildr.yaml
├── README.md
└── src
    └── sql
        └── inc
            └── README.md
            └── 000001-ddl-crear-procedure.sql
        └── rep
            └── README.md
            └── procedure-pro.sql
```

Siendo el contenido de `src/sql/inc/000001-ddl-crear-procedure.sql` algo como:

```sql
@@../rep/procedure-pro.sql

EXIT
```

Agregar los scripts creados al index de git y crear un commit ejecutando los siguientes comandos:

```bash
git add src/sql/rep/procedure-pro.sql
git commit
```

## Publicar

Una vez creados los SQLs necesarios y realizadas las pruebas correspondientes, se desea publicar la versión 1.0.0 de este proyecto.

Crear un tag de la siguiente manera:

```bash
git tag v1.0.0 -a -m "Version 1.0.0"
```

Ejecutar el comando publish del gestor de proyectos std-buildr, el cual empaquetará y publicará el proyecto.

```bash
std-buildr publish
```

## Desarrollar y publicar versión con paquetes incrementales

Suponiendo que se realizaron sucesivos desarrollos y publicaciones de la aplicación, se desea desarrollar una versión subsiguiente (2.0.0) con nuevos scripts y publicar el paquete completo y dos apaquetes incrementales conteniendo las novedades desde la versión 1.0.0 y desde la versión 1.5.0.

Crear uno o más scripts en `src/sql`:

```tree
factura-blockchain-sql
└── src
    └── sql
        └── inc
            └── 000010-ddl-crear-tabla.sql
```

Agregar el script creado al index de git y crear un commit ejecutando los siguientes comandos:

```bash
git add src/sql/inc/000010-ddl-crear-tabla.sql
git commit
```

## Indicar a std-buildr la publicación de los incrementales

Agregar al archivo buildr.yaml la sección "from" conteniendo la lista de versiones base de los incrementales a publicar:

```yaml
from:
  - 1.0.0
  - 1.5.0
```

## Publicar

Crear un tag de la siguiente manera:

```bash
git tag v2.0.0 -a -m "Version 2.0.0"
```

Ejecutar el comando publish del gestor de proyectos std-buildr, el cual creará y publicará el paquete completo y los dos incrementales especificados.

```bash
std-buildr publish
```

Luego de la ejecución exitosa del comando quedarán publicados en Nexus los paquetes:

```tree
factura-blockchain-sql-2.0.0.zip
factura-blockchain-sql-2.0.0-from-1.0.0.zip
factura-blockchain-sql-2.0.0-from-1.5.0.zip
```

## Próximos pasos

Cada vez que se requiera publicar una nueva versión se deben repetir los pasos anteriores según el caso.
