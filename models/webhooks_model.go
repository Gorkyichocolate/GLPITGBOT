package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type StringOrNumber string

func (v *StringOrNumber) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		*v = ""
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		*v = StringOrNumber(strings.TrimSpace(asString))
		return nil
	}

	var asNumber json.Number
	if err := json.Unmarshal(data, &asNumber); err == nil {
		*v = StringOrNumber(strings.TrimSpace(asNumber.String()))
		return nil
	}

	var asBool bool
	if err := json.Unmarshal(data, &asBool); err == nil {
		*v = StringOrNumber(strconv.FormatBool(asBool))
		return nil
	}

	return fmt.Errorf("unsupported value type: %s", string(data))
}

func (v StringOrNumber) String() string {
	return strings.TrimSpace(string(v))
}

type StatusValue struct {
	ID    StringOrNumber
	Name  string
	Value StringOrNumber
}

func (s *StatusValue) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		s.ID = ""
		s.Name = ""
		s.Value = ""
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		s.Value = StringOrNumber(strings.TrimSpace(asString))
		s.ID = ""
		s.Name = ""
		return nil
	}

	var asNumber json.Number
	if err := json.Unmarshal(data, &asNumber); err == nil {
		s.Value = StringOrNumber(strings.TrimSpace(asNumber.String()))
		s.ID = ""
		s.Name = ""
		return nil
	}

	var asObject struct {
		ID   StringOrNumber `json:"id"`
		Name string         `json:"name"`
	}
	if err := json.Unmarshal(data, &asObject); err == nil {
		s.ID = asObject.ID
		s.Name = strings.TrimSpace(asObject.Name)
		s.Value = ""
		return nil
	}

	return fmt.Errorf("unsupported status value type: %s", string(data))
}

func (s StatusValue) String() string {
	if strings.TrimSpace(s.Name) != "" {
		return strings.TrimSpace(s.Name)
	}

	if s.ID.String() != "" {
		return s.ID.String()
	}

	return s.Value.String()
}

type BaseWebhookPayload struct {
	Event string          `json:"event"`
	Item  json.RawMessage `json:"item"`
}

type WebhookItemTypeProbe struct {
	ItemType string `json:"itemtype"`
}

type UpdateEvent struct {
	Event string     `json:"event"`
	Item  UpdateItem `json:"item"`
}

type UpdateItem struct {
	ID                 StringOrNumber `json:"id"`
	Name               string         `json:"name"`
	Content            string         `json:"content"`
	Status             StatusValue    `json:"status"`
	StatusName         string         `json:"status_name"`
	Type               StringOrNumber `json:"type"`
	Priority           StringOrNumber `json:"priority"`
	Urgency            StringOrNumber `json:"urgency"`
	Impact             StringOrNumber `json:"impact"`
	EntitiesID         StringOrNumber `json:"entities_id"`
	DateMod            string         `json:"date_mod"`
	UsersIDLastUpdater StringOrNumber `json:"users_id_lastupdater"`
}

type AddFollowupEvent struct {
	Event      string        `json:"event"`
	Item       FollowupItem  `json:"item"`
	ParentItem ParentItemRef `json:"parent_item"`
}

type FullFollowupWebhookPayload struct {
	Event      string        `json:"event"`
	Item       FollowupItem  `json:"item"`
	ParentItem ParentItemRef `json:"parent_item,omitempty"`
}

type FullTicketWebhookPayload struct {
	Event string     `json:"event"`
	Item  UpdateItem `json:"item"`
}

type FollowupItem struct {
	ID             StringOrNumber `json:"id"`
	ItemType       string         `json:"itemtype"`
	ItemsID        StringOrNumber `json:"items_id"`
	Content        string         `json:"content"`
	IsPrivate      StringOrNumber `json:"is_private"`
	Date           string         `json:"date"`
	UsersID        StringOrNumber `json:"users_id"`
	RequestTypesID StringOrNumber `json:"requesttypes_id"`
}

type ParentItemRef struct {
	ID         StringOrNumber `json:"id"`
	Name       string         `json:"name"`
	Status     StatusValue    `json:"status"`
	Type       StringOrNumber `json:"type"`
	EntitiesID StringOrNumber `json:"entities_id"`
}

func normalizeTicketStatusRU(status string) string {
	s := strings.TrimSpace(strings.ToLower(status))

	switch s {
	case "1", "new", "новая", "новый":
		return "Новая"
	case "2", "in_progress", "processing", "в работе":
		return "В работе"
	case "3", "pending", "waiting", "ожидание", "в ожидании":
		return "Ожидание"
	case "4", "solved", "resolved", "решена", "решено":
		return "Решена"
	case "5", "closed", "закрыто", "закрыта":
		return "Закрыто"
	case "6", "cancelled", "canceled", "отменена", "отменено":
		return "Отменена"
	default:
		if status == "" {
			return "Неизвестно"
		}

		return status
	}
}

func BuildTicketCommentNotification(ticketName, commentText, status string) string {
	ticketName = strings.TrimSpace(ticketName)
	if ticketName == "" {
		ticketName = "Без названия"
	}

	commentText = strings.TrimSpace(commentText)
	if commentText == "" {
		commentText = "(пустой комментарий)"
	}

	return fmt.Sprintf(
		"💬 Комментарий к заявке\nЗаявка: %s\nКомментарий: %s\nСтатус: %s",
		ticketName,
		commentText,
		normalizeTicketStatusRU(status),
	)
}

func BuildTicketStatusChangedNotification(ticketName, newStatus string) string {
	ticketName = strings.TrimSpace(ticketName)
	if ticketName == "" {
		ticketName = "Без названия"
	}

	return fmt.Sprintf(
		"🔄 Статус заявки изменён\nЗаявка: %s\nНовый статус: %s",
		ticketName,
		normalizeTicketStatusRU(newStatus),
	)
}
