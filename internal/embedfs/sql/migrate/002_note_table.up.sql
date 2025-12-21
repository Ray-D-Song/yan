-- Migration: note_table
-- Created at: 2025-12-21 19:39:54
-- Description: Create notes table
-- Write your UP migration here
CREATE TABLE IF NOT EXISTS notes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  parent_id INTEGER, -- If null, it indicates the root
  user_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  icon TEXT, -- emoji or icon identifier
  is_favorite INTEGER NOT NULL DEFAULT 0, -- 1 favorite, 0 not favorite
  position INTEGER NOT NULL DEFAULT 0, -- for custom ordering
  status INTEGER NOT NULL DEFAULT 1, -- 1 normal, 0 trashed
  created_at TIMESTAMP NOT NULL DEFAULT (datetime ('now')),
  updated_at TIMESTAMP NOT NULL DEFAULT (datetime ('now'))
);
