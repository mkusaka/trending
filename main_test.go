package main

import (
    "testing"
    "time"
    "os"
)

func TestPruneRecents_RemovesOlderThanCutoffAndKeepsOrder(t *testing.T) {
    base := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
    cutoff := base // entries older than this should be pruned; equal kept

    rs := []Recent{
        {Logged: base.AddDate(0, 0, -2), Uri: "a", Period: "daily", Language: "go"},  // prune
        {Logged: base.Add(-time.Hour), Uri: "b", Period: "daily", Language: "go"},      // prune
        {Logged: base, Uri: "c", Period: "daily", Language: "go"},                     // keep (equal)
        {Logged: base.Add(time.Hour), Uri: "d", Period: "daily", Language: "go"},      // keep
        {Logged: base.AddDate(0, 0, 1), Uri: "e", Period: "daily", Language: "go"},    // keep
    }

    got := pruneRecents(rs, cutoff)

    wantURIs := []string{"c", "d", "e"}
    if len(got) != len(wantURIs) {
        t.Fatalf("unexpected length: got %d want %d", len(got), len(wantURIs))
    }
    for i, uri := range wantURIs {
        if got[i].Uri != uri {
            t.Fatalf("order mismatch at %d: got %s want %s", i, got[i].Uri, uri)
        }
    }
}

func TestHas_MatchesOnUriPeriodLanguage(t *testing.T) {
    rs := []Recent{
        {Uri: "u1", Period: "daily", Language: "go"},
        {Uri: "u2", Period: "weekly", Language: "go"},
    }

    if !has(rs, Recent{Uri: "u1", Period: "daily", Language: "go"}) {
        t.Fatalf("expected match for u1/daily/go")
    }
    if has(rs, Recent{Uri: "u1", Period: "weekly", Language: "go"}) {
        t.Fatalf("did not expect match for u1/weekly/go")
    }
    if has(rs, Recent{Uri: "u3", Period: "daily", Language: "go"}) {
        t.Fatalf("did not expect match for u3/daily/go")
    }
}

func TestGetRetentionMonthsFromEnv_DefaultAndOverrides(t *testing.T) {
    prev, had := os.LookupEnv("RETENTION_MONTHS")
    if had {
        t.Cleanup(func() { _ = os.Setenv("RETENTION_MONTHS", prev) })
    } else {
        t.Cleanup(func() { _ = os.Unsetenv("RETENTION_MONTHS") })
    }

    // default
    _ = os.Unsetenv("RETENTION_MONTHS")
    if got := getRetentionMonthsFromEnv(); got != 3 {
        t.Fatalf("default retention months: got %d want %d", got, 3)
    }

    // valid override
    _ = os.Setenv("RETENTION_MONTHS", "6")
    if got := getRetentionMonthsFromEnv(); got != 6 {
        t.Fatalf("override retention months: got %d want %d", got, 6)
    }

    // invalid values fallback to default
    _ = os.Setenv("RETENTION_MONTHS", "0")
    if got := getRetentionMonthsFromEnv(); got != 3 {
        t.Fatalf("zero should fallback: got %d want %d", got, 3)
    }
    _ = os.Setenv("RETENTION_MONTHS", "-1")
    if got := getRetentionMonthsFromEnv(); got != 3 {
        t.Fatalf("negative should fallback: got %d want %d", got, 3)
    }
    _ = os.Setenv("RETENTION_MONTHS", "abc")
    if got := getRetentionMonthsFromEnv(); got != 3 {
        t.Fatalf("non-int should fallback: got %d want %d", got, 3)
    }
}
