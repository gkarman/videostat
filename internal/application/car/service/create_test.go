package service

import (
	"context"
	"errors"
	"testing"

	"github.com/gkarman/demo/internal/infrastructure/dispatcher"
	"github.com/stretchr/testify/require"

	"github.com/gkarman/demo/internal/application/car/requestdto"
	"github.com/gkarman/demo/internal/domain/car"
	carrepo "github.com/gkarman/demo/internal/infrastructure/repository/car"
)

func TestCreateService_Execute_WithExpectedEvents(t *testing.T) {
	tests := []struct {
		name          string
		input         *requestdto.CreateCar
		wantErr       error
		wantSaved     bool
		expectedEvent any
	}{
		{
			name:          "empty name",
			input:         &requestdto.CreateCar{Name: ""},
			wantErr:       car.ErrEmptyName,
			wantSaved:     false,
			expectedEvent: nil,
		},
		{
			name:          "valid name",
			input:         &requestdto.CreateCar{Name: "Tesla"},
			wantErr:       nil,
			wantSaved:     true,
			expectedEvent: car.Created{Name: "Tesla"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := carrepo.NewInMemoryRepo()
			disp := dispatcher.NewFakeDispatcher()
			svc := NewCreate(repo, disp)

			resp, err := svc.Execute(ctx, tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.wantErr))
				require.Nil(t, resp)
				require.Empty(t, disp.Events)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			if tt.wantSaved {
				saved, err := repo.GetByID(ctx, resp.ID)
				require.NoError(t, err)
				require.Equal(t, tt.input.Name, saved.Name)
			}

			if tt.expectedEvent != nil {
				require.Len(t, disp.Events, 1)
				actualEvents := disp.Events[0]
				require.Len(t, actualEvents, 1)

				switch e := tt.expectedEvent.(type) {
				case car.Created:
					act := actualEvents[0].(*car.Created)
					require.Equal(t, e.Name, act.Name)
					require.NotEmpty(t, act.ID)
				default:
					t.Fatalf("unexpected expectedEvent type: %T", e)
				}
			} else {
				require.Empty(t, disp.Events)
			}
		})
	}
}
