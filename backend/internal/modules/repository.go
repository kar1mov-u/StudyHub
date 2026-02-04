package modules

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ModuleRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewModuleRepositoryPostgres(p *pgxpool.Pool) *ModuleRepositoryPostgres {
	return &ModuleRepositoryPostgres{pool: p}
}

func (r *ModuleRepositoryPostgres) Create(ctx context.Context, module Module) error {
	query := `INSERT INTO modules (id, code, name, department_name) VALUES ($1, $2, $3, $4 )`
	_, err := r.pool.Exec(ctx, query, module.ID, module.Code, module.Name, module.DepartmentName)
	if err != nil {
		return fmt.Errorf("InsertModule err: %w", err)
	}
	return nil
}

func (r *ModuleRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (Module, error) {
	var module Module
	query := `SELECT id, code, name, department_name, created_at, updated_at FROM modules WHERE id=$1`
	row := r.pool.QueryRow(ctx, query, id)
	err := row.Scan(&module.ID, &module.Code, &module.Name, &module.DepartmentName, &module.CreatedAt, &module.UpdatedAt)
	if err != nil {
		return Module{}, fmt.Errorf("GetModule err: %w", err)
	}
	return module, nil
}

func (r *ModuleRepositoryPostgres) List(ctx context.Context) ([]Module, error) {
	modules := make([]Module, 0)
	query := `SELECT id, code, name, department_name FROM modules`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return []Module{}, fmt.Errorf("ListModules query err: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var module Module

		err := rows.Scan(&module.ID, &module.Code, &module.Name, &module.DepartmentName)
		if err != nil {
			return []Module{}, fmt.Errorf("ListModules scan err: %w", err)
		}

		modules = append(modules, module)
	}

	return modules, nil
}

func (r *ModuleRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM modules WHERE id=$1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("DeleteModules err: %w", err)
	}
	return nil
}

type ModuleRunRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewModuleRunRepositoryPostgres(p *pgxpool.Pool) *ModuleRunRepositoryPostgres {
	return &ModuleRunRepositoryPostgres{pool: p}
}

