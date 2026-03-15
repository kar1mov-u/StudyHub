// API Response wrapper
export interface ApiResponse<T> {
  data: T
}

export interface ApiError {
  error: {
    code: number
    message: string
  }
}

// Module types
export interface Module {
  ID: string
  Code: string
  Name: string
  DepartmentName: string
  CreatedAt: string
  UpdatedAt: string
}

export interface ModuleRun {
  ID: string
  ModuleID: string
  Year: number
  Semester: string
  CreatedAt: string
}

export interface Week {
  ID: string
  ModuleRunID: string
  Number: number
}

export interface ModulePage {
  Module: Module
  Run: ModuleRun
  Weeks: Week[]
}

export interface ModuleRunPage {
  Run: ModuleRun
  Weeks: Week[]
}

// Academic Term types
export interface AcademicTerm {
  ID: string
  Year: number
  Semester: string
}

// Request DTOs
export interface CreateModuleRequest {
  code: string
  name: string
  department_name: string
}

export interface UpdateModuleRequest {
  code?: string
  name?: string
  department_name?: string
}

export interface CreateAcademicTermRequest {
  year: number
  semester: string
}

// Resource types
export type ResourceType = 'file' | 'link' | 'note'

export interface Resource {
  ID: string
  WeekID: string
  UserID: string
  UserName: string
  ResourceType: ResourceType
  Hash: string
  Name: string
  Url: string
  ObjectID: string
  CreatedAt: string
  UpdatedAt: string
}

export interface UserResource {
  ID: string
  WeekID: string
  UserID: string
  ModuleName: string
  Semester: string
  Year: number
  WeekNumber: number
  ObjectID: string | null
  ExternalLink: string | null
  ResourceType: ResourceType
  Name: string
  CreatedAt: string
}

// User types
export interface User {
  ID: string
  Email: string
  FirstName: string
  LastName: string
  IsAdmin: boolean
  CreatedAt: string
  UpdatedAt: string
}

// Auth types
export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
}

export interface RegisterRequest {
  email: string
  password: string
  first_name: string
  last_name: string
}

// Flashcard types
export interface Flashcard {
  ID: string
  ObjectID: string | null
  UserID: string | null
  WeekID: string | null
  Front: string
  Back: string
}

// Comment types
export interface Comment {
  id: string
  user_id: string
  week_id: string
  reply_id?: string | null
  content: string
  upvote: number
  downvote: number
  created_at: string
}

export interface CreateCommentRequest {
  user_id: string
  week_id: string
  content: string
  reply_id?: string
}

// Response DTOs
export interface CreateResponse {
  id: string
}
