// Profile represents a user in the system
export interface Profile {
    id: string;
    name: string;
    experience: number;
    level: number;
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
// Note: field names match Go struct serialization (PascalCase from database models)
export interface UserProgress {
    ID: string;
    UserID: string;
    ContentItemID: string;
    Completed: boolean;
    ProgressPct: number; // 0-100
    LastPosition?: { Int32: number; Valid: boolean }; // sql.NullInt32
    LastAccessed?: { Time: string; Valid: boolean }; // sql.NullTime
    CreatedAt?: { Time: string; Valid: boolean };
    UpdatedAt?: { Time: string; Valid: boolean };
}

// API response wrapper
export interface ApiResponse<T> {
    success: boolean;
    data?: T;
    message?: string;
    error?: string;
}

// Course existence check response
export interface CourseExistsResponse {
    exists: boolean;
    missing_paths: string[];
}

// Course directory info from filesystem scan
export interface CourseDirectory {
    path: string
    relative_path: string
    name: string
    size: number
    is_dir: boolean
}

// For creating courses during import
export interface CreateCourseInput {
    title: string
    relative_path: string
    description?: string
}

// Scan response
export interface ScanResponse {
    count: number
    directories: CourseDirectory[]
}

// Batch import response
export interface BatchImportResponse {
    success_count: number
    failure_count: number
    imported_courses: Course[]
    errors?: string[]
}

// Progress for a single module
export interface ModuleProgress {
    module_id: string
    user_id: string
    completed_items: number
    total_items: number
    completion_pct: number
    is_completed: boolean
    last_accessed_at?: string
}

// Progress for a single content item
export interface ContentItemProgress {
    content_item_id: string
    user_id: string
    completed: boolean
}

// Overall course progress
export interface CourseProgress {
    course_id: string
    user_id: string
    completed_modules: number
    total_modules: number
    completed_items: number
    total_items: number
    completion_pct: number
    is_completed: boolean
    last_accessed_at?: string
    modules?: ModuleProgress[]
    items?: ContentItemProgress[]
}

// Request to save video progress
export interface SaveProgressRequest {
    user_id: string
    last_position: number
    progress_pct: number
}

// Request to mark content as complete
export interface CompleteContentRequest {
    user_id: string
}

//XP Award Data
export interface XPAwardData {
    xp_awarded: boolean
    xp_amount: number
    leveled_up?: boolean
    old_level?: number
    new_level?: number
    gems_awarded?: number
}

export interface CompletionResponse {
    progress: UserProgress
    xp_awarded: boolean
    xp_amount: number
}

