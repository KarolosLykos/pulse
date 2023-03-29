package pulse

import (
	"github.com/KarolosLykos/pulse/proto"
	"math"
)

// A Sink is an output device.
type Sink struct {
	info proto.GetSinkInfoReply
}

// ListSinks returns a list of all available output devices.
func (c *Client) ListSinks() ([]*Sink, error) {
	var reply proto.GetSinkInfoListReply
	err := c.c.Request(&proto.GetSinkInfoList{}, &reply)
	if err != nil {
		return nil, err
	}
	sinks := make([]*Sink, len(reply))
	for i := range sinks {
		sinks[i] = &Sink{info: *reply[i]}
	}
	return sinks, nil
}

// DefaultSink returns the default output device.
func (c *Client) DefaultSink() (*Sink, error) {
	var sink Sink
	err := c.c.Request(&proto.GetSinkInfo{SinkIndex: proto.Undefined}, &sink.info)
	if err != nil {
		return nil, err
	}
	return &sink, nil
}

// SetDefaultSink sets the default sink.
func (c *Client) SetDefaultSink(sinkName string) error {
	return c.RawRequest(&proto.SetDefaultSink{SinkName: sinkName}, nil)
}

// SetSinkPort sets port to the sink.
func (c *Client) SetSinkPort(sinkIndex uint32, sinkName, port string) error {
	return c.RawRequest(&proto.SetSinkPort{SinkIndex: sinkIndex, SinkName: sinkName, Port: port}, nil)
}

// SinkMuteToggle toggles mute on a sink.
func (c *Client) SinkMuteToggle(sinkIndex uint32, sinkName string) error {
	sink, err := c.DefaultSink()
	if err != nil {
		return err
	}

	return c.RawRequest(&proto.SetSinkMute{SinkIndex: sinkIndex, SinkName: sinkName, Mute: !sink.info.Mute}, nil)
}

// SetSinkVolume toggles mute on a sink.
func (c *Client) SetSinkVolume(sinkIndex uint32, sinkName string, volume float64) error {
	return c.RawRequest(&proto.SetSinkVolume{
		SinkIndex:      sinkIndex,
		SinkName:       sinkName,
		ChannelVolumes: proto.ChannelVolumes{uint32(math.Round(volume * float64(proto.VolumeNorm) / 100))},
	}, nil)
}

// SinkByID looks up a sink id.
func (c *Client) SinkByID(name string) (*Sink, error) {
	var sink Sink
	err := c.c.Request(&proto.GetSinkInfo{SinkIndex: proto.Undefined, SinkName: name}, &sink.info)
	if err != nil {
		return nil, err
	}
	return &sink, nil
}

// ID returns the sink name. Sink names are unique identifiers, but not necessarily human-readable.
func (s *Sink) ID() string {
	return s.info.SinkName
}

// Name is a human-readable name describing the sink.
func (s *Sink) Name() string {
	return s.info.Device
}

// Channels returns the default channel map.
func (s *Sink) Channels() proto.ChannelMap {
	return s.info.ChannelMap
}

// SampleRate returns the default sample rate.
func (s *Sink) SampleRate() int {
	return int(s.info.Rate)
}

// Volume returns the sink volume.
func (s *Sink) Volume() float64 {
	return float64(s.info.ChannelVolumes[0]) / float64(proto.VolumeNorm) * 100.0
}

// SinkIndex returns the sink index.
// This should only be used together with (*Client).RawRequest.
func (s *Sink) SinkIndex() uint32 {
	return s.info.SinkIndex
}

// Info is a helper method that exposes Sink properties.
func (s *Sink) Info() proto.GetSinkInfoReply {
	return s.info
}