func (r *ModuleRunRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (ModuleRun, error) {
	var moduleRun ModuleRun

	query := `SELECT id, module_id, year, semester, created_at FROM module_runs WHERE id=$1`
	row := r.pool.QueryRow(ctx, query, id)
	err := row.Scan(&moduleRun.ID, &moduleRun.ModuleID, &moduleRun.Year, &moduleRun.Semester, &moduleRun.CreatedAt)
	if err != nil {
		return ModuleRun{}, fmt.Errorf("GetModuleRun err: %w", err)
	}
	return moduleRun, nil
}

func (r *ModuleRunRepositoryPostgres) Create(ctx context.Context, moduleRun ModuleRun) error {
	query := `INSERT INTO module_runs (id, module_id, year, semester, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.pool.Exec(ctx, query, moduleRun.ID, moduleRun.ModuleID, moduleRun.Year, moduleRun.Semester, moduleRun.CreatedAt)
	if err != nil {
		return fmt.Errorf("InsertModuleRun err: %w", err)
	}
	return nil
}

func (r *ModuleRunRepositoryPostgres) GetLatestModuleRun(ctx context.Context, moduleID uuid.UUID) (ModuleRun, error) {
	var moduleRun ModuleRun

	query := `SELECT id, module_id, year, semester, created_at FROM module_runs WHERE module_id=$1 ORDER BY created_at DESC LIMIT 1`
	row := r.pool.QueryRow(ctx, query, moduleID)
	err := row.Scan(&moduleRun.ID, &moduleRun.ModuleID, &moduleRun.Year, &moduleRun.Semester, &moduleRun.CreatedAt)
	if err != nil {
		return ModuleRun{}, fmt.Errorf("GetLatestModuleRun err: %w", err)
	}
	return moduleRun, nil
}

func (r *ModuleRunRepositoryPostgres) ListByModuleID(ctx context.Context, moduleID uuid.UUID) ([]ModuleRun, error) {
	moduleRuns := make([]ModuleRun, 0)
	query := `SELECT id, module_id, year, semester, created_at FROM module_runs WHERE module_id=$1`
	rows, err := r.pool.Query(ctx, query, moduleID)
	if err != nil {
		return []ModuleRun{}, fmt.Errorf("ListModuleRuns query err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var moduleRun ModuleRun
		err := rows.Scan(&moduleRun.ID, &moduleRun.ModuleID, &moduleRun.Year, &moduleRun.Semester, &moduleRun.CreatedAt)
		if err != nil {
			return []ModuleRun{}, fmt.Errorf("ListModuleRuns scan err: %w", err)
		}
		moduleRuns = append(moduleRuns, moduleRun)
	}

	return moduleRuns, nil
}

func (r *ModuleRunRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM module_runs WHERE id=$1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("DeleteModuleRun err: %w", err)
	}
	return nil
}

type WeekRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewWeekRepositoryPostgres(p *pgxpool.Pool) *WeekRepositoryPostgres {
	return &WeekRepositoryPostgres{pool: p}
}

func (r *WeekRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (Week, error) {
	var week Week

	query := `SELECT id, module_run_id, number FROM weeks WHERE id=$1`
	row := r.pool.QueryRow(ctx, query, id)
	err := row.Scan(&week.ID, &week.ModuleRunID, &week.Number)
	if err != nil {
		return Week{}, fmt.Errorf("GetWeek err: %w", err)
	}
	return week, nil
}

func (r *WeekRepositoryPostgres) ListByModuleRun(ctx context.Context, id uuid.UUID) ([]Week, error) {
	weeks := make([]Week, 0)
	query := `SELECT id, module_run_id, number FROM weeks WHERE module_run_id=$1`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return []Week{}, fmt.Errorf("ListWeeks query err: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var week Week
		err := rows.Scan(&week.ID, &week.ModuleRunID, &week.Number)
		if err != nil {
			return []Week{}, fmt.Errorf("ListWeeks scan: %w", err)
		}

		weeks = append(weeks, week)
	}

	return weeks, err
}

func (r *WeekRepositoryPostgres) CreateWeeksForMoudleRun(ctx context.Context, moduleRunID uuid.UUID) error {
	batch := pgx.Batch{}
	query := `INSERT INTO weeks(module_run_id, id, number) VALUES ($1, $2, $3)`
	for i := 1; i <= 15; i++ {
		batch.Queue(query, moduleRunID, uuid.New(), i)
	}

	err := r.pool.SendBatch(ctx, &batch).Close()
	if err != nil {
		return fmt.Errorf("CreateWeeks batch err: %w", err)
	}

	return nil
}

type AcademicCalendarRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewAcademicCalendarRepositoryPostgres(p *pgxpool.Pool) *AcademicCalendarRepositoryPostgres {
	return &AcademicCalendarRepositoryPostgres{pool: p}
}

func (r *AcademicCalendarRepositoryPostgres) GetActive(ctx context.Context) (AcademicTerm, error) {
	var term AcademicTerm

	query := `SELECT id, year, semester, is_active FROM academic_terms WHERE is_active=true`
	row := r.pool.QueryRow(ctx, query)
	err := row.Scan(&term.ID, &term.Year, &term.Semester, &term.IsActive)
	if err != nil {
		return AcademicTerm{}, fmt.Errorf("GetActiveAcademicTerm err: %w", err)
	}
	return term, nil
}

func (r *AcademicCalendarRepositoryPostgres) Create(ctx context.Context, term AcademicTerm) error {
	query := `INSERT INTO academic_terms (id, year, semester, is_active) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, term.ID, term.Year, term.Semester, term.IsActive)
	if err != nil {
		return fmt.Errorf("InsertAcademicTerm err: %w", err)
	}
	return nil
}

func (r *AcademicCalendarRepositoryPostgres) List(ctx context.Context) ([]AcademicTerm, error) {
	terms := make([]AcademicTerm, 0)
	query := `SELECT id, year, semester, is_active FROM academic_terms`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return []AcademicTerm{}, fmt.Errorf("ListAcademicTerms query err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var term AcademicTerm
		err := rows.Scan(&term.ID, &term.Year, &term.Semester, &term.IsActive)
		if err != nil {
			return []AcademicTerm{}, fmt.Errorf("ListAcademicTerms scan err: %w", err)
		}
		terms = append(terms, term)
	}

	return terms, nil
}

func (r *AcademicCalendarRepositoryPostgres) DeActivate(ctx context.Context) error {
	query := `UPDATE academic_terms SET is_active=false`
	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("DeactivateAcademicTerm err: %w", err)
	}
	return nil
}

func (r *AcademicCalendarRepositoryPostgres) Activate(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE academic_terms SET is_active=true WHERE id=$1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ActivateAcademicTerm err: %w", err)
	}
	return nil
}
