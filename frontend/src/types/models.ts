// Profile represents a user in the system
export interface Profile {
  id: string;
  name: string;
  experience: number;
  gems: number;
  streak: number;
  last_active_date?: string;
  created_at?: string;
  updated_at?: string;
}

// Course represents a complete learning course
export interface Course {
  id: string;
  title: string;
  description: string;
  creator?: string;
  creator_id?: string;
  relative_path: string;
  modules?: Module[];
  created_at?: string;
  updated_at?: string;
}

// Module represents a section within a course
export interface Module {
  id: string;
  course_id: string;
  title: string;
  description?: string;
  relative_path: string;
  order: number;
  content_items?: ContentItem[];
  created_at?: string;
  updated_at?: string;
}

// ContentItem represents individual learning content (videos, PDFs, etc.)
// This is what gives XP rewards when completed
export interface ContentItem {
  id: string;
  module_id: string;
  title: string;
  description?: string;
  relative_path: string;
  content_type: string; // video, pdf, text, etc.
  duration?: number; // seconds (for videos)
  size?: number; // file size in bytes
  order: number;
  created_at?: string;
  updated_at?: string;
}

// Task represents a background job for async operations (batch imports, etc.)
export interface Task {
  id: string;
  type: string; // e.g., "batch_import"
  status: 'pending' | 'processing' | 'completed' | 'failed';
  progress: number; // 0-100 percent
  created_at: string;
  started_at?: string;
  completed_at?: string;
  message?: string;
  error_message?: string;
  result?: unknown;
}

// UserProgress tracks how far a user has gotten through content
export interface UserProgress {
  id: string;
  user_id: string;
  content_item_id: string;
  completed: boolean;
  progress_pct: number; // 0-100
  last_position?: number; // seconds (for videos)
  last_accessed?: string;
  created_at?: string;
  updated_at?: string;
}

// API response wrapper
export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}
