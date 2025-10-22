package nuxeo

import (
	"time"
)

type AuditLogEntry struct {
	entity
	ID            int            `json:"id"`
	Category      string         `json:"category"`
	PrincipalName string         `json:"principalName"`
	Comment       string         `json:"comment"`
	DocLifeCycle  string         `json:"docLifeCycle"`
	DocPath       string         `json:"docPath"`
	DocType       string         `json:"docType"`
	DocUUID       string         `json:"docUUID"`
	EventID       string         `json:"eventId"`
	RepositoryID  string         `json:"repositoryId"`
	EventDate     time.Time      `json:"eventDate"`
	LogDate       time.Time      `json:"logDate"`
	Extended      map[string]any `json:"extended"`
}

type Audit paginableEntities[AuditLogEntry]
