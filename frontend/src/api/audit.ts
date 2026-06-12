import type { NoteEvent } from '../types/audit'

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

export const auditApi = {
  listEvents: () =>
    fetch(`${API_URL}/api/audit/events`).then(async (res) => {
      if (!res.ok) throw new Error('Falha ao carregar auditoria')
      return res.json() as Promise<NoteEvent[]>
    }),
}
