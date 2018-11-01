# DESIGN

## Tipos de proyecto a implementar

- oracle-sql-evolutional ✔
- oracle-sql-deferred ✔
- oracle-sql-eventual ✔
- microsoft-net-web
- microsoft-net-web-core
- microsoft-net-lib
- go-web
- go-command

## Determinación de identificador de versión desde una WC

Se debe determinar la versión correspondiente a la WC, para eso se basa en
información disponible en git:

- Determina si la WC está limpia (DIRTY=false) o tiene cambios (DIRTY=true)
- Busca el último tag del branch actual (TAG)
- Busca el último commit del branch actual (COMMIT)
- VERSION=$(git describe --abbrev=40 HEAD)[1:]
- Y si DIRTY => VERSION = VERSION + "-dirty"
