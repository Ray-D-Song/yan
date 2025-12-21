-- Migration: user_table
-- Created at: 2025-12-15 11:21:53
-- Description: Create user table
-- Write your UP migration here
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  status INTEGER NOT NULL DEFAULT 1, -- 1 normal, 0 disable
  is_admin INTEGER NOT NULL DEFAULT 0, -- 1 true, 0 false
  created_at TIMESTAMP NOT NULL DEFAULT (datetime ('now')),
  updated_at TIMESTAMP NOT NULL DEFAULT (datetime ('now'))
);

