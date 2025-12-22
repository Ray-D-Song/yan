import { fetcher } from '../lib/fetcher'

import type { JsonValue } from '../lib/types'

// Note model
export interface Note {
  id: number
  parentId: number | null
  userId: number
  title: string
  content: string
  icon: string | null
  isFavorite: number
  position: number
  status: number
  createdAt: string
  updatedAt: string
}

// Note status constants
export const NoteStatus = {
  Normal: 1,
  Trashed: 2,
} as const

// Request types
export interface CreateNoteRequest {
  parentId?: number | null
  title: string
  content?: string
  icon?: string | null
  isFavorite?: number
  position?: number
  [key: string]: JsonValue | undefined
}

export interface UpdateNoteRequest {
  parentId?: number | null
  title: string
  content?: string
  icon?: string | null
  isFavorite?: number
  position?: number
  status?: number
  [key: string]: JsonValue | undefined
}

export interface UpdatePositionRequest {
  position: number
  [key: string]: JsonValue | undefined
}

export interface ListNotesParams {
  parentId?: number | null | 'null'
  status?: number
  favorite?: boolean
}

// API methods
export const notesApi = {
  /**
   * Create a new note
   * POST /api/v1/notes
   */
  create(data: CreateNoteRequest): Promise<Note> {
    return fetcher<Note>('/v1/notes', {
      method: 'POST',
      body: data,
    })
  },

  /**
   * Get note by ID
   * GET /api/v1/notes/:id
   */
  getById(id: number): Promise<Note> {
    return fetcher<Note>(`/v1/notes/${id}`, {
      method: 'GET',
    })
  },

  /**
   * List notes with optional filters
   * GET /api/v1/notes
   */
  list(params?: ListNotesParams): Promise<Note[]> {
    const searchParams = new URLSearchParams()

    if (params?.parentId !== undefined) {
      searchParams.append('parent_id', params.parentId === null ? 'null' : String(params.parentId))
    }

    if (params?.status !== undefined) {
      searchParams.append('status', String(params.status))
    }

    if (params?.favorite !== undefined) {
      searchParams.append('favorite', params.favorite ? 'true' : 'false')
    }

    const queryString = searchParams.toString()
    const url = queryString ? `/v1/notes?${queryString}` : '/v1/notes'

    return fetcher<Note[]>(url, {
      method: 'GET',
    })
  },

  listAll(): Promise<Note[]> {
    return this.list()
  },

  /**
   * Update note
   * PUT /api/v1/notes/:id
   */
  update(id: number, data: UpdateNoteRequest): Promise<Note> {
    return fetcher<Note>(`/v1/notes/${id}`, {
      method: 'PUT',
      body: data,
    })
  },

  /**
   * Delete note permanently
   * DELETE /api/v1/notes/:id
   */
  delete(id: number): Promise<null> {
    return fetcher<null>(`/v1/notes/${id}`, {
      method: 'DELETE',
    })
  },

  /**
   * Move note to trash (soft delete)
   * PUT /api/v1/notes/:id/trash
   */
  trash(id: number): Promise<null> {
    return fetcher<null>(`/v1/notes/${id}/trash`, {
      method: 'PUT',
    })
  },

  /**
   * Restore note from trash
   * PUT /api/v1/notes/:id/restore
   */
  restore(id: number): Promise<null> {
    return fetcher<null>(`/v1/notes/${id}/restore`, {
      method: 'PUT',
    })
  },

  /**
   * Toggle favorite status
   * PUT /api/v1/notes/:id/favorite
   */
  toggleFavorite(id: number): Promise<null> {
    return fetcher<null>(`/v1/notes/${id}/favorite`, {
      method: 'PUT',
    })
  },

  /**
   * Update note position
   * PUT /api/v1/notes/:id/position
   */
  updatePosition(id: number, position: number): Promise<null> {
    return fetcher<null>(`/v1/notes/${id}/position`, {
      method: 'PUT',
      body: { position },
    })
  },
}
