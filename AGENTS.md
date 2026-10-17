# AGENTS.md

Pokyny pro AI agenty (Claude Code, Codex, Cursor…) i lidi, kteří v repozitáři pracují. Přehled projektu a spuštění najdeš v [README.md](README.md).

## Projekt v kostce

- Kvízová aplikace STEM_HUB, fork projektu [jovVix](https://github.com/Improwised/jovVix).
- `app/` je frontend v Nuxt 3 (Vue 3, Pinia, Tailwind v4, shadcn-vue).
- `api/` je backend v Go (Fiber, goqu, PostgreSQL, Redis, Ory Kratos).
- Dva režimy hry:
  - **Živý** (`live`): admin řídí postup přes websockety (`/socket/admin/arrange/:session_id`, `/socket/join/...`).
  - **Vlastním tempem** (`self_paced`): hráč si sám stahuje otázky přes REST (`GET /quiz/self/next`, `POST /quiz/answer`), čas měří server. Režim je uložený v `active_quizzes.mode`.

## Zásadní pravidla

1. **Žádné cloudové služby.** Nepřidávej AI, OAuth/Google přihlášení, externí API, analytiku, Google Fonts ani CDN. Všechno musí běžet lokálně. Fonty se berou z `@fontsource`.
2. **Čeština.** Všechny texty pro uživatele (UI, toasty, chyby z API v `api/constants`, e-maily) piš česky a tykej. Množná čísla řeš přes `czPlural` / `czCount` z `app/composables/plural.js`. Kód, názvy proměnných a komentáře zůstávají anglicky.
3. **Design jen přes tokeny.** Barvy, stíny a písma ber z `app/assets/css/main.css` (`jv-primary`, `jv-sky`, `jv-violet`, `jv-orange`, `jv-mauve`, `jv-ink`, `jv-canvas`, `jv-line`, `shadow-soft*`). Nepiš natvrdo hex barvy. Styl je moderní: zaoblené karty, tenké okraje `border-jv-line` a jemné stíny. Starý „brutalistický“ styl (silné černé okraje, natočení, washi pásky) už nepoužívej.
4. **Značka.** Logo je `app/public/images/stem-logo.png` a vedle něj nápis `STEM_HUB`. Titulky stránek mají tvar `… – STEM_HUB`. Kontakty: https://stemhub.cz a info@stemhub.cz. Patička je komponenta `components/common/SiteFooter.vue`.
5. **Oprávnění.** Tvorba kvízů je jen pro přihlášené adminy (role `admin-user`). Na straně frontendu skrývej odkazy (`v-if`), na backendu kontroluj middlewary.

## Frontend (`app/`)

- Composables a komponenty se importují automaticky (Nuxt). Ikony bereš z `lucide-vue-next`, toasty z `notivue`.
- UI komponenty shadcn přidávej příkazem `npx shadcn-vue@2.2.0 add <název>` a pak je uprav na tokeny.
- Kopírování do schránky dělej přes `composables/copy_to_clipboard.js`. Má záložní řešení pro HTTP v lokální síti, kde `navigator.clipboard` neexistuje.
- Kontrola: `npm run lint` musí projít bez varování. `npx vitest run`: na `main` dlouhodobě padá zhruba 33 testů. Nové testy musí procházet a počet padajících nesmí stoupnout.
- Dev server (`npm run dev`) je pomalý a stará otevřená záložka po restartu umí zlobit, takže ji obnov natvrdo. Na ověřování v prohlížeči použij `npx nuxt build && node .output/server/index.mjs`. Build neběží, když zároveň běží dev server.

## Backend (`api/`)

- Struktura: `routes/main.go` → `controllers/api/v1/*` → `models/*` (dotazy přes goqu, ne ručně psané SQL řetězce).
- Odpovědi vracej přes `utils.JSONSuccess` / `JSONFail` / `JSONError`. Texty chyb patří do `constants/constant.go`, česky.
- **Migrace:** nový soubor `database/migrations/<YYYYMMDDHHMMSS>_<název>.up.sql` a k němu `.down.sql`. Už aplikované migrace neupravuj.
- Když měníš endpointy, aktualizuj `utils/swagger_structs.go` a `assets/swagger.json`.
- Kontrola: `go vet ./... && go test ./...`. Bez lokálního Go:
  ```bash
  docker run --rm -v "$PWD":/src -v jovvix-gomod:/go/pkg/mod -w /src golang:1.24 \
    sh -c 'go build ./... && go vet ./... && go test ./...'
  ```

## Lokální prostředí

- Otevírej `http://127.0.0.1:…`, ne `localhost`. Cookies Kratosu mají doménu `127.0.0.1`.
- Kratos má vypnuté `haveibeenpwned`, protože jde o externí volání. Heslo nesmí obsahovat jméno ani e-mail uživatele, jinak ho Kratos odmítne.
- Přístupové údaje ani vlastní `.env` necommituj. Lokální hesla patří mimo repozitář.

## Git a PR

- Commity podle Conventional Commits (`feat:`, `fix:`, `docs:`, `refactor:` …). Husky před commitem spouští lint.
- Popisy PR piš česky: co se mění, proč a jak to bylo otestované.
- Neodesílej nic přímo do `main`. Větší změny rozděl do menších PR podle témat.
