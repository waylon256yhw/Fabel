// Package resource provides the unified resource index layer.
// Domain CRUD owns truth; the resource index is a projection for search and relations.
package resource

import (
	"context"
	"fmt"
	"strings"

	"fabel/internal/dbq"
)

// Kind enumerates the resource types that can be indexed.
type Kind string

const (
	KindCharacter    Kind = "character"
	KindPreset       Kind = "preset"
	KindConversation Kind = "conversation"
	KindLorebook     Kind = "lorebook"
	KindLorebookEntry Kind = "lorebook_entry"
	KindPersona      Kind = "persona"
)

// Relation types between resources.
const (
	RelContains   = "contains"    // parent contains child (lorebook → entry)
	RelAttachedTo = "attached_to" // entity is attached to another (conversation → character)
	RelUses       = "uses"        // entity uses another (conversation → preset)
)

// Document represents a resource to be indexed.
type Document struct {
	ID      string
	UserID  string
	Kind    Kind
	Title   string
	Content string // full text for search (description, prompt text, etc.)
	Tags    []string
	Meta    string // JSON string, default "{}"
}

// Indexer writes resource projections into the unified index.
// Domain CRUD code calls these after successful writes.
type Indexer interface {
	Upsert(ctx context.Context, doc Document) error
	Delete(ctx context.Context, resourceID string) error
}

// SQLiteIndexer implements Indexer using sqlc-generated queries.
type SQLiteIndexer struct {
	q *dbq.Queries
}

// NewIndexer creates a new SQLite-backed resource indexer.
func NewIndexer(q *dbq.Queries) *SQLiteIndexer {
	return &SQLiteIndexer{q: q}
}

func (idx *SQLiteIndexer) Upsert(ctx context.Context, doc Document) error {
	meta := doc.Meta
	if meta == "" {
		meta = "{}"
	}
	return idx.q.UpsertResource(ctx, dbq.UpsertResourceParams{
		ID:      doc.ID,
		UserID:  doc.UserID,
		Kind:    string(doc.Kind),
		Title:   doc.Title,
		Content: doc.Content,
		Tags:    strings.Join(doc.Tags, " "),
		Meta:    meta,
	})
}

func (idx *SQLiteIndexer) Delete(ctx context.Context, resourceID string) error {
	// Clean up relations in both directions, then the resource itself.
	if err := idx.q.DeleteRelationsBySrc(ctx, resourceID); err != nil {
		return fmt.Errorf("delete relations (src): %w", err)
	}
	if err := idx.q.DeleteRelationsByDst(ctx, resourceID); err != nil {
		return fmt.Errorf("delete relations (dst): %w", err)
	}
	return idx.q.DeleteResource(ctx, resourceID)
}

// Relate creates a directed relation between two resources.
func (idx *SQLiteIndexer) Relate(ctx context.Context, srcID, dstID, relation string) error {
	return idx.q.UpsertRelation(ctx, dbq.UpsertRelationParams{
		SrcID:    srcID,
		DstID:    dstID,
		Relation: relation,
	})
}
