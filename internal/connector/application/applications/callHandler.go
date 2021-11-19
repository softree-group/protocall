package applications

import "github.com/CyCoreSystems/ari/v5"

type CallHandler interface {
	Handle(channel *ari.ChannelHandle)
}
