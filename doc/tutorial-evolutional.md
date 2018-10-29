# Paso a paso: SQL Evolutivo

Este documento esta destinado a guiar al desarrollador en el creación, desarrollo, empaquetado y publicación de un proyectos SQL Evolutivo según los estándares.

## Crear proyecto

Se desea crear un proyecto SQL Evolutivo para el sistema `factura-blockchain` que corresponde a la aplicacion `factura-blockchain-sql`

Se ejecuta el comando create-project del gestor de proyectos std-buildr:

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

**El directorio creado, `factura-blockchain-sql` es un proyecto git con un `commit` inicial que incluye a todos los archivos generados.  Se deberá agregar un repositorio `remote` y por último grabar la nueva estructura en el repositorio remoto de gitlab:**

```bash
git push origin master
```

## Desarrollar

En caso de no ser el creador del proyecto y este exista en gitlab (si no existe es necesario [crear un proyecto](#Crear-proyecto) y continuar con la creacion del script) habrá que clonarlo desde gitlab con el siguiente comando:

```bash
git clone git@gitlab.cloudint.afip.gob.ar:factura-blockchain/factura-blockchain-sql.git
```

Siguiendo cno el ejemplo, crear un nuevo script en `src/sql`:

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

Siendo el contenido de `src/sql/inc/000001-ddl-crear-procedure.sql` el siguiente:

```sql
@@../rep/procedure-pro.sql

EXIT
```

Por último, se agregan los archivos al indice de git y se realiza un comit ejecutando los siguientes comandos :

```bash
git add src/sql/rep/procedure-pro.sql
git commit
```

## Publicar

Suponiendo que ya se agregaron los SQLs necesarios y se hicieron las pruebas correspondientes, ahora se desea publicar la versión 1.0.0 de este proyecto.

Para ello es necesario modificar el archivo `buildr.yaml` para agregar las configuraciones de empaquetamiento y eventualmente los datos del repositorio de nexus:

```yaml
system-id: factura-blockchain
application-id: factura-blockchain-sql
type: oracle-sql-evolutional
from:
    - 1.0.0
package:
    format: "zip"
nexus:
    url: "https://nexus.cloudint.afip.gob.ar/nexus/repository/"
```

Realizar un commit para asegurarse de que todos los cambios fueron impactados:

```bash
git add buildr.yaml
git commit -m "COMENTARIO DEL COMMIT"
```

Y luego crear un tag de la siguiente manera:

```bash
git tag v1.0.0 -a -m "Version 1.0.0"
```

Para finalizar se ejecuta el comando publish del gestor de proyectos std-buildr, el cual empaquetara (remplazando la instrucción `@@../rep/procedure-pro.sql` de `src/sql/inc/000001-ddl-crear-procedure.sql` con el contenido de `src/sql/rep/procedure-pro.sql`) y publicará el proyecto:

```bash
std-buildr publish
```

El comando package de `std-buildr` creara lo siguiente:

```tree
└── factura-blockchain-sql-1.0.0-from-1.0.0.zip
└── factura-blockchain-sql-1.0.0.zip
```

Cuyo contenido sera:

```tree
└── factura-blockchain-sql-000001-ddl-crear-procedure.sql
```

## Próximos pasos

Para publicar la version `3.0.0` se modifica la propiedad `from` de  `buildr.yaml` para empaquetar solo los cambios producidos desde las versiones configurada hasta las actual.

```yaml
system-id: factura-blockchain
application-id: factura-blockchain-sql
type: oracle-sql-evolutional
from:
    - 1.0.0
    - 2.0.0
    - 3.0.0
package:
    format: "zip"
nexus:
    url: "https://nexus.cloudint.afip.gob.ar/nexus/repository/"
```

El comando package de `std-buildr` creara lo siguiente:

```tree
└── factura-blockchain-sql-3.0.0-from-1.0.0.zip
└── factura-blockchain-sql-3.0.0-from-2.0.0.zip
└── factura-blockchain-sql-3.0.0-from-3.0.0.zip
└── factura-blockchain-sql-3.0.0.zip
```