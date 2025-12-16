package service

import (
    "context"
    "testing"
    "time"

    "industry-api/internal/models"
)

type mockRoomMaintRepo struct{
    add func(ctx context.Context, rm *models.RoomMaintenance) (*models.RoomMaintenance, error)
}
func (m *mockRoomMaintRepo) AddRoomMaintenance(ctx context.Context, rm *models.RoomMaintenance) (*models.RoomMaintenance, error) { return m.add(ctx,rm) }

func TestAddRoomMaintenance_Validation(t *testing.T) {
    svc := &RoomMaintenanceService{repo: &mockRoomMaintRepo{add: func(ctx context.Context, rm *models.RoomMaintenance) (*models.RoomMaintenance, error){ rm.ID=1; return rm,nil }}}
    rm := &models.RoomMaintenance{}
    if _, err := svc.AddRoomMaintenance(context.Background(), rm); err == nil { t.Fatalf("expected validation error") }

    rm = &models.RoomMaintenance{RoomID:1, StartDate:time.Now(), EndDate:time.Now().Add(time.Hour), Reason:"fix", Status:"scheduled", CreatedBy:1}
    got, err := svc.AddRoomMaintenance(context.Background(), rm)
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if got.ID == 0 { t.Fatalf("expected ID set") }
}
