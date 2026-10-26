# Bezpečnostní zásady

## Podporované verze

Bezpečnostní opravy dostává jen aktuální verze ve větvi `main`.

## Nahlášení zranitelnosti

**Bezpečnostní chyby prosím nehlas přes veřejné GitHub issues.**

Pošli je e-mailem na **info@stemhub.cz** s předmětem „Bezpečnost“. Ozveme se co nejdřív.

Do hlášení uveď:

- typ problému (např. SQL injection, XSS, obejití přihlášení)
- cesty ke zdrojovým souborům, kterých se problém týká
- umístění kódu (větev, commit nebo přímý odkaz)
- případné zvláštní nastavení potřebné k vyvolání
- postup krok za krokem, jak problém vyvolat
- ukázku nebo proof-of-concept (pokud je to možné)
- dopad: jak by mohl útočník chybu zneužít

Díky těmto informacím hlášení rychleji posoudíme.

## Co můžeš čekat

- **Potvrzení:** dáme vědět, že jsme hlášení dostali.
- **Posouzení:** ověříme zranitelnost a její dopad.
- **Oprava:** připravíme opravu.
- **Zveřejnění:** termín zveřejnění s tebou domluvíme.

## Bezpečný přístav

Proti výzkumníkům nepodnikneme právní kroky, pokud:

- se v dobré víře snaží neporušit soukromí, nezničit data a nenarušit provoz
- pracují jen s vlastními účty nebo s výslovným souhlasem jejich majitele
- zranitelnost nezneužijí víc, než je nutné k ověření, že existuje
- chybu nahlásí včas a nezveřejní ji dřív, než je k dispozici oprava

## Působnost

Zásady platí pro:

- aplikaci STEM_HUB kvízy (frontend i backend)
- konfiguraci nasazení v Dockeru
- API endpointy a přihlašování

Zásady neplatí pro:

- služby a závislosti třetích stran (hlas je jejich správcům)
- sociální inženýrství vůči správcům nebo přispěvatelům

## Chyby v původním jovVixu

Pokud se zranitelnost týká i původního projektu [jovVix](https://github.com/Improwised/jovVix), nahlas ji také jeho autorům podle jejich bezpečnostních zásad (security@improwised.com).

---

Tyto zásady vycházejí z bezpečnostních zásad projektu jovVix (© Improwised Technologies a přispěvatelé).
