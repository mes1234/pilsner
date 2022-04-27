package consumer

import "pilsner/internal/stream"

type Callback func(item stream.Item) error

type RetryPolicy func(callback Callback, item stream.Item) error
