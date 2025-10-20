package hook

type NoOpHook struct{}

func NewNoOpHook() NoOpHook {
	return NoOpHook{}
}

func (h NoOpHook) BeforeRequest(method, url string) {}

func (h NoOpHook) AfterResponse(method, url string, status int) {}
