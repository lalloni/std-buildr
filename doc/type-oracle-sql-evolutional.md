# Proyecto Oracle SQL Evolutivo

## Estructura de código fuente

```tree
src/
  sql/
    inc[remental]/
      000001-ddl.sql
      000002-dml.sql
      000003-dcl.sql
    rep[laceable]/
      procedure-foo.sql
      view-bar.sql
      package-baz.sql
buildr.yaml
```

Teniendo en `buildr.yaml`:

```yaml
system-id: "factura-blockchain"
application-id: "factura-blockchain-sql"
type: "oracle-sql-evolutional"
from:
  - "1.0.0"
  - "1.2.0"
```

En `src/sql/incremental/000001.sql`:

```sql
@@../replaceable/procedure-foo.sql
```
