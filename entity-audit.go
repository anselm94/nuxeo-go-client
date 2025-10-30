package nuxeo

// AuditLogEntry represents a single audit log entry for a Nuxeo document or event.
// It contains metadata about the event, principal, document, and extended fields.
type AuditLogEntry struct {
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

// Audit is a paginated collection of AuditLogEntry objects.
type Audit paginableEntities[AuditLogEntry]
