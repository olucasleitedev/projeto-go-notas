# projeto-go-notas

API de notas em **arquitetura de microserviços** — Gateway, serviços independentes e mensageria com **Kafka** (Redpanda).

**Lucas Leite** · [github.com/olucasleitedev](https://github.com/olucasleitedev)

---

## Arquitetura

```
                    ┌─────────────┐
                    │   React     │
                    │  :5173      │
                    └──────┬──────┘
                           │ HTTP
                    ┌──────▼──────┐
                    │   Gateway   │  :8080  (entrada única + CORS)
                    └──────┬──────┘
              ┌────────────┼────────────┐
              │            │            │
       ┌──────▼─────┐      │     ┌──────▼──────┐
       │   Notes    │      │     │   Audit     │
       │  Service   │      │     │  Service    │
       │   :8081    │      │     │   :8082     │
       └──────┬─────┘      │     └──────▲──────┘
              │ publish    │            │ consume
              └────────────┼────────────┘
                    ┌──────▼──────┐
                    │    Kafka    │  topic: note.events
                    │  (Redpanda) │
                    └─────────────┘
```

| Serviço | Porta | Responsabilidade |
|---------|-------|------------------|
| **gateway** | 8080 | Proxy reverso, CORS, ponto único para o frontend |
| **notes-service** | 8081 | CRUD de notas, publica eventos no Kafka |
| **audit-service** | 8082 | Consome eventos e expõe trilha de auditoria |
| **redpanda** | 19092 | Broker Kafka-compatible |

---

## Eventos (Kafka)

Tópico `note.events`:

| Evento | Quando |
|--------|--------|
| `note.created` | Nota criada |
| `note.updated` | Nota atualizada |
| `note.deleted` | Nota excluída |

O **audit-service** consome esses eventos de forma assíncrona — desacoplado do fluxo HTTP principal.

---

## Stack

| Camada | Tecnologia |
|--------|------------|
| Gateway / serviços | Go 1.22+, `net/http` |
| Mensageria | Kafka (Redpanda) · `segmentio/kafka-go` |
| Frontend | React + TypeScript + Vite |
| Infra local | Docker Compose |

---

## Rodar com Docker (recomendado)

```powershell
cd "C:\Users\Lucas M2Z Creative\Documents\estudos-golang"
docker compose up --build
```

Serviços disponíveis:

- Gateway: http://localhost:8080
- Notes: http://localhost:8081
- Audit: http://localhost:8082

### Frontend

```powershell
cd frontend
npm install
npm run dev
```

→ http://localhost:5173

---

## Rodar local (sem Docker nos serviços Go)

**Terminal 1 — Kafka (só o broker):**

```powershell
docker compose up redpanda -d
```

**Terminal 2 — Notes:**

```powershell
$env:KAFKA_BROKERS="localhost:19092"
go run ./cmd/notes-service
```

**Terminal 3 — Audit:**

```powershell
$env:KAFKA_BROKERS="localhost:19092"
go run ./cmd/audit-service
```

**Terminal 4 — Gateway:**

```powershell
go run ./cmd/gateway
```

---

## Estrutura do repositório

```
cmd/
├── gateway/           # API Gateway
├── notes-service/     # microserviço de notas
├── audit-service/     # microserviço de auditoria
└── api/               # monólito legado (dev simples)
internal/              # domínio, use cases, handlers
pkg/
├── events/            # contratos de eventos
└── messaging/         # Kafka producer/consumer
frontend/              # React
docker-compose.yml
```

---

## API (via Gateway :8080)

| Método | Rota | Serviço |
|--------|------|---------|
| `GET` | `/health` | gateway |
| `GET/POST` | `/api/notes` | notes-service |
| `GET/PUT/DELETE` | `/api/notes/{id}` | notes-service |
| `GET` | `/api/audit/events` | audit-service |

---

## Variáveis de ambiente

| Variável | Padrão | Uso |
|----------|--------|-----|
| `KAFKA_BROKERS` | `localhost:19092` | Brokers Kafka |
| `KAFKA_ENABLED` | `true` | `false` desliga mensageria |
| `NOTES_SERVICE_URL` | `http://localhost:8081` | Gateway → notes |
| `AUDIT_SERVICE_URL` | `http://localhost:8082` | Gateway → audit |

---

## Roadmap

- [x] Microserviços (gateway, notes, audit)
- [x] Kafka / event-driven
- [x] Docker Compose
- [x] Frontend com painel de auditoria
- [ ] Postgres por serviço
- [ ] Testes de integração
- [ ] CI/CD
- [ ] Observabilidade (logs estruturados, métricas)

---

## Autor

**Lucas Leite** — [olucasleitedev](https://github.com/olucasleitedev)
