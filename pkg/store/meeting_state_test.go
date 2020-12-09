package store

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"gitlab.com/infra.run/public/b3scale/pkg/bbb"
)

func meetingStateFactory(pool *pgxpool.Pool, init *MeetingState) (*MeetingState, error) {
	// We start with a blank meeting
	if init == nil {
		init = &MeetingState{
			ID:         uuid.New().String(),
			InternalID: uuid.New().String(),
		}
	}
	// A meeting should not exist without a backend and
	// frontend.
	if init.frontend == nil {
		init.frontend = frontendStateFactory(pool)
		if err := init.frontend.Save(); err != nil {
			return nil, err
		}
		init.FrontendID = &init.frontend.ID
	}
	if init.backend == nil {
		init.backend = backendStateFactory(pool)
		if err := init.backend.Save(); err != nil {
			return nil, err
		}
		init.BackendID = &init.backend.ID
	}

	// Prepare meeting state
	if init.Meeting == nil {
		init.Meeting = &bbb.Meeting{
			MeetingName: "MyMeetingName-" + uuid.New().String(),
		}
	}
	init.Meeting.MeetingID = init.ID
	init.Meeting.InternalMeetingID = init.InternalID

	return InitMeetingState(pool, init), nil
}

func TestGetMeetingStates(t *testing.T) {
	pool := connectTest(t)

	m1, err := meetingStateFactory(pool, &MeetingState{
		ID:         uuid.New().String(),
		InternalID: uuid.New().String(),
		Meeting: &bbb.Meeting{
			Running: true,
		}})
	t.Log(m1)
	if err != nil {
		t.Error(err)
		return
	}
	err = m1.Save()
	if err != nil {
		t.Error(err)
		return
	}

	// Get running meetings
	states, err := GetMeetingStates(pool, Q().
		Where("id = ?", m1.ID).
		Where("state->'Running' = ?", true))
	if err != nil {
		t.Error(err)
	}
	if len(states) != 1 {
		t.Error("expected meeting to be in result set")
	}

}

func TestMeetingStateSave(t *testing.T) {
	pool := connectTest(t)

	state, err := meetingStateFactory(pool, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("meeting state after factory:", state)
	t.Log(
		"backendID:", state.BackendID,
		"frontendID:", state.FrontendID)

	err = state.Save()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("New meeting state id:", state.ID)
}

func TestMeetingStateSaveUpdate(t *testing.T) {
	pool := connectTest(t)
	state, err := meetingStateFactory(pool, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if err := state.Save(); err != nil {
		t.Error(err)
		return
	}

	state.Meeting = &bbb.Meeting{
		MeetingName: "bar",
	}
	if err := state.Save(); err != nil {
		t.Error(err)
		return
	}
}

func TestMeetingStateIsStale(t *testing.T) {
	pool := connectTest(t)
	state, err := meetingStateFactory(pool, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("before SavE()")
	if err := state.Save(); err != nil {
		t.Error(err)
		return
	}
	t.Log("after SavE()")

	state.SyncedAt = time.Now().UTC()
	err = state.Save()
	if err != nil {
		t.Error(err)
	}

	if state.IsStale() {
		t.Error("state should be fresh")
	}

	state.SyncedAt = time.Now().UTC().Add(-10 * time.Minute)
	err = state.Save()
	if err != nil {
		t.Error(err)
	}

	if !state.IsStale() {
		t.Error("state should be stale")
	}
}

func TestDeleteMeetingStateByInternalID(t *testing.T) {
	pool := connectTest(t)
	state, err := meetingStateFactory(pool, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if err := state.Save(); err != nil {
		t.Error(err)
	}

	// Now delete the meeting state
	if err := DeleteMeetingStateByInternalID(pool, state.InternalID); err != nil {
		t.Error(err)
	}
}

func TestDeleteMeetingStateByID(t *testing.T) {
	pool := connectTest(t)
	state, err := meetingStateFactory(pool, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if err := state.Save(); err != nil {
		t.Error(err)
	}

	// Now delete the meeting state
	if err := DeleteMeetingStateByID(pool, state.ID); err != nil {
		t.Error(err)
	}
}
