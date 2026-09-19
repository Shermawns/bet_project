CREATE TABLE users (
  id UUID PRIMARY KEY,
  nome TEXT NOT NULL,
  saldo_centavos BIGINT NOT NULL CHECK (saldo_centavos >= 0)
);

CREATE TABLE events (
  id UUID PRIMARY KEY,
  nome TEXT NOT NULL,
  odds_mil BIGINT NOT NULL CHECK (odds_mil >= 1000),
  status TEXT NOT NULL CHECK (status IN ('aberto','fechado','finalizado')),
  time_a_venceu BOOLEAN
);

CREATE TABLE bets (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  event_id UUID NOT NULL REFERENCES events(id),
  valor_apostado_centavos BIGINT NOT NULL CHECK (valor_apostado_centavos > 0),
  odd_no_momento_mil BIGINT NOT NULL CHECK (odd_no_momento_mil >= 1000),
  escolha_time_a BOOLEAN NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('pendente','ganhou','perdeu')),
  created_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX bets_user_id_idx ON bets(user_id);
CREATE INDEX bets_event_id_status_idx ON bets(event_id, status);
