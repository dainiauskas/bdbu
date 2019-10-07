# BDBU - B�tent DB Utilities

## Konfigravimas

-   Source - nurodoma i� kokios duomen� baz�s bus kopijuojama
-   Destination - nurodoma � koki� duomen� baz� bus kopijuojama
    Source ir destination dialektai gali b�ti tiek mssql, tiek ir mysql

## Komandos

Pagalbos i�kvietimas:
```cmd
bdbu help
```

Versijos parodymas:
```cmd
bdbu version
```

Duomen� kopijavimas vykdomas:
```cmd
bdbu copy
```

> Pastaba: �vykus bet kokiai klaidai kopijavimo procesas i� karto sustabdomas ir � ekran� i�vedamas klaidos prane�imas.

## Paleidimas

Paleidus program� be joki� parametr� bus i�kviesta pagalba

```cmd
Usage:
  bdbu [command]

Available Commands:
  copy        copy database to another
  help        Help about any command
  version     Print the version number

Flags:
      --config string   config file (default is bdbu.yaml)
  -h, --help            help for bdbu
  -v, --verbose         verbose output

Use "bdbu [command] --help" for more information about a command.
```
