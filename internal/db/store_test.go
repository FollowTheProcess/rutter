package db_test

import (
	"path/filepath"
	"slices"
	"testing"
	"time"

	"go.followtheprocess.codes/rutter/internal/db"
	"go.followtheprocess.codes/test"
)

func TestStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	store, err := db.Open(t.Context(), path)
	test.Ok(t, err)

	if store == nil {
		t.Fatal("store returned from db.Open was nil")
	}

	t.Cleanup(func() { store.Close() })

	now := time.Now().UTC()
	session := "c59e0707-7eba-4110-88ec-910ccd4f541e"

	id, err := store.Query.StartCommand(t.Context(), db.StartCommandParams{
		Cmd:       "echo 'hello world'",
		Cwd:       "some/fake/dir",
		Session:   session,
		StartedAt: now,
	})

	test.Ok(t, err)
	test.Equal(t, id, 1)

	// List it back again
	got, err := store.Query.ListCandidates(t.Context(), 1)
	test.Ok(t, err)

	want := []db.History{
		{
			ID:        id,
			Cmd:       "echo 'hello world'",
			Cwd:       "some/fake/dir",
			StartedAt: now,
			Duration:  -1, // Never finished so Duration is -1
			Exit:      -1, // Same with exit
			Session:   session,
		},
	}

	test.EqualFunc(t, got, want, slices.Equal)

	// Test the FinishCommand
	err = store.Query.FinishCommand(t.Context(), db.FinishCommandParams{
		Duration: 2 * time.Second,
		Exit:     1,
		ID:       id,
	})

	test.Ok(t, err)

	// If we get it again, we should see the duration and exit code.
	got, err = store.Query.ListCandidates(t.Context(), 1)
	test.Ok(t, err)

	want = []db.History{
		{
			ID:        id,
			Cmd:       "echo 'hello world'",
			Cwd:       "some/fake/dir",
			StartedAt: now,
			Duration:  2 * time.Second,
			Exit:      1,
			Session:   session,
		},
	}

	test.EqualFunc(t, got, want, slices.Equal)
}
