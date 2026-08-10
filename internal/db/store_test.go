package db_test

import (
	"path/filepath"
	"slices"
	"testing"
	"time"

	"go.followtheprocess.codes/rutter/internal/db"
	"go.followtheprocess.codes/test"
)

func TestStartFinishCommand(t *testing.T) {
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

func TestListByDirectory(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(1 * time.Second)
	evenLater := later.Add(2 * time.Second)
	session := "8d921a03-d07a-48da-8cef-beb033ba9268"

	tests := []struct {
		name  string       // Name of the test case
		dir   string       // Directory to filter candidates by
		seed  []db.History // Records to seed the DB with
		want  []db.History // Expected returned entries
		limit int          // Maximum number of matches to return
	}{
		{
			name:  "empty",
			seed:  []db.History{},
			dir:   "",
			limit: 1,
			want:  []db.History{},
		},
		{
			name: "all same dir",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        3,
					Cmd:       "false",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: evenLater,
					Duration:  18 * time.Millisecond,
					Exit:      1,
				},
			},
			dir:   "/some/dir",
			limit: 3,
			want: []db.History{
				// Reverse order, latest first
				{
					ID:        3,
					Cmd:       "false",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: evenLater,
					Duration:  18 * time.Millisecond,
					Exit:      1,
				},
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
			},
		},
		{
			name: "filters by dir",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/dir/we/want",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/dir/we/want",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        3,
					Cmd:       "false",
					Cwd:       "/somewhere/else",
					Session:   session,
					StartedAt: evenLater,
					Duration:  18 * time.Millisecond,
					Exit:      1,
				},
			},
			dir:   "/dir/we/want",
			limit: 3,
			want: []db.History{
				// Reverse order, latest first
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/dir/we/want",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/dir/we/want",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
			},
		},
		{
			name: "keeps only the latest of a repeated command",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        2,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  412 * time.Millisecond,
					Exit:      0,
				},
			},
			dir:   "/some/dir",
			limit: 3,
			want: []db.History{
				{
					ID:        2,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  412 * time.Millisecond,
					Exit:      0,
				},
			},
		},
		{
			name: "limit caps the results",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
			},
			dir:   "/some/dir",
			limit: 1,
			want: []db.History{
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "test.db")

			store, err := db.Open(t.Context(), path)
			test.Ok(t, err)

			if store == nil {
				t.Fatal("store returned from db.Open was nil")
			}

			t.Cleanup(func() { store.Close() })

			for _, entry := range tt.seed {
				var id int64

				id, err = store.Query.StartCommand(t.Context(), db.StartCommandParams{
					Cmd:       entry.Cmd,
					Cwd:       entry.Cwd,
					Session:   entry.Session,
					StartedAt: entry.StartedAt,
				})

				test.Ok(t, err)

				err = store.Query.FinishCommand(t.Context(), db.FinishCommandParams{
					Duration: entry.Duration,
					Exit:     entry.Exit,
					ID:       id,
				})

				test.Ok(t, err)
			}

			got, err := store.Query.ListCandidatesInDirectory(t.Context(), db.ListCandidatesInDirectoryParams{
				Cwd:   tt.dir,
				Limit: int64(tt.limit),
			})

			test.Ok(t, err)

			test.EqualFunc(t, got, tt.want, slices.Equal)
		})
	}
}

