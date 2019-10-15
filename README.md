# BDBU - Bûtent DB Utilities

## Versijos

-   [1.2.1]
    -   [+] Benchmark option

    -   [+] Tuner option
    
    -   [+] MySQL Option decoder


-   [1.1.0]

    -   [+] Kopijavimas vienos lentelës su parametru --table

    -   [- ] Panaikintas string tipo laukams trim funkcionalumas

-   [1.0.0] Pirminë versija, kopijuoja visas lenteles

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

arba vienai lentelei kopijuoti:

```cmd
bdbu copy --table klientai
```

> Pastaba: Ávykus bet kokiai klaidai kopijavimo procesas ið karto sustabdomas ir á ekranà iðvedamas klaidos praneðimas.
