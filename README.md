# buildr

## Tipos de construcción

- oracle-sql-evolutional
- oracle-sql-deferred
- oracle-sql-eventual
- microsoft-net-web
- microsoft-net-web-core
- microsoft-net-lib
- go-web
- go-command

## UC 1: Empaquetar app Oracle SQL

Proyecto:

```tree
src/
  sql/
    incremental/
      000001.sql
      000002.sql
      000003.sql
    replaceable/
      procedure-foo.sql
      view-bar.sql
      package-baz.sql
buildr.yaml
```

Teniendo en `buildr.yaml`:

```yaml
system-id: "factura-blockchain"
application-id: "sql"
type: "oracle-sql-evolutional"
from:
  - "1.0.0"
  - "1.2.0"
```

En `src/sql/incremental/000001.sql`:

```sql
@@../replaceable/procedure-foo.sql
```

Y en el proyecto:

El comando:

```sh
buildr package
```

Produce:

```tree
target/
  factura-blockchain-sql-1.2.3.zip
  factura-blockchain-sql-1.2.3-from-1.0.0.zip
  factura-blockchain-sql-1.2.3-from-1.2.0.zip
```

Conteniendo `factura-blockchain-sql-1.2.3.zip`:

```tree
000001.sql
000002.sql
000003.sql
```

Conteniendo `factura-blockchain-sql-1.2.3-from-1.2.0.zip`:

```tree
000003.sql
```

Todos los script incrementales tendrán reemplazadas las directivas `@@<file>` por el contenido del archivo `<file>`.

## Notas

### Determinación de identificador de versión desde una WC

buildr debe determinar la versión sobre la que está trabajando, para eso se basa en información disponible en git:

- Determina si la WC está limpia (DIRTY=false) o tiene cambios (DIRTY=true)
- Busca el último tag del branch actual (TAG)
- Busca el último commit del branch actual (COMMIT)
- VERSION=$(git describe --abbrev=40 HEAD)[1:]
- Y si DIRTY => VERSION = VERSION + "-dirty"
