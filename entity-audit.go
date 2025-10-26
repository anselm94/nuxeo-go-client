package nuxeo

type entityAuditLogEntry struct {
	entity
	ID            int              `json:"id"`
	Category      string           `json:"category"`
	PrincipalName string           `json:"principalName"`
	Comment       string           `json:"comment"`
	DocLifeCycle  string           `json:"docLifeCycle"`
	DocPath       string           `json:"docPath"`
	DocType       string           `json:"docType"`
	DocUUID       string           `json:"docUUID"`
	EventID       string           `json:"eventId"`
	RepositoryID  string           `json:"repositoryId"`
	EventDate     *ISO8601Time     `json:"eventDate"`
	LogDate       *ISO8601Time     `json:"logDate"`
	Extended      map[string]Field `json:"extended"`
}

type entityAudit paginableEntities[entityAuditLogEntry]
