# BDBU - B�tent DB Utilities

## Versijos
-   [1.3.0]
  - [+]  �d�ta apsauga nuo nety�inio duomen� kopijavimo.
         Naudojant copy funkcij� ir norint, kad b�t� perra�ytos lentel�s b�tina
         nurodyti parametr� -D arba --with-drop
  - [+]  Su parametru -e arba --empty kopijuojamos tu��ios lentel�s, i�skyrus
         B�tent programos sistemines lenteles.

-   [1.2.2]

    -   [+] Konfig�racijos klaid� prane�imai

-   [1.2.1]

    -   [+] Benchmark option

    -   [+] Tuner option

    -   [+] MySQL Option decoder


-   [1.1.0]

    -   [+] Kopijavimas vienos lentel�s su parametru --table

    -   [- ] Panaikintas string tipo laukams trim funkcionalumas

-   [1.0.0] Pirmin� versija, kopijuoja visas lenteles

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

arba vienai lentelei kopijuoti:

```cmd
bdbu copy --table klientai
```

norint perra�yti lentel�:

```cmd
bdbu copy --table klientai -D
```

> Pastaba: �vykus bet kokiai klaidai kopijavimo procesas i� karto sustabdomas ir � ekran� i�vedamas klaidos prane�imas.
