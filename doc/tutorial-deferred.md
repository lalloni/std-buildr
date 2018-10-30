# Paso a paso: SQL Diferido

Este documento esta destinado a guiar al desarrollador en la preparación del entorno de trabajo, el desarrollo, el empaquetado y la publicación de un proyectos SQL Diferidos según los estándares.

## Preparación del entorno de trabajo

El primer paso necesario es tener una copia local del proyecto de la Aplicación Oracle SQL Diferida en el entorno del desarrollador.

Para esto hay dos alternativas según la fase del proyecto en el que se desea trabajar, si el proyecto ya existe será suficiente con clonar el repositorio existente en GitLab al entorno del desarrollador, mientras que si el proyecto no existe habrá que crear uno nuevo en el entorno del desarrollador y subirlo a un nuevo repositorio en GitLab.

A continuación se detallan ambas alternativas:

### Clonar proyecto

Se desea **clonar** un proyecto SQL Diferido del sistema `factura-blockchain` que corresponde a la aplicacíón  `factura-blockchain-sql-process`.

Para ello se ejecuta el comando clone de git:

```bash
git clone git@gitlab.cloudint.afip.gob.ar:factura-blockchain/factura-blockchain-sql-process.git
cd factura-blockchain-sql-process
```

Luego de lo cual podremos empezar a realizar los cambios requeridos en el directorio actual.

### Crear proyecto

Se desea **crear** un proyecto SQL Diferido para el sistema `factura-blockchain` que corresponde a la aplicacion `factura-blockchain-sql-process`

Para ello se ejecuta el comando create-project del gestor de proyectos std-buildr:

```bash
std-buildr create-project --application-id factura-blockchain-sql-process  --system-id factura-blockchain --type oracle-sql-deferred
```

Como resultado se obtiene la siguiente estructura de directorio:

```tree
factura-blockchain-sql-process
├── buildr.yaml
├── README.md
└── src
    └── sql
        └── README.md
```

Y en `buildr.yaml`

```yaml
system-id: factura-blockchain
application-id: factura-blockchain-sql-process
type: oracle-sql-deferred
```

**El directorio creado, `factura-blockchain-sql-process` es un proyecto git con un `commit` inicial que incluye a todos los archivos generados.  Se deberá agregar un repositorio `remote` y por último grabar la nueva estructura en el repositorio remoto de gitlab:**

```bash
git push -u origin master
```

## Desarrollar

Crear un nuevo script en `src/sql`:

```tree
factura-blockchain-sql-process
├── buildr.yaml
├── README.md
└── src
    └── sql
        └── README.md
        └── proceso-diario.sql
```

Se agregan los archivos al indice de git y se realiza un comit ejecutando los siguientes comandos :

```bash
git add src/sql/proceso-diario.sql
git commit
```

## Publicar

Se desea publicar la versión 1.0.0 de este proyecto.

Crear un tag de la siguiente manera:

```bash
git tag v1.0.0 -a -m "Version 1.0.0"
```

Para finalizar se ejecuta el comando publish del gestor de proyectos std-buildr, el cual empaquetara y publicara el proyecto:

```bash
std-buildr publish
```

El comando package de `std-buildr` creara lo siguiente:

```tree
factura-blockchain-sql-process
└── target/
    └── factura-blockchain-sql-process-1.0.0.zip
```

Cuyo contenido sera:

```tree
factura-blockchain-sql-process-proceso-diario.sql
```

Y publicará el archivo zip al repositorio Nexus correspondiente.
