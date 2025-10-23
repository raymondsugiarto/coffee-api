package entity

import (
	"time"

	b "github.com/getbrevo/brevo-go/lib"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type NotificationInputDto struct {
	TemplateID int64
	From       *b.SendSmtpEmailSender
	To         []b.SendSmtpEmailTo
	Data       map[string]interface{}
}

type NotificationDto struct {
	ID             string
	OrganizationID string
	UserID         string
	User           *UserDto
	RefModule      string
	RefTable       string
	RefID          string
	RefCode        string
	Description    string
	NotifyAt       time.Time
	ReadAt         time.Time
}

func (n *NotificationDto) ToModel() *model.Notification {
	return &model.Notification{
		CommonWithIDs: concern.CommonWithIDs{
			ID: n.ID,
		},
		OrganizationID: n.OrganizationID,
		UserID:         n.UserID,
		RefModule:      n.RefModule,
		RefTable:       n.RefTable,
		RefID:          n.RefID,
		RefCode:        n.RefCode,
		Description:    n.Description,
		NotifyAt:       n.NotifyAt,
		ReadAt:         n.ReadAt,
	}
}

func (n *NotificationDto) FromModel(m *model.Notification) *NotificationDto {
	dto := &NotificationDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		UserID:         m.UserID,
		RefModule:      m.RefModule,
		RefTable:       m.RefTable,
		RefID:          m.RefID,
		RefCode:        m.RefCode,
		Description:    m.Description,
		NotifyAt:       m.NotifyAt,
		ReadAt:         m.ReadAt,
	}
	if dto.User != nil {
		dto.User = new(UserDto).FromModel(m.User)
	}

	return dto
}

type NotificationFindAllRequest struct {
	FindAllRequest
	UserID string
}

func (r *NotificationFindAllRequest) GenerateFilter() {
	if r.UserID != "" {
		r.FindAllRequest.AddFilter(pagination.FilterItem{
			Field: "user_id",
			Op:    "eq",
			Val:   r.UserID,
		})
	}
}
