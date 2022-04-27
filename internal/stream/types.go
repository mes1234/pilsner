package stream

// Item is a single portion of data in memoryStream
type Item struct {
	content []byte //Raw content of item
}

// Context is container for all metadata regarding memoryStream
type Context struct {
}