func TestListBySession(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(1 * time.Second)
	evenLater := later.Add(2 * time.Second)
	session := "2f0c2b4f-2c3a-4a2e-9d1b-6a0f4c8e7b11"
	other := "b7d1e5c8-9a3f-4c6e-8f2d-1e4a7b0c3d59"

	tests := []struct {
		name    string       // Name of the test case
		session string       // Session to filter candidates by
		seed    []db.History // Records to seed the DB with
		want    []db.History // Expected returned entries
		limit   int          // Maximum number of matches to return
	}{
		{
			name:    "empty",
			seed:    []db.History{},
			session: "",
			limit:   1,
			want:    []db.History{},
		},
		{
			name: "all same session",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/another/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        3,
					Cmd:       "false",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: evenLater,
					Duration:  18 * time.Millisecond,
					Exit:      1,
				},
			},
			session: session,
			limit:   3,
			want: []db.History{
				// Reverse order, latest first
				{
					ID:        3,
					Cmd:       "false",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: evenLater,
					Duration:  18 * time.Millisecond,
					Exit:      1,
				},
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/another/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
			},
		},
		{
			name: "filters by session",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        3,
					Cmd:       "false",
					Cwd:       "/some/dir",
					Session:   other,
					StartedAt: evenLater,
					Duration:  18 * time.Millisecond,
					Exit:      1,
				},
			},
			session: session,
			limit:   3,
			want: []db.History{
				// Reverse order, latest first
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
			},
		},
		{
			name: "keeps only the latest of a repeated command",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        2,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  412 * time.Millisecond,
					Exit:      0,
				},
			},
			session: session,
			limit:   3,
			want: []db.History{
				{
					ID:        2,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  412 * time.Millisecond,
					Exit:      0,
				},
			},
		},
		{
			name: "limit caps the results",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
			},
			session: session,
			limit:   1,
			want: []db.History{
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "test.db")

			store, err := db.Open(t.Context(), path)
			test.Ok(t, err)

			if store == nil {
				t.Fatal("store returned from db.Open was nil")
			}

			t.Cleanup(func() { store.Close() })

			for _, entry := range tt.seed {
				var id int64

				id, err = store.Query.StartCommand(t.Context(), db.StartCommandParams{
					Cmd:       entry.Cmd,
					Cwd:       entry.Cwd,
					Session:   entry.Session,
					StartedAt: entry.StartedAt,
				})

				test.Ok(t, err)

				err = store.Query.FinishCommand(t.Context(), db.FinishCommandParams{
					Duration: entry.Duration,
					Exit:     entry.Exit,
					ID:       id,
				})

				test.Ok(t, err)
			}

			got, err := store.Query.ListCandidatesInSession(t.Context(), db.ListCandidatesInSessionParams{
				Session: tt.session,
				Limit:   int64(tt.limit),
			})

			test.Ok(t, err)

			test.EqualFunc(t, got, tt.want, slices.Equal)
		})
	}
}

func TestSuggestByPrefix(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(1 * time.Second)
	session := "2ea72024-4ad1-4663-9fda-bd4d13a5b27d"

	tests := []struct {
		name    string       // Name of the test case
		prefix  string       // Prefix to match
		want    string       // Expected suggestion
		seed    []db.History // Records to seed the test DB with
		wantErr bool         // Whether we want an error
	}{
		{
			name:    "empty",
			prefix:  "echo",
			want:    "",
			seed:    []db.History{},
			wantErr: true,
		},
		{
			name: "no match",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
			},
			prefix:  "git",
			wantErr: true,
		},
		{
			name: "latest match wins",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
				{
					ID:        2,
					Cmd:       "echo 'again'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
			},
			prefix: "echo",
			want:   "echo 'again'",
		},
		{
			name: "must match at the start",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "echo 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
				{
					// Later, and contains the prefix, but not at the start
					ID:        2,
					Cmd:       "sudo echo 'again'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: later,
					Duration:  267 * time.Millisecond,
					Exit:      0,
				},
			},
			prefix: "echo",
			want:   "echo 'hello'",
		},
		{
			name: "case sensitive",
			seed: []db.History{
				{
					ID:        1,
					Cmd:       "ECHO 'hello'",
					Cwd:       "/some/dir",
					Session:   session,
					StartedAt: now,
					Duration:  375 * time.Millisecond,
					Exit:      0,
				},
			},
			prefix:  "echo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "test.db")

			store, err := db.Open(t.Context(), path)
			test.Ok(t, err)

			if store == nil {
				t.Fatal("store returned from db.Open was nil")
			}

			t.Cleanup(func() { store.Close() })

			for _, entry := range tt.seed {
				var id int64

				id, err = store.Query.StartCommand(t.Context(), db.StartCommandParams{
					Cmd:       entry.Cmd,
					Cwd:       entry.Cmd,
					Session:   entry.Session,
					StartedAt: entry.StartedAt,
				})

				test.Ok(t, err)

				err = store.Query.FinishCommand(t.Context(), db.FinishCommandParams{
					Duration: entry.Duration,
					Exit:     entry.Exit,
					ID:       id,
				})

				test.Ok(t, err)
			}

			got, err := store.Query.SuggestByPrefix(t.Context(), tt.prefix)
			test.WantErr(t, err, tt.wantErr)
			test.Equal(t, got, tt.want)
		})
	}
}
