package services

import (
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/challenge_manager/internal/core/ports"
	"github.com/google/uuid"
)

type ChallengeService struct {
	repo ports.ChallengeRepository
}

// Factory for creating a new ChallengeService
func NewChallengeService(repo ports.ChallengeRepository) *ChallengeService {
	return &ChallengeService{
		repo: repo,
	}
}

// Challenge Services

func (svc *ChallengeService) CreateChallenge(req *domain.Challenge) (*domain.Challenge, error) {
	ch, err := domain.NewChallenge(req.Name, req.Description, req.BadgeURL, req.Criteria, req.Goal, req.Start, req.End)
	if err != nil {
		return &domain.Challenge{}, ports.ErrorCreateChallengeFailed
	}

	challenge, err := svc.repo.CreateChallenge(ch)
	if err != nil {
		return &domain.Challenge{}, ports.ErrorCreateChallengeFailed
	}

	return challenge, nil

}

func (svc *ChallengeService) GetChallengeByID(id uuid.UUID) (*domain.Challenge, error) {
	ch, err := svc.repo.GetChallengeByID(id)
	if err != nil {
		return &domain.Challenge{}, err
	}

	return ch, nil
}

func (svc *ChallengeService) UpdateChallenge(req *domain.Challenge) (*domain.Challenge, error) {
	ch, err := svc.repo.UpdateChallenge(req)
	if err != nil {
		return &domain.Challenge{}, nil
	}

	return ch, nil
}

func (svc *ChallengeService) ListChallenges(status string) ([]*domain.Challenge, error) {
	chs, err := svc.repo.ListChallenges()
	if err != nil {
		return []*domain.Challenge{}, err
	}

	if status == "active" {
		var activeChs []*domain.Challenge
		for _, ch := range chs {
			if ch.IsActive() {
				activeChs = append(activeChs, ch)
			}
		}
		return activeChs, nil
	}

	return chs, nil
}

// Badge Services

func (svc *ChallengeService) CreateBadge(cs *domain.ChallengeStats) (*domain.Badge, error) {
	b, err := domain.NewBadge(cs)
	if err != nil {
		return &domain.Badge{}, err
	}
	badge, err := svc.repo.CreateBadge(b)
	if err != nil {
		return &domain.Badge{}, err
	}

	return badge, nil
}

// TODO: This should be part of the Badge Service
func (svc *ChallengeService) ListBadgesByPlayerID(pid uuid.UUID) (*domain.Badge, error) {
	// b, err := domain.NewBadge(pid, ch, score)
	// if err != nil {
	// 	return &domain.Badge{}, err
	// }
	// badge, err := svc.repo.CreateBadge(b)
	// if err != nil {
	// 	return &domain.Badge{}, err
	// }

	return &domain.Badge{}, nil
}

func (svc *ChallengeService) SubscribeToActiveChallenges(pid uuid.UUID, dc float32, ef uint8, ee uint8) error {
	// TODO: This should be in the context
	activeChs, err := svc.ListChallenges("active")
	if err != nil {
		return err
	}

	for _, ch := range activeChs {
		cs := domain.NewChallengeStats(ch, pid, dc, ef, ee)
		err = svc.repo.CreateOrUpdateChallengeStats(cs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *ChallengeService) ListChallengeStatsByPlayerID(pid uuid.UUID) ([]*domain.ChallengeStats, error) {
	csArr, err := svc.repo.ListChallengeStatsByPlayerID(pid)
	if err != nil {
		return []*domain.ChallengeStats{}, err
	}
	return csArr, nil
}

// This function runs when a challenge ends
func (svc *ChallengeService) XXXXXXX(ch *domain.Challenge) ([]*domain.Badge, error) {
	// 1. Fetch Challenge Stats
	// 2. Create Badge
	// 3. Delete all Challenge Stats
}
