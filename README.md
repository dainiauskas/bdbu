# BDBU - Bûtent DB Utilities

## Konfigravimas

-   Source - nurodoma ið kokios duomenø bazës bus kopijuojama
-   Destination - nurodoma á kokià duomenø bazæ bus kopijuojama
    Source ir destination dialektai gali bûti tiek mssql, tiek ir mysql

## Komandos

Pagalbos iðkvietimas:
```cmd
bdbu help
```

Versijos parodymas:
```cmd
bdbu version
```

Duomenø kopijavimas vykdomas:
```cmd
bdbu copy
```

> Pastaba: Ávykus bet kokiai klaidai kopijavimo procesas ið karto sustabdomas ir á ekranà iðvedamas klaidos praneðimas.

## Paleidimas

Paleidus programà be jokiø parametrø bus iðkviesta pagalba

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
