# Notes Frontend (React + Vite)

Consome a API Go em `http://localhost:8080`.

Design system: [Vercel-inspired](https://github.com/VoltAgent/awesome-design-md/blob/main/design-md/vercel/DESIGN.md) — tokens em `src/styles/vercel-tokens.css`, mesh gradient em `src/styles/mesh-gradient.css`.

## Rodar

Terminal 1 — backend:

```powershell
cd ..
go run ./cmd/api
```

Terminal 2 — frontend:

```powershell
npm run dev
```

Abra a URL que o Vite mostrar (geralmente `http://localhost:5173`).

## Variáveis

Copie `.env.example` para `.env` se precisar mudar a URL da API.
