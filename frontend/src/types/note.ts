export interface Note {
  id: string
  title: string
  content: string
  created_at: string
  updated_at: string
}

export interface CreateNoteInput {
  title: string
  content: string
}
