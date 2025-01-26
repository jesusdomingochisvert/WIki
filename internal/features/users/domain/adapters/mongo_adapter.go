package repositories

import (
	"context"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoCursorAdapter struct {
	cursor *mongo.Cursor
}

func (m *mongoCursorAdapter) Next(ctx context.Context) bool {
	return m.cursor.Next(ctx)
}

func (m *mongoCursorAdapter) Decode(v interface{}) error {
	return m.cursor.Decode(v)
}

func (m *mongoCursorAdapter) Close(ctx context.Context) error {
	return m.cursor.Close(ctx)
}

func (m *mongoCursorAdapter) Err() error {
	return m.cursor.Err()
}

type mongoCollectionAdapter struct {
	collection *mongo.Collection
}

func (m *mongoCollectionAdapter) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (interfaces.Cursor, error) {
	cur, err := m.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	return &mongoCursorAdapter{cursor: cur}, nil
}

type mongoDatabaseAdapter struct {
	database *mongo.Database
}

func (m *mongoDatabaseAdapter) Collection(name string) interfaces.Collection {
	return &mongoCollectionAdapter{collection: m.database.Collection(name)}
}

func NewMongoDatabaseAdapter(db *mongo.Database) interfaces.Database {
	return &mongoDatabaseAdapter{database: db}
}
