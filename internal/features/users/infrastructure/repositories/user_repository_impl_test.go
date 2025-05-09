package repositories

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/domain/mocks"
	"github.com/jesusdomingochisvert/WIki/internal/features/users/infrastructure/db/schema"
	"testing"
)

func TestGetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockCollection(ctrl)
	mockDatabase := mocks.NewMockDatabase(ctrl)

	mockDatabase.EXPECT().Collection("users").Return(mockCollection).AnyTimes()

	repository := NewUserRepository(mockDatabase)

	t.Run("Success to get all users", func(t *testing.T) {
		mockCursor := mocks.NewMockCursor(ctrl)

		mockCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mockCursor, nil)
		mockCursor.EXPECT().Next(gomock.Any()).Return(true).Times(2)
		mockCursor.EXPECT().Next(gomock.Any()).Return(false)

		userSchema1 := schema.UserSchema{ID: "1", Name: "test1", Email: "test1@test.com", Username: "test1", Password: "test1"}
		userSchema2 := schema.UserSchema{ID: "2", Name: "test2", Email: "test2@test.com", Username: "test2", Password: "test2"}

		mockCursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(userSchema *schema.UserSchema) error {
			*userSchema = userSchema1
			return nil
		}).Times(1)
		mockCursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(userSchema *schema.UserSchema) error {
			*userSchema = userSchema2
			return nil
		}).Times(1)

		mockCursor.EXPECT().Close(gomock.Any()).Return(nil)

		mockCursor.
			EXPECT().
			Err().
			Return(nil)

		users, err := repository.GetAllUsers(context.Background())
		if err != nil {
			t.Fatalf("Se esperaba éxito, pero ocurrió error: %v", err)
		}

		if len(users) != 2 {
			t.Fatalf("Se esperaban 2 usuarios, pero se obtuvieron %d", len(users))
		}

		if users[0].ID != "1" || users[1].ID != "2" {
			t.Errorf("Los usuarios obtenidos no coinciden con los esperados")
		}
	})

	t.Run("Error al ejecutar Find", func(t *testing.T) {
		mockCollection.
			EXPECT().
			Find(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errors.New("error de Find"))

		users, err := repository.GetAllUsers(context.Background())
		if err == nil {
			t.Fatalf("Se esperaba error, pero se obtuvo nil")
		}

		if users != nil {
			t.Fatalf("Se esperaba nil para users, pero se obtuvo: %v", users)
		}
	})

	t.Run("No documents found", func(t *testing.T) {
		mockCursor := mocks.NewMockCursor(ctrl)

		mockCollection.
			EXPECT().
			Find(gomock.Any(), gomock.Any()).
			Return(mockCursor, nil)

		mockCursor.
			EXPECT().
			Next(gomock.Any()).
			Return(false)

		mockCursor.
			EXPECT().
			Err().
			Return(nil)

		mockCursor.
			EXPECT().
			Close(gomock.Any()).
			Return(nil)

		users, err := repository.GetAllUsers(context.Background())
		if err != nil {
			t.Fatalf("No se esperaba error, pero ocurrió: %v", err)
		}

		if len(users) != 0 {
			t.Fatalf("Se esperaban 0 usuarios, pero se obtuvieron %d", len(users))
		}
	})

	t.Run("Error al decodificar usuario", func(t *testing.T) {
		mockCursor := mocks.NewMockCursor(ctrl)

		mockCollection.
			EXPECT().
			Find(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockCursor, nil)

		mockCursor.
			EXPECT().
			Next(gomock.Any()).
			Return(true).
			Times(1)

		mockCursor.
			EXPECT().
			Decode(gomock.Any()).
			Return(errors.New("error de Decode"))

		mockCursor.
			EXPECT().
			Close(gomock.Any()).
			Return(nil)

		users, err := repository.GetAllUsers(context.Background())
		if err == nil {
			t.Fatalf("Se esperaba error de Decode, pero se obtuvo nil")
		}

		if users != nil {
			t.Fatalf("Se esperaba nil para users, pero se obtuvo: %v", users)
		}
	})

	t.Run("Cursor returns error after iteration", func(t *testing.T) {
		mockCursor := mocks.NewMockCursor(ctrl)

		mockCollection.
			EXPECT().
			Find(gomock.Any(), gomock.Any()).
			Return(mockCursor, nil)

		mockCursor.
			EXPECT().
			Next(gomock.Any()).
			Return(false)

		mockCursor.
			EXPECT().
			Err().
			Return(errors.New("mocked cursor error"))

		mockCursor.
			EXPECT().
			Close(gomock.Any()).
			Return(nil)

		users, err := repository.GetAllUsers(context.Background())
		if err == nil {
			t.Fatalf("Se esperaba error en cur.Err(), pero no se obtuvo")
		}
		if users != nil {
			t.Fatalf("Se esperaba nil en users, pero se obtuvo: %v", users)
		}
	})

	t.Run("Contexto cancelado", func(t *testing.T) {
		mockCursor := mocks.NewMockCursor(ctrl)

		mockCollection.
			EXPECT().
			Find(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockCursor, nil)

		mockCursor.
			EXPECT().
			Next(gomock.Any()).
			Return(false)

		mockCursor.
			EXPECT().
			Close(gomock.Any()).
			Return(nil)

		mockCursor.
			EXPECT().
			Err().
			Return(nil)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		users, err := repository.GetAllUsers(ctx)
		if err == nil {
			t.Fatalf("Se esperaba error de contexto cancelado, pero se obtuvo nil")
		}

		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Se esperaba error de contexto cancelado, pero se obtuvo: %v", err)
		}

		if users != nil {
			t.Fatalf("Se esperaba nil para users, pero se obtuvo: %v", users)
		}
	})
}
