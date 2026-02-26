package models

import (
	"encoding/json"
)

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
	ID                 json.RawMessage `json:"id"`
	Name               string          `json:"name"`
	Content            string          `json:"content"`
	Status             json.RawMessage `json:"status"`
	StatusName         string          `json:"status_name"`
	Type               json.RawMessage `json:"type"`
	Priority           json.RawMessage `json:"priority"`
	Urgency            json.RawMessage `json:"urgency"`
	Impact             json.RawMessage `json:"impact"`
	EntitiesID         json.RawMessage `json:"entities_id"`
	DateMod            string          `json:"date_mod"`
	UsersIDLastUpdater json.RawMessage `json:"users_id_lastupdater"`
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
	ID             json.RawMessage `json:"id"`
	ItemType       string          `json:"itemtype"`
	ItemsID        json.RawMessage `json:"items_id"`
	Content        string          `json:"content"`
	IsPrivate      json.RawMessage `json:"is_private"`
	Date           string          `json:"date"`
	UsersID        json.RawMessage `json:"users_id"`
	RequestTypesID json.RawMessage `json:"requesttypes_id"`
}

type ParentItemRef struct {
	ID         json.RawMessage `json:"id"`
	Name       string          `json:"name"`
	Status     json.RawMessage `json:"status"`
	Type       json.RawMessage `json:"type"`
	EntitiesID json.RawMessage `json:"entities_id"`
}
