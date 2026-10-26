# Jak přispívat do STEM_HUB kvízů

Díky, že chceš pomoct! Tady najdeš, jak na to.

## Kodex chování

Projekt se řídí [kodexem chování](./CODE_OF_CONDUCT.md). Účastí v projektu souhlasíš s jeho dodržováním.

## Začínáme

1. Projdi [otevřené issues](https://github.com/Mirek321/jovVix-stemQuiz/issues), nebo založ nové s popisem chyby či nápadu.
2. Napiš do issue, že na něm pracuješ.
3. Naklonuj repozitář a rozjeď projekt podle [README](./README.md#vývoj-ze-zdrojáků).
4. Vytvoř si větev z `main` (viz [Názvy větví](#názvy-větví)).
5. Udělej změny, ověř je a pošli pull request.

Před prací si přečti také [AGENTS.md](./AGENTS.md). Jsou v něm konvence projektu: žádné cloudové služby, texty česky, design přes tokeny `jv-*` a další.

## Názvy větví

Anglicky, malými písmeny a se spojovníky, krátce podle obsahu:

| Typ změny | Příklad |
|---|---|
| nová funkce | `self-paced-mode` |
| oprava chyby | `fix-scoreboard-calculation` |
| dokumentace | `docs-readme-agents` |
| úklid, údržba | `cleanup-unused-files` |

## Commity

Píšeme podle [Conventional Commits](https://www.conventionalcommits.org/):

```
<typ>(<oblast>): <popis>
```

Typy: `feat` (nová funkce), `fix` (oprava), `docs` (dokumentace), `test` (testy), `refactor` (úprava kódu bez změny chování), `chore` (údržba), `ci` (CI/CD), `style` (formátování).

Příklady:

```
feat(quiz): režim vlastním tempem
fix(reports): oprava výpočtu skóre u záporných bodů
docs(readme): nový návod na spuštění
```

Před každým commitem se přes husky spouští frontendový lint.

## Pull requesty

1. Název a popis PR piš **anglicky**: co se mění, proč a jak to bylo otestované.
2. Propoj PR s issue (`Closes #123`).
3. Přiřaď PR @Mirek321 a požádej ho o review.
4. Zapracuj připomínky z review a pushni nové commity.
5. Po schválení PR sloučí správce.
6. Větší změny rozděl do menších PR podle témat.

### Kontrolní seznam

- [ ] Kód odpovídá konvencím projektu (viz [AGENTS.md](./AGENTS.md))
- [ ] Lint a testy lokálně prochází
- [ ] Nepřibyla žádná nová varování
- [ ] Texty pro uživatele jsou česky
- [ ] Dokumentace je aktualizovaná (pokud je potřeba)

## Testy a lint

```bash
# frontend
cd app
npm run lint           # kontrola
npm run lint-fix       # automatické opravy
npx vitest run         # testy

# backend
cd api
go vet ./...
go test ./...
golangci-lint run
```

Na `main` dlouhodobě padá asi 33 frontendových testů, které hledají starou podobu stránek. Nové testy musí procházet a počet padajících nesmí stoupnout.

## Hlášení chyb a nápady

Založ [issue](https://github.com/Mirek321/jovVix-stemQuiz/issues/new) a uveď:

- **u chyby:** postup, jak ji vyvolat, co jsi čekal a co se stalo, prostředí (systém, prohlížeč) a případně screenshot
- **u nápadu:** jaký problém řeší, navrhované řešení a případné alternativy

Bezpečnostní chyby nehlas veřejně, postupuj podle [SECURITY.md](./SECURITY.md).

## Licence příspěvků

Projekt vychází z [jovVix](https://github.com/Improwised/jovVix) (© Improwised Technologies a přispěvatelé) a je šířen pod licencí [GNU AGPL 3.0](./LICENSE.txt). Odesláním příspěvku souhlasíš, že bude šířen pod stejnou licencí.

---

Tento návod vychází z průvodce pro přispěvatele projektu jovVix.
