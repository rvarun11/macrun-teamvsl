package httphandler

import "github.com/google/uuid"

type StartWorkout struct {
	// TrailID chosen by the Player
	TrailID uuid.UUID `json:"trail_id"`
	// PlayerID of the player starting the workout session
	PlayerID uuid.UUID `json:"player_id"`
	// Whether HRM is connected or not
	HRMConnected bool `json:"hrm_connected"`
	// If HRM is connected then HRM ID otherwise garbage
	HRMId uuid.UUID `json:"hrm_id"`
	// HardCore Mode of User
	HardCoreMode bool `json:"hardcore_mode"`
}

type StartWorkoutOption struct {
	// WorkoutID for which the workout option is to be stopped
	Option string `json:"option"`
}

type StopWorkoutOption struct {
	// WorkoutID for which the workout option is to be stopped
	WorkoutID uuid.UUID `json:"workout_id"`
}
