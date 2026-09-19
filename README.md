# bet-api

API simples em Go para cadastrar usuários, eventos e apostas. Usa Uber Fx para montar as dependências e PostgreSQL em Docker para guardar os dados.

## Iniciar o PostgreSQL com Docker

Suba somente o banco:

```bash
docker compose up -d
```

O container cria o banco `bet_api` e executa automaticamente a migration inicial na primeira vez que o volume é criado.

Use este comando para acompanhar a inicialização:

```bash
docker compose logs -f postgres
```

O banco fica disponível para a API local nesta URL:

```text
postgres://postgres:postgres@localhost:5432/bet_api?sslmode=disable
```

Para parar o banco sem apagar os dados:

```bash
docker compose down
```

Para apagar o banco e criá-lo novamente do zero, incluindo a execução da migration inicial:

```bash
docker compose down -v
docker compose up -d
```

## Iniciar a API

```bash
go run ./cmd
```

A API fica em `http://localhost:8000`. `GET /health` verifica se está respondendo.

## Rotas

- `POST /users` cria usuário; saldo é enviado em centavos.
- `GET /users/{id}` consulta usuário.
- `POST /events` cria evento.
- `GET /events` lista eventos; pode filtrar com `?status=aberto`.
- `PATCH /events/{id}` altera odd ou status.
- `POST /bets` cria aposta simples.
- `GET /bets` lista apostas.
- `GET /bets/{id}` consulta aposta.
- `PATCH /bets/{id}` altera valor ou escolha da aposta.
- `DELETE /bets/{id}` exclui aposta.

## Exemplo de aposta

```json
{
  "user_id": "UUID_DO_USUARIO",
  "event_id": "UUID_DO_EVENTO",
  "valor_apostado_centavos": 1000,
  "escolha_time_a": true
}
```

`true` significa escolher o time A. Ao criar a aposta, a API copia a odd atual do evento. Este exemplo CRUD não debita nem credita saldo e não calcula vencedores; ele só demonstra criar, consultar, editar e apagar apostas.
