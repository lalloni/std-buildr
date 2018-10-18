# Proyecto Oracle SQL Diferido

## Estructura de código fuente

```tree
src/
  sql/
    control-diario.sql
    control-mensual.sql
    control-eventual.sql
buildr.yaml
```

Teniendo en `buildr.yaml`:

```yaml
system-id: "factura-blockchain"
application-id: "factura-blockchain-sql-process"
type: "oracle-sql-deferred"
```
