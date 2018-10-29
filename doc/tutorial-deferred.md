# Paso a paso: SQL Diferido

Este documento esta destinado a guiar al desarrollador en la creación, desarrollo, empaquetado y publicación de un proyectos SQL Diferidos según los estándares.

## Crear proyecto

Se desea crear un proyecto SQL Diferido para el sistema `factura-blockchain` que corresponde a la aplicacion `factura-blockchain-sql-process`

Se ejecuta el comando create-project del gestor de proyectos std-buildr:

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
git push origin master
```

## Desarrollar

En caso de no ser el creador del proyecto y este exista en gitlab (si no existe es necesario [crear un proyecto](#Crear-proyecto) y continuar con la creacion del script) habrá que clonarlo desde gitlab con el siguiente comando:

```bash
git clone git@gitlab.cloudint.afip.gob.ar:factura-blockchain/factura-blockchain-sql-process.git
```

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
└── factura-blockchain-sql-process-1.0.0.zip
```

Cuyo contenido sera:

```tree
└── factura-blockchain-sql-process-proceso-diario.sql
```
