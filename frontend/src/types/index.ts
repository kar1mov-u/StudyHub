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
  IsActive: boolean
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
  IsActive: boolean
}

// Request DTOs
export interface CreateModuleRequest {
  code: string
  name: string
  department_name: string
}

export interface CreateModuleRunRequest {
  year: number
  semester: string
  is_active: boolean
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

// Response DTOs
export interface CreateResponse {
  id: string
}
