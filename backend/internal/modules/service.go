package internal

import (
	"context"

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
}

type WeekRepository interface {
	GetByID(context.Context, uuid.UUID) (Week, error)
	ListByModuleRun(context.Context, uuid.UUID) ([]Week, error)
}

type ModuleService struct {
	moduleRepo    ModuleRepository
	moduleRunRepo ModuleRunRepository
	weekRepo      WeekRepository
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
