package services

import (
	"fmt"
	"time"

	logger "github.com/CAS735-F23/macrun-teamvsl/challenge/log"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/domain"
	"github.com/CAS735-F23/macrun-teamvsl/workout/internal/core/ports"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/umahmood/haversine"
)

type ActiveWorkoutsLastLocation struct {
	// Latitude of the Player
	Latitude float64 `json:"latitude"`
	// Longitude of the Player
	Longitude float64 `json:"longitude"`
	// Time of location
	TimeOfLocation time.Time `json:"time_of_location"`
}

type ActiveWorkoutsHeartRate struct {
	// HRM Connection
	HRMConnected bool
}

type WorkoutService struct {
	repo                       ports.WorkoutRepository
	peripheral                 ports.PeripheralDeviceClient
	user                       ports.UserServiceClient
	activeWorkoutsLastLocation map[uuid.UUID]ActiveWorkoutsLastLocation
	activeWorkoutsHeartRate    map[uuid.UUID]ActiveWorkoutsHeartRate
}

// Factory for creating a new WorkoutService
func NewWorkoutService(repo ports.WorkoutRepository, peripheral ports.PeripheralDeviceClient, user ports.UserServiceClient) *WorkoutService {
	return &WorkoutService{
		repo:                       repo,
		peripheral:                 peripheral,
		user:                       user,
		activeWorkoutsLastLocation: make(map[uuid.UUID]ActiveWorkoutsLastLocation),
		activeWorkoutsHeartRate:    make(map[uuid.UUID]ActiveWorkoutsHeartRate),
	}
}

func (s *WorkoutService) List() ([]*domain.Workout, error) {
	return []*domain.Workout{}, nil
}

