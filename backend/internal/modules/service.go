package modules

import (
	"context"
	"log/slog"
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
	GetLatestModuleRun(context.Context, uuid.UUID) (ModuleRun, error)
	ListByModuleID(context.Context, uuid.UUID) ([]ModuleRun, error)
	Delete(context.Context, uuid.UUID) error
}

type WeekRepository interface {
	GetByID(context.Context, uuid.UUID) (Week, error)
	ListByModuleRun(context.Context, uuid.UUID) ([]Week, error)
	CreateWeeksForMoudleRun(context.Context, uuid.UUID) error
}

type AcademicCalendarRepository interface {
	GetActive(context.Context) (AcademicTerm, error)
	Create(context.Context, AcademicTerm) error
	List(context.Context) ([]AcademicTerm, error)
	DeActivate(ctx context.Context) error
}

type ModuleService struct {
	moduleRepo    ModuleRepository
	moduleRunRepo ModuleRunRepository
	weekRepo      WeekRepository
	calendarRepo  AcademicCalendarRepository
}

func NewModuleService(moduleRepo ModuleRepository, weekRepo WeekRepository, moduleRunRepo ModuleRunRepository, calendarRepo AcademicCalendarRepository) *ModuleService {
	return &ModuleService{
		moduleRepo:    moduleRepo,
		moduleRunRepo: moduleRunRepo,
		weekRepo:      weekRepo,
		calendarRepo:  calendarRepo,
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

	moduleRun, err := s.moduleRunRepo.GetLatestModuleRun(ctx, id)
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

	//get current semester
	term, err := s.calendarRepo.GetActive(ctx)
	if err != nil {
		return err
	}

	//create the module
	if err := s.moduleRepo.Create(ctx, module); err != nil {
		return err
	}

	//create a modelRun
	moduleRun := ModuleRun{
		ID:        uuid.New(),
		ModuleID:  module.ID,
		Year:      term.Year,
		Semester:  term.Semester,
		CreatedAt: time.Now(),
	}

	err = s.moduleRunRepo.Create(ctx, moduleRun)
	if err != nil {
		return err
	}

	return s.weekRepo.CreateWeeksForMoudleRun(ctx, moduleRun.ID)

}

func (s *ModuleService) DeleteModule(ctx context.Context, id uuid.UUID) error {
	return s.moduleRepo.Delete(ctx, id)
}

func (s *ModuleService) ListModuleRuns(ctx context.Context, moduleID uuid.UUID) ([]ModuleRun, error) {
	return s.moduleRunRepo.ListByModuleID(ctx, moduleID)
}

func (s *ModuleService) CreateModuleRun(ctx context.Context, moduleRun ModuleRun) error {
	err := s.moduleRunRepo.Create(ctx, moduleRun)
	if err != nil {
		return err
	}
	return s.weekRepo.CreateWeeksForMoudleRun(ctx, moduleRun.ID)
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

// Academic Calendar methods

func (s *ModuleService) GetActiveAcademicTerm(ctx context.Context) (AcademicTerm, error) {
	return s.calendarRepo.GetActive(ctx)
}

func (s *ModuleService) ListAcademicTerms(ctx context.Context) ([]AcademicTerm, error) {
	return s.calendarRepo.List(ctx)
}

func (s *ModuleService) StartNewTerm(ctx context.Context, term AcademicTerm) error {
	//1 deactive all the other terms
	err := s.calendarRepo.DeActivate(ctx)
	if err != nil {
		return err
	}

	//2 create new academic term
	err = s.calendarRepo.Create(ctx, term)
	if err != nil {
		return err
	}
	//3 create moduleRun for each module
	//get the list of modules
	modules, err := s.moduleRepo.List(ctx)
	if err != nil {
		return err
	}
	for _, module := range modules {
		moduleRun := ModuleRun{
			ID:        uuid.New(),
			ModuleID:  module.ID,
			Year:      term.Year,
			Semester:  term.Semester,
			CreatedAt: time.Now(),
		}
		err = s.moduleRunRepo.Create(ctx, moduleRun)
		if err != nil {
			//if we cant create the 1 module, should still continue the processs
			slog.Error("failed to create the moduleRun in NewTerm", "err", err)
		}
		err = s.weekRepo.CreateWeeksForMoudleRun(ctx, moduleRun.ID)
		if err != nil {
			slog.Error("failed to create the weeks in NewTerm", "err", err)
		}

	}
	return nil
}
