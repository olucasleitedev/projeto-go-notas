import { useCallback, useEffect, useState, type FormEvent } from 'react'
import { auditApi } from './api/audit'
import { notesApi } from './api/notes'
import type { NoteEvent } from './types/audit'
import type { Note } from './types/note'
import './App.css'

const emptyForm = { title: '', content: '' }

function App() {
  const [notes, setNotes] = useState<Note[]>([])
  const [form, setForm] = useState(emptyForm)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [apiOnline, setApiOnline] = useState<boolean | null>(null)
  const [auditEvents, setAuditEvents] = useState<NoteEvent[]>([])

  const loadAudit = useCallback(async () => {
    try {
      const data = await auditApi.listEvents()
      setAuditEvents(data ?? [])
    } catch {
      setAuditEvents([])
    }
  }, [])

  const loadNotes = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await notesApi.list()
      setNotes(data ?? [])
      setApiOnline(true)
    } catch (err) {
      setApiOnline(false)
      setError(err instanceof Error ? err.message : 'Falha ao carregar notas')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadNotes()
    void loadAudit()
  }, [loadNotes, loadAudit])

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    if (!form.title.trim()) {
      setError('Título é obrigatório')
      return
    }

    setSaving(true)
    setError(null)
    try {
      if (editingId) {
        await notesApi.update(editingId, form)
      } else {
        await notesApi.create(form)
      }
      setForm(emptyForm)
      setEditingId(null)
      await loadNotes()
      await loadAudit()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao salvar')
    } finally {
      setSaving(false)
    }
  }

  function startEdit(note: Note) {
    setEditingId(note.id)
    setForm({ title: note.title, content: note.content })
    setError(null)
  }

  function cancelEdit() {
    setEditingId(null)
    setForm(emptyForm)
    setError(null)
  }

  async function handleDelete(id: string) {
    if (!confirm('Excluir esta nota?')) return

    setError(null)
    try {
      await notesApi.remove(id)
      if (editingId === id) cancelEdit()
      await loadNotes()
      await loadAudit()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao excluir')
    }
  }

  return (
    <>
      <nav className="nav-bar">
        <div className="nav-brand">
          <span className="nav-brand-mark" aria-hidden />
          Notes
        </div>
        <div className="nav-actions">
          <span
            className={`badge-secondary ${
              apiOnline === true ? 'online' : apiOnline === false ? 'offline' : ''
            }`}
          >
            <span className="dot" />
            {apiOnline === null ? 'Verificando API' : apiOnline ? 'API online' : 'API offline'}
          </span>
          <button
            type="button"
            className="btn btn-secondary"
            onClick={() => {
              void loadNotes()
              void loadAudit()
            }}
          >
            Atualizar
          </button>
        </div>
      </nav>

      <section className="feature-mesh-band mesh-backdrop">
        <div className="mesh-inner">
          <p className="eyebrow-mono">microservices · kafka · go</p>
          <h1 className="hero-title">Your notes, delivered.</h1>
          <p className="hero-lead">
            API Gateway, Notes Service e Audit Service desacoplados via eventos no Kafka.
          </p>
        </div>
      </section>

      {error && (
        <div className="showcase-band" style={{ paddingBottom: 0 }}>
          <div className="showcase-inner">
            <div className="alert" style={{ gridColumn: '1 / -1' }}>
              {error}
            </div>
          </div>
        </div>
      )}

      <section className="showcase-band">
        <div className="showcase-inner">
          <div className="card-marketing">
            <h2>{editingId ? 'Editar nota' : 'Nova nota'}</h2>
            <form onSubmit={handleSubmit} className="form">
              <label>
                Título
                <input
                  className="form-input"
                  value={form.title}
                  onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))}
                  placeholder="Ex: Deploy do backend"
                />
              </label>
              <label>
                Conteúdo
                <textarea
                  className="form-input"
                  value={form.content}
                  onChange={(e) => setForm((f) => ({ ...f, content: e.target.value }))}
                  placeholder="Detalhes da nota..."
                  rows={5}
                />
              </label>
              <div className="form-actions">
                <button type="submit" className="btn btn-primary" disabled={saving}>
                  {saving ? 'Salvando...' : editingId ? 'Atualizar' : 'Criar nota'}
                </button>
                {editingId && (
                  <button type="button" className="btn btn-secondary" onClick={cancelEdit}>
                    Cancelar
                  </button>
                )}
              </div>
            </form>
          </div>

          <div className="card-marketing">
            <div className="list-header">
              <h2>Suas notas</h2>
            </div>

            {loading ? (
              <p className="muted">Carregando...</p>
            ) : notes.length === 0 ? (
              <p className="muted">Nenhuma nota ainda. Crie a primeira ao lado.</p>
            ) : (
              <ul className="note-list">
                {notes.map((note) => (
                  <li key={note.id} className="note-card">
                    <div>
                      <h3>{note.title}</h3>
                      {note.content && <p>{note.content}</p>}
                      <time className="note-meta">
                        {new Date(note.updated_at).toLocaleString('pt-BR')}
                      </time>
                    </div>
                    <div className="card-actions">
                      <button type="button" className="btn btn-ghost" onClick={() => startEdit(note)}>
                        Editar
                      </button>
                      <button
                        type="button"
                        className="btn btn-danger"
                        onClick={() => void handleDelete(note.id)}
                      >
                        Excluir
                      </button>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>

        <div className="card-marketing audit-panel">
          <div className="list-header">
            <h2>Auditoria (Kafka)</h2>
            <span className="footer-mono">audit-service</span>
          </div>
          {auditEvents.length === 0 ? (
            <p className="muted">Nenhum evento ainda. Crie ou edite uma nota.</p>
          ) : (
            <ul className="note-list">
              {auditEvents.map((evt, index) => (
                <li key={`${evt.note_id}-${evt.timestamp}-${index}`} className="note-card">
                  <div>
                    <h3>{evt.type}</h3>
                    <p>{evt.title || evt.note_id}</p>
                    <time className="note-meta">
                      {new Date(evt.timestamp).toLocaleString('pt-BR')}
                    </time>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>
      </section>

      <footer className="footer">
        <div className="footer-inner">
          <span>Gateway · Notes · Audit · Kafka</span>
          <span className="footer-mono">projeto-go-notas</span>
        </div>
      </footer>
    </>
  )
}

export default App
