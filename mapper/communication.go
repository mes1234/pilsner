package mapper

import (
	"fmt"
	"pilsner/internal/communication"
	"pilsner/proto/pb"
)

func MapItemToProto(pbFunc func(*pb.Item) error) func(*communication.Item) error {
	return func(item *communication.Item) error {

		p, ok := item.Content.([]byte)

		if ok != true {
			return fmt.Errorf("wrong type")
		}

		pbItem := pb.Item{
			Content: p,
		}
		return pbFunc(&pbItem)
	}
}

func MapConsumerSetupProtoToInternal(pbSetup *pb.ConsumerSetup) communication.ConsumerSetup {
	return communication.ConsumerSetup{
		ReplayMode:          pbSetup.ReplayMode,
		StreamName:          pbSetup.StreamName,
		ConsumerName:        pbSetup.ConsumerName,
		RetryPolicy:         pbSetup.RetryPolicy.String(),
		TimeoutMilliSeconds: 0,
	}
}