func (s *WorkoutService) GetWorkout(id uuid.UUID) (*domain.Workout, error) {
	workout, err := s.repo.GetWorkout(id)
	if err != nil {
		logger.Error("failed to get workout", zap.String("id", id.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get workout with ID %s: %w", id, err)
	}
	return workout, nil
}

func (s *WorkoutService) Start(workout *domain.Workout, HRMID uuid.UUID, HRMConnected bool) (string, error) {
	// Retrieve user profile details
	profile, err := s.user.GetWorkoutPreferenceOfUser(workout.PlayerID)
	if err != nil {
		logger.Error("failed to get user profile", zap.String("playerID", workout.PlayerID.String()), zap.Error(err))
		return "", fmt.Errorf("failed to get profile for user %s: %w", workout.PlayerID, err)
	}

	if err != nil {
		logger.Error("failed to get hardcore mode for user", zap.String("playerID", workout.PlayerID.String()), zap.Error(err))
		return "", fmt.Errorf("failed to get hardcore mode for user %s: %w", workout.PlayerID, err)
	}

	// Set workout profile
	workout.Profile = profile

	// Update workout details in the repository
	_, err = s.repo.UpdateWorkout(workout)
	if err != nil {
		logger.Error("failed to update workout", zap.String("workoutID", workout.WorkoutID.String()), zap.Error(err))
		return "", fmt.Errorf("failed to update workout %s: %w", workout.WorkoutID, err)
	}

	shelterNeeded := !workout.HardcoreMode

	err = s.peripheral.BindPeripheralData(workout.PlayerID, workout.WorkoutID, HRMID, HRMConnected, shelterNeeded)
	if err != nil {
		logger.Error("failed to bind peripheral data", zap.String("HRMID", HRMID.String()), zap.Error(err))
		return "", fmt.Errorf("failed to bind HRM device %s for workout %s: %w", HRMID, workout.WorkoutID, err)
	}

	// Record the heart rate monitor connection status
	s.activeWorkoutsHeartRate[workout.WorkoutID] = ActiveWorkoutsHeartRate{
		HRMConnected: HRMConnected,
	}
	var currentWorkoutOption int8
	if shelterNeeded {
		currentWorkoutOption = 7
	} else {
		currentWorkoutOption = 6
	}

	// Create workout options for the new workout
	workoutOptions := &domain.WorkoutOptions{
		WorkoutID:             workout.WorkoutID,
		CurrentWorkoutOption:  currentWorkoutOption,
		FightsPushDown:        false,
		IsWorkoutOptionActive: false,
	}

	// Create the workout and its options in the repository
	err = s.repo.Create(workout, workoutOptions)
	if err != nil {
		logger.Error("failed to create workout", zap.String("workoutID", workout.WorkoutID.String()), zap.Error(err))
		return "", fmt.Errorf("failed to create workout %s: %w", workout.WorkoutID, err)
	}

	// Log the successful creation of the workout
	logger.Info("Workout created with ID", zap.String("workoutID", workout.WorkoutID.String()))

	// Generate and return the link to the workout options
	linkURL := fmt.Sprintf("/workoutOptions?workoutID=%s", workout.WorkoutID)
	return linkURL, nil // Return nil explicitly to indicate no error occurred
}

func (s *WorkoutService) GetWorkoutOptions(workoutID uuid.UUID) ([]domain.WorkoutOptionLink, error) {
	// Retrieve workout options from the repository
	pworkoutOptions, err := s.repo.GetWorkoutOptions(workoutID)
	if err != nil {
		logger.Debug("failed to get workout options", zap.String("workoutID", workoutID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get workout options for workout %s: %w", workoutID, err)
	}

	// If no options are found for the given workout, handle the condition appropriately
	if pworkoutOptions == nil {
		logger.Debug("no workout options found", zap.String("workoutID", workoutID.String()))
		return nil, fmt.Errorf("no workout options found for workout %s", workoutID)
	}

	// Determine the order of workout options based on FightsPushDown
	optionsOrder := computeOptionsOrder(pworkoutOptions)

	// Generating HATEOAS links for StartWorkoutOption based on the computed order
	links := generateStartWorkoutOptionLinks(workoutID, optionsOrder)

	return links, nil
}

// Constants for bit positions
const (
	ShelterBit = 0
	FightBit   = 1
	EscapeBit  = 2
)

func getWorkoutType(bit int8) string {
	if bit == ShelterBit {
		return "Shelter"
	} else if bit == FightBit {
		return "Fight"
	}
	return "Escape"
}

// Computing the order of options based on FightsPushDown and whether the Shelter bit is set.
func computeOptionsOrder(pworkoutOptions *domain.WorkoutOptions) []uint8 {
	order := []uint8{}

	// Add Shelter to the order only if the bit is set for the current workout option
	if pworkoutOptions.CurrentWorkoutOption&ShelterBit != 0 {
		order = append(order, ShelterBit)
	}

	// Add Fight and Escape to the order. If FightsPushDown is set, Escape comes before Fight.
	if pworkoutOptions.FightsPushDown {
		order = append(order, EscapeBit, FightBit)
	} else {
		order = append(order, FightBit, EscapeBit)
	}

	return order
}

func generateStartWorkoutOptionLinks(workoutID uuid.UUID, optionsOrder []uint8) []domain.WorkoutOptionLink {
	var links []domain.WorkoutOptionLink // A slice to hold ordered links

	for _, option := range optionsOrder {
		var optionName string
		switch option {
		case ShelterBit:
			optionName = "Shelter"
		case FightBit:
			optionName = "Fight"
		case EscapeBit:
			optionName = "Escape"
		default:
			continue
		}
		linkName := "option" + optionName
		linkURL := fmt.Sprintf("/startWorkoutOption?workoutID=%s&option=%d", workoutID, option)
		links = append(links, domain.WorkoutOptionLink{Name: linkName, URL: linkURL})
	}

	return links
}

func (s *WorkoutService) UpdateDistanceTravelled(workoutID uuid.UUID, latitude float64, longitude float64, timeOfLocation time.Time) error {
	// Check if the workout ID exists in the location map
	lastLocation, locationExists := s.activeWorkoutsLastLocation[workoutID]

	if locationExists {
		// Calculate the distance between existing and new location
		distanceCovered := 0.0
		if lastLocation.Latitude != latitude || lastLocation.Longitude != longitude {
			point1 := haversine.Coord{Lat: lastLocation.Latitude, Lon: lastLocation.Longitude}
			point2 := haversine.Coord{Lat: latitude, Lon: longitude}

			// Calculate the distance using the Haversine formula
			_, distanceCovered = haversine.Distance(point1, point2)
		}

		// Update the workout distance if the distance covered is greater than 0
		if distanceCovered > 0 {

			s.activeWorkoutsLastLocation[workoutID] = ActiveWorkoutsLastLocation{
				Latitude:       latitude,
				Longitude:      longitude,
				TimeOfLocation: timeOfLocation,
			}

			// Get the workout from the repository
			workout, err := s.repo.GetWorkout(workoutID)
			if err != nil {
				return err // Propagate the error from the repository
			}

			// Update the workout distance
			workout.DistanceCovered += distanceCovered

			// Update the workout in the repository
			_, err = s.repo.UpdateWorkout(workout)
			if err != nil {
				return err // Propagate the error from the repository
			}
		}
	} else {
		// If the location doesn't exist, add it to the map
		workout, err := s.repo.GetWorkout(workoutID)

		if err == nil && !workout.IsCompleted {
			s.activeWorkoutsLastLocation[workoutID] = ActiveWorkoutsLastLocation{
				Latitude:       latitude,
				Longitude:      longitude,
				TimeOfLocation: timeOfLocation,
			}
		}
	}

	return nil // Return nil to indicate success
}

func (s *WorkoutService) UpdateShelter(workoutID uuid.UUID, DistanceToShelter float64) error {
	// Get the workout options from the repository
	workoutOptions, err := s.repo.GetWorkoutOptions(workoutID)
	if err != nil {
		return err // Propagate the error from the repository
	}

	workoutOptions.DistanceToShelter = DistanceToShelter

	s.repo.UpdateWorkoutOptions(workoutOptions)

	if err != nil {
		return err // Propagate the error from the repository
	}

	return nil // Return nil to indicate success
}

func (s *WorkoutService) StartWorkoutOption(workoutID uuid.UUID, workoutType uint8) error {
	// Get the workout options from the repository
	workoutOptions, err := s.repo.GetWorkoutOptions(workoutID)
	if err != nil {
		return err // Propagate the error from the repository
	}

	// Check if the workout option is already active
	if workoutOptions.IsWorkoutOptionActive {
		return domain.ErrWorkoutOptionAlreadyActive // Return an error with a custom message
	}

	// Update the workout option to make it active (you need to set appropriate fields)
	workoutOptions.IsWorkoutOptionActive = true
	workoutOptions.CurrentWorkoutOption = int8(workoutType)

	// Update the workout options in the repository
	_, err = s.repo.UpdateWorkoutOptions(workoutOptions)

	if err != nil {
		return err // Propagate the error from the repository
	}
	logger.Info("Workout Option OPT started for ID : workoutID", zap.String("workoutID", workoutOptions.WorkoutID.String()), zap.String("OPT", getWorkoutType(workoutOptions.CurrentWorkoutOption)))
	return nil // Return nil to indicate success
}

func (s *WorkoutService) StopWorkoutOption(workoutID uuid.UUID) error {
	// Get the workout from the repository
	workout, err := s.repo.GetWorkout(workoutID)
	if err != nil {
		logger.Error("failed to get workout for stopping options", zap.String("workoutID", workoutID.String()), zap.Error(err))
		return fmt.Errorf("failed to get workout %s for stopping options: %w", workoutID, err)
	}

	// Get the workout options from the repository
	workoutOptions, err := s.repo.GetWorkoutOptions(workoutID)
	if err != nil {
		logger.Error("failed to get workout options for workout", zap.String("workoutID", workoutID.String()), zap.Error(err))
		return fmt.Errorf("failed to get workout options for workout %s: %w", workoutID, err)
	}

	// Check if the workout option is already inactive
	if !workoutOptions.IsWorkoutOptionActive {
		logger.Info("workout option already inactive", zap.String("workoutID", workoutID.String()))
		return domain.ErrWorkoutOptionAlreadyInActive
	}

	// Update the Shelters, Fights and Escapes
	if workoutOptions.CurrentWorkoutOption == ShelterBit {
		workout.Shelters++
	} else if workoutOptions.CurrentWorkoutOption == FightBit {
		workout.Fights++
	} else if workoutOptions.CurrentWorkoutOption == EscapeBit {
		workout.Escapes++
	}

	// Update the workout option to make it inactive
	workoutOptions.IsWorkoutOptionActive = false
	workoutOptions.CurrentWorkoutOption = -1

	// Update the workout options in the repository
	_, err = s.repo.UpdateWorkoutOptions(workoutOptions)
	if err != nil {
		logger.Error("failed to update workout options on stop", zap.String("workoutID", workoutID.String()), zap.Error(err))
		return fmt.Errorf("failed to update workout options for workout %s on stop: %w", workoutID, err)
	}

	// Update the workout in the repository
	_, err = s.repo.UpdateWorkout(workout)
	if err != nil {
		logger.Error("failed to update workout on stop", zap.String("workoutID", workoutID.String()), zap.Error(err))
		return fmt.Errorf("failed to update workout %s on stop: %w", workoutID, err)
	}

	// Log the successful stopping of the workout option
	workoutType := getWorkoutType(workoutOptions.CurrentWorkoutOption) // Assuming getWorkoutType is a valid function
	logger.Info("Workout option stopped", zap.String("workoutID", workoutID.String()), zap.String("optionType", workoutType))

	return nil
}

func (s *WorkoutService) Stop(id uuid.UUID) (*domain.Workout, error) {
	// Retrieve the workout to be stopped
	tempWorkout, err := s.repo.GetWorkout(id)
	if err != nil {
		logger.Error("failed to get workout for stopping", zap.String("workoutID", id.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get workout %s for stopping: %w", id, err)
	}

	// Set the workout as completed and mark the end time
	tempWorkout.EndedAt = time.Now()
	tempWorkout.IsCompleted = true

	// Update the workout's status in the repository
	_, err = s.repo.UpdateWorkout(tempWorkout)
	if err != nil {
		logger.Debug("failed to update workout on stop", zap.String("workoutID", tempWorkout.WorkoutID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to update workout %s on stop: %w", tempWorkout.WorkoutID, err)
	}

	// Delete the workout options associated with the workout
	err = s.repo.DeleteWorkoutOptions(tempWorkout.WorkoutID)
	if err != nil {
		logger.Debug("failed to delete workout options on stop", zap.String("workoutID", tempWorkout.WorkoutID.String()), zap.Error(err))
		// need not return the error
	}

	// Remove the workout from active workouts tracking
	delete(s.activeWorkoutsLastLocation, tempWorkout.WorkoutID)
	delete(s.activeWorkoutsHeartRate, tempWorkout.WorkoutID)

	// Unbind peripheral data associated with the workout
	err = s.peripheral.UnbindPeripheralData(tempWorkout.WorkoutID)
	if err != nil {
		logger.Debug("failed to unbind peripheral data on stop", zap.String("workoutID", tempWorkout.WorkoutID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to unbind peripheral data on stop")
	}

	return tempWorkout, nil
}

func (s *WorkoutService) GetDistanceById(workoutID uuid.UUID) (float64, error) {
	return s.repo.GetDistanceByID(workoutID)
}

func (s *WorkoutService) GetDistanceCoveredBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (float64, error) {
	return s.repo.GetDistanceCoveredBetweenDates(playerID, startDate, endDate)
}

func (s *WorkoutService) GetEscapesMadeById(workoutID uuid.UUID) (uint16, error) {
	return s.repo.GetEscapesMadeByID(workoutID)
}

func (s *WorkoutService) GetEscapesMadeBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error) {
	return s.repo.GetEscapesMadeBetweenDates(playerID, startDate, endDate)
}

func (s *WorkoutService) GetFightsFoughtById(workoutID uuid.UUID) (uint16, error) {
	return s.repo.GetFightsFoughtByID(workoutID)
}

func (s *WorkoutService) GetFightsFoughtBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error) {
	return s.repo.GetFightsFoughtBetweenDates(playerID, startDate, endDate)
}

func (s *WorkoutService) GetSheltersTakenById(workoutID uuid.UUID) (uint16, error) {
	return s.repo.GetSheltersTakenByID(workoutID)
}

func (s *WorkoutService) GetSheltersTakenBetweenDates(playerID uuid.UUID, startDate time.Time, endDate time.Time) (uint16, error) {
	return s.repo.GetSheltersTakenBetweenDates(playerID, startDate, endDate)
}

// Function to run periodically
func (s *WorkoutService) RunPeriodicTask() {
	for {
		// Iterate through active workouts in locationMap
		for workoutID := range s.activeWorkoutsLastLocation {
			_ = s.ComputeWorkoutOptionsOrder(workoutID)
			// TODO
			//if err != nil {
			// Handle the error (e.g., log it)
			//}
		}

		// Sleep for 4 seconds before the next iteration
		time.Sleep(4 * time.Second)
	}
}

// ComputeWorkoutOptionsOrder is modified to take profile directly
func (s *WorkoutService) ComputeWorkoutOptionsOrder(workoutID uuid.UUID) error {
	// Get the workout options from the repository
	workoutOptions, err := s.repo.GetWorkoutOptions(workoutID)
	if err != nil {
		return err // Propagate the error from the repository
	}

	// Get the workout from the repository
	workout, err := s.repo.GetWorkout(workoutID)
	if err != nil {
		return err // Propagate the error from the repository
	}

	// Calculate the weight based on cardio
	weight := calculateWeight(workout.Profile)

	// Get the number of fights and escapes
	fights, err := s.repo.GetFightsFoughtByID(workoutID)
	if err != nil {
		return err
	}
	escapes, err := s.repo.GetEscapesMadeByID(workoutID)
	if err != nil {
		return err
	}

	// Calculate the weight based on fights and escapes
	fightEscapeWeight := calculateFightEscapeWeight(fights, escapes)

	avgHeartRate, err := s.peripheral.GetAverageHeartRateOfUser(workout.WorkoutID)

	if err != nil {
		return err
	}

	age, err := s.user.GetUserAge(workout.PlayerID)

	if err != nil {
		return err
	}

	// Calculate the weight based on average heart rate
	heartRateWeight := calculateHeartRateWeight(avgHeartRate, age)

	// Sum the weights
	totalWeight := weight + fightEscapeWeight + heartRateWeight
	// Update FightsPushDown based on the total weight
	if totalWeight >= 75 {
		workoutOptions.FightsPushDown = true
	} else {
		workoutOptions.FightsPushDown = false
	}

	// Update the workout options in the repository
	_, err = s.repo.UpdateWorkoutOptions(workoutOptions)
	if err != nil {
		return err // Propagate the error from the repository
	}

	return nil // Return nil to indicate success
}

// Helper function to calculate the weight based on cardio
func calculateWeight(profile string) int {
	if profile == "cardio" {
		return 50
	}
	return 0
}

// Helper function to calculate the weight based on fights and escapes
func calculateFightEscapeWeight(fights, escapes uint16) int {
	if fights-escapes >= 2 {
		return 25
	}
	return 0
}

// Helper function to calculate the weight based on average heart rate
func calculateHeartRateWeight(avgHeartRate uint8, age uint8) int {
	// TODO: Fetch user age from profile or another service
	// age := profile.Age

	// Mock user age for testing

	maxHeartRate := 220 - age
	percentageMaxHeartRate := float64(avgHeartRate) / float64(maxHeartRate) * 100

	if percentageMaxHeartRate < 70 {
		return 25
	}

	return 0
}
