package translator

import (
	"fmt"
	"pilsner/internal/communication"
	"pilsner/proto/pb"
)

func mapItemToProto(pbFunc func(*pb.Item) error) func(*communication.Item) error {
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

func mapConsumerSetupProtoToInternal(pbSetup *pb.ConsumerSetup) communication.ConsumerSetup {
	return communication.ConsumerSetup{
		ReplayMode:          pbSetup.ReplayMode,
		ConsumerName:        pbSetup.ConsumerName,
		RetryPolicy:         pbSetup.RetryPolicy.String(),
		TimeoutMilliSeconds: 0,
	}
}

func mapConsumerAckProtoToInternal(pbAck *pb.ConsumerAck) communication.ConsumerAck {
	return communication.ConsumerAck{
		Status: pbAck.Status.String(),
	}
}

func mapItem(pbItem *pb.Item) communication.Item {
	return communication.Item{
		Content: pbItem.Content,
	}
}

func mapCommunicationItemToProto(communicationItem *communication.Item) pb.Item {

	switch communicationItem.Content.(type) {
	case []byte:
		return pb.Item{
			Content: communicationItem.Content.([]byte),
		}
	case string:
		return pb.Item{
			Content: []byte(communicationItem.Content.(string)),
		}
	default:
		return pb.Item{}
	}

}

func Translate[Out interface{}](input interface{}) (error, Out) {

	var out Out

	switch input.(type) {
	case *pb.ConsumerSetup:
		switch any(out).(type) {
		case communication.ConsumerSetup:
			mappedInput := input.(*pb.ConsumerSetup)
			result := mapConsumerSetupProtoToInternal(mappedInput)
			return nil, any(result).(Out)
		default:
			return fmt.Errorf("no transformation function "), out
		}

	case *pb.ConsumerAck:
		switch any(out).(type) {
		case *communication.ConsumerAck:
			mappedInput := input.(*pb.ConsumerAck)
			result := mapConsumerAckProtoToInternal(mappedInput)
			return nil, any(result).(Out)
		default:
			return fmt.Errorf("no transformation function "), out
		}

	case *pb.Item:
		switch any(out).(type) {
		case communication.Item:
			mappedInput := input.(*pb.Item)
			result := mapItem(mappedInput)
			return nil, any(result).(Out)
		default:
			return fmt.Errorf("no transformation function "), out
		}

	case func(*pb.Item) error:
		switch any(out).(type) {
		case func(*communication.Item) error:
			mappedInput := input.(func(*pb.Item) error)
			result := mapItemToProto(mappedInput)
			return nil, any(result).(Out)
		default:
			return fmt.Errorf("no transformation function "), out
		}

	case *communication.Item:
		switch any(out).(type) {
		case pb.Item:
			mappedInput := input.(*communication.Item)
			result := mapCommunicationItemToProto(mappedInput)
			return nil, any(result).(Out)
		default:
			return fmt.Errorf("no transformation function "), out
		}

	default:
		return fmt.Errorf("no transformation function"), out

	}
}
