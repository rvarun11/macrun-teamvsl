package postgres

import (
	"fmt"
	"time"

	"github.com/CAS735-F23/macrun-teamvsl/challenge/config"
	"github.com/CAS735-F23/macrun-teamvsl/challenge/internal/core/domain"
	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type postgresChallenge struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"unique;not null"`
	Description string    `gorm:"not null"`
	BadgeURL    string
	Criteria    string `gorm:"not null"`
	Goal        float64
	Start       time.Time
	End         time.Time
	CreatedAt   time.Time
}

type postgresBadge struct {
	// ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	PlayerID    uuid.UUID `gorm:"not null;primaryKey"`
	ChallengeID uuid.UUID `gorm:"not null;primaryKey"`
	Score       float64
	CompletedOn time.Time
}

type postgresChallengeStats struct {
	// ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	PlayerID        uuid.UUID `gorm:"not null;primaryKey"`
	ChallengeID     uuid.UUID `gorm:"not null;primaryKey"`
	DistanceCovered float64
	EnemiesFought   uint8
	EnemiesEscaped  uint8
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(cfg *config.Postgres) *Repository {

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable client_encoding=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
		cfg.Encoding,
	)

	logLevel := getLogLevel(cfg.LogLevel)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
	})
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	db.AutoMigrate(&postgresChallenge{}, &postgresBadge{}, &postgresChallengeStats{})

	return &Repository{
		db: db,
	}
}

// Helper for converting to domain Challenge
func (pc *postgresChallenge) toAggregate() *domain.Challenge {
	return &domain.Challenge{
		ID:          pc.ID,
		Name:        pc.Name,
		Description: pc.Description,
		Criteria:    domain.Criteria(pc.Criteria),
		Goal:        pc.Goal,
		Start:       pc.Start,
		End:         pc.End,
		BadgeURL:    pc.BadgeURL,
		CreatedAt:   pc.CreatedAt,
	}
}

// Helper function to convert to domain Badge
func (pb *postgresBadge) toAggregate(ch *domain.Challenge) *domain.Badge {
	return &domain.Badge{
		PlayerID:    pb.PlayerID,
		Challenge:   ch,
		CompletedOn: pb.CompletedOn,
		Score:       pb.Score,
	}
}

// Helper function to convert to domain Badge
func (pcs *postgresChallengeStats) toAggregate(ch *domain.Challenge) *domain.ChallengeStats {
	return &domain.ChallengeStats{
		PlayerID:        pcs.PlayerID,
		Challenge:       ch,
		DistanceCovered: pcs.DistanceCovered,
		EnemiesFought:   pcs.EnemiesFought,
		EnemiesEscaped:  pcs.EnemiesEscaped,
	}
}

// getLogLevel returns the GORM Log Level
func getLogLevel(l string) gormLogger.LogLevel {
	switch l {
	case "silent":
		return gormLogger.Silent
	case "info":
		return gormLogger.Info
	case "error":
		return gormLogger.Error
	case "warn":
		return gormLogger.Warn
	default:
		return gormLogger.Warn
	}
}
