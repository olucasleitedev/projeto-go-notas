import type { CreateNoteInput, Note } from '../types/note'

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })

  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: 'request failed' }))
    throw new Error(body.error ?? `HTTP ${response.status}`)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return response.json() as Promise<T>
}

export const notesApi = {
  list: () => request<Note[]>('/api/notes'),
  create: (input: CreateNoteInput) =>
    request<Note>('/api/notes', { method: 'POST', body: JSON.stringify(input) }),
  update: (id: string, input: CreateNoteInput) =>
    request<Note>(`/api/notes/${id}`, { method: 'PUT', body: JSON.stringify(input) }),
  remove: (id: string) => request<void>(`/api/notes/${id}`, { method: 'DELETE' }),
  health: () => request<{ status: string }>('/health'),
}
