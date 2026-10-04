# Zadanie 1: `Sender.Send` — wyślij webhook i zmierz wynik

**Kto:** ja implementuję. Agent przygotował stronę, handler `/send`, odbiorcę demo i testy.

## Problem

Panel ma już formularz, ale kliknięcie „Send” pokazuje błąd `Send: not implemented yet`.
Brakuje jednej funkcji: `Send` w `internal/webhook/sender.go`, która wysyła JSON
POST-em do odbiorcy i opisuje, co się stało.

## Efekt w aplikacji

Po zadaniu na `http://localhost:8080`:

- `demo-ok` → zielone pole, status 200, czas w ms, odpowiedź `{"received":true}`;
- `demo-fail` → czerwone pole, status 500 i treść błędu odbiorcy;
- `demo-slow` → „no response”, czas ~3 s i błąd timeoutu.

## Kryteria ukończenia

1. `go test ./...` przechodzi (5 testów w `internal/webhook/sender_test.go`).
2. `./scripts/check.sh` bez uwag (gofmt, vet, golangci-lint, testy).
3. Trzy scenariusze wyżej działają w przeglądarce.
4. Potrafisz w 2–3 zdaniach wyjaśnić, dlaczego HTTP 500 **nie** jest błędem zwracanym przez `Send`, a timeout jest.

## Co ma zrobić `Send`

1. Ograniczyć czas całej próby do `s.Timeout`.
2. Zbudować żądanie POST z `payload` w body i nagłówkiem `Content-Type: application/json`.
3. Wysłać je klientem `s.Client` i zmierzyć czas.
4. Przy błędzie wysyłki zwrócić `Result` z ustawionym `Duration` oraz błąd opakowany kontekstem.
5. Przy odpowiedzi przeczytać najwyżej `s.MaxResponseBytes` bajtów body, ustawić `Truncated`
   i zawsze zamknąć body.

## Rzeczy do sprawdzenia w dokumentacji

- `context.WithTimeout` i dlaczego zawsze robi się `defer cancel()`;
- `http.NewRequestWithContext`;
- `io.LimitReader` (podpowiedź: jak poznać, że body było dłuższe niż limit?);
- `fmt.Errorf` z `%w`.

## Gdy utkniesz

Napisz w wątku, który punkt blokuje. Dostaniesz najpierw wskazówkę, potem pseudokod,
a pełne rozwiązanie tylko na prośbę.

## Uruchamianie

```sh
go test ./...                 # testy
go run ./cmd/web              # serwer na http://localhost:8080
```
