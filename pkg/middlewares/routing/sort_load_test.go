package routing

import (
	"sort"
	"testing"

	"gitlab.com/infra.run/public/b3scale/pkg/cluster"
	"gitlab.com/infra.run/public/b3scale/pkg/store"
)

func TestSortBackendByLoad(t *testing.T) {

	b := []*cluster.Backend{
		cluster.NewBackend(nil, &store.BackendState{
			ID:             "A",
			MeetingsCount:  20,
			LoadFactor:     1,
			AttendeesCount: 12,
		}),
		cluster.NewBackend(nil, &store.BackendState{
			ID:             "B",
			MeetingsCount:  10,
			LoadFactor:     1,
			AttendeesCount: 12,
		}),
		cluster.NewBackend(nil, &store.BackendState{
			ID:             "C",
			MeetingsCount:  0,
			LoadFactor:     1,
			AttendeesCount: 0,
		}),
	}

	sort.Sort(BackendsByLoad(b))

	if b[0].ID() != "C" {
		t.Error("unexpected:", b[0])
	}
	if b[1].ID() != "B" {
		t.Error("unexpected:", b[1])
	}
	if b[2].ID() != "A" {
		t.Error("unexpected:", b[2])
	}

}