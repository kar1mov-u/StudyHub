package modules

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ModuleRepository interface {
	Create(context.Context, Module) error
	GetByID(context.Context, uuid.UUID) (Module, error)
	List(context.Context) ([]Module, error) //#TO-DO  implement better fetching later
	Delete(context.Context, uuid.UUID) error
}

type ModuleRunRepository interface {
	GetByID(context.Context, uuid.UUID) (ModuleRun, error)
	Create(context.Context, ModuleRun) error
	GetActiveByModuleID(context.Context, uuid.UUID) (ModuleRun, error)
	ListByModuleID(context.Context, uuid.UUID) ([]ModuleRun, error)
	DeactivateByModulID(context.Context, uuid.UUID) error
	Delete(context.Context, uuid.UUID) error
}

type WeekRepository interface {
	GetByID(context.Context, uuid.UUID) (Week, error)
	ListByModuleRun(context.Context, uuid.UUID) ([]Week, error)
}

type AcademicCalendarRepository interface {
	GetActive(context.Context) (AcademicTerm, error)
}

type ModuleService struct {
	moduleRepo    ModuleRepository
	moduleRunRepo ModuleRunRepository
	weekRepo      WeekRepository
	calendarRepo  AcademicCalendarRepository
}

func NewModuleService(moduleRepo ModuleRepository, weekRepo WeekRepository, moduleRunRepo ModuleRunRepository) *ModuleService {
	return &ModuleService{
		moduleRepo:    moduleRepo,
		moduleRunRepo: moduleRunRepo,
		weekRepo:      weekRepo,
	}
}

// GET modules/  -- returns the list of all modules, will be the default API to hit??
func (s *ModuleService) ListModules(ctx context.Context) ([]Module, error) {
	return s.moduleRepo.List(ctx)
}

// GET modules/<module_id>   -- returns the singe module info, not just module data, but ModuleData, LatestModuleRun, Weeks of that ModuleRun
func (s *ModuleService) GetModuleFull(ctx context.Context, id uuid.UUID) (ModulePage, error) {

	module, err := s.moduleRepo.GetByID(ctx, id)
	if err != nil {
		return ModulePage{}, err
	}

	moduleRun, err := s.moduleRunRepo.GetActiveByModuleID(ctx, id)
	if err != nil {
		return ModulePage{}, err
	}

	weeks, err := s.weekRepo.ListByModuleRun(ctx, moduleRun.ID)
	if err != nil {
		return ModulePage{}, err
	}

	return ModulePage{
		Module: module,
		Run:    moduleRun,
		Weeks:  weeks,
	}, nil
}

// when the new module is created, we will automatically create the newModuleRun also
func (s *ModuleService) CreateModule(ctx context.Context, module Module) error {
	//create the module
	if err := s.CreateModule(ctx, module); err != nil {
		return err
	}

	//get current semester
	term, err := s.calendarRepo.GetActive(ctx)
	if err != nil {
		return err
	}

	//create a modelRun
	moduleRun := ModuleRun{
		ID:        uuid.New(),
		ModuleID:  module.ID,
		Year:      term.Year,
		Semester:  term.Semester,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	return s.moduleRunRepo.Create(ctx, moduleRun)

}

func (s *ModuleService) DeleteModule(ctx context.Context, id uuid.UUID) error {
	return s.moduleRepo.Delete(ctx, id)
}

func (s *ModuleService) ListModuleRuns(ctx context.Context, moduleID uuid.UUID) ([]ModuleRun, error) {
	return s.moduleRunRepo.ListByModuleID(ctx, moduleID)
}

func (s *ModuleService) CreateModuleRun(ctx context.Context, moduleRun ModuleRun) error {
	return s.moduleRunRepo.Create(ctx, moduleRun)
}

func (s *ModuleService) GetModuleRun(ctx context.Context, id uuid.UUID) (ModuleRunPage, error) {
	moduleRun, err := s.moduleRunRepo.GetByID(ctx, id)
	if err != nil {
		return ModuleRunPage{}, err
	}

	weeks, err := s.weekRepo.ListByModuleRun(ctx, id)
	if err != nil {
		return ModuleRunPage{}, err
	}

	return ModuleRunPage{
		Run:   moduleRun,
		Weeks: weeks,
	}, nil
}

func (s *ModuleService) DeleteModuleRun(ctx context.Context, id uuid.UUID) error {
	return s.moduleRunRepo.Delete(ctx, id)
}
