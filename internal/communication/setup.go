package communication

type ConsumerSetup struct {
	ReplayMode          bool
	StreamName          string
	ConsumerName        string
	RetryPolicy         string
	TimeoutMilliSeconds int32
}
type RetryPolicy string
