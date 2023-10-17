package repository

import (
	"fmt"
	"sync"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/ports"

	"github.com/CAS735-F23/macrun-teamvs_/hrm/internal/core/domain"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	hrms map[uuid.UUID]domain.HRM
	sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		hrms: make(map[uuid.UUID]domain.HRM),
	}
}

func (r *MemoryRepository) AddHRMIntance(hrm domain.HRM) error {
	if r.hrms == nil {
		r.Lock()
		r.hrms = make(map[uuid.UUID]domain.HRM)
		r.Unlock()
	}

	if _, ok := r.hrms[hrm.GetID()]; ok {
		return fmt.Errorf("hrm already connected: %w", ports.ErrorCreateHRMFailed)
	}
	r.Lock()
	r.hrms[hrm.GetID()] = hrm
	r.Unlock()
	return nil

}

func (r *MemoryRepository) DeleteHRMInstance(hrmid uuid.UUID) error {

	if _, ok := r.hrms[hrmid]; !ok {
		return fmt.Errorf("hrm is not connected: %w", ports.ErrorCreateHRMFailed)
	}
	r.Lock()
	delete(r.hrms, hrmid)
	r.Unlock()
	return nil

}

func (r *MemoryRepository) Get(hrmId uuid.UUID) (*domain.HRM, error) {
	if hrm, ok := r.hrms[hrmId]; ok {
		return &hrm, nil
	}
	return &domain.HRM{}, ports.ErrorHRMNotFound
}

func (r *MemoryRepository) Update(hrm domain.HRM) error {
	if _, ok := r.hrms[hrm.GetID()]; !ok {
		return fmt.Errorf("hrm does not exist: %w", ports.ErrorUpdateHRMFailed)
	}
	r.Lock()
	r.hrms[hrm.GetID()] = hrm
	r.Unlock()
	return nil
}

func (r *MemoryRepository) List() ([]*domain.HRM, error) {
	if r.hrms == nil {
		// If r.workouts is nil, return an error or handle the case accordingly
		return nil, fmt.Errorf("hrms map doesn't exit %w", ports.ErrorListHRMSFailed)
	}
	hrms := make([]*domain.HRM, 0, len(r.hrms))
	for _, hrm := range r.hrms {
		hrms = append(hrms, &hrm)
	}
	return hrms, nil
}
