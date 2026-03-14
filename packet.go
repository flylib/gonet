package gonet

import (
	"encoding/binary"
	"errors"
)

// Packet wire format:
//
// +-------------------+--------------+------------------+
// | body length (2B)  | msg ID (4B)  | body (N bytes)   |
// +-------------------+--------------+------------------+
//
// The "body length" field stores the combined size of msgID(4) + body(N).

const (
	// PktSizeOffset is the byte length of the size field (uint16).
	PktSizeOffset = 2
	// MsgIDOffset is the byte length of the message ID field (uint32).
	MsgIDOffset  = 4
	HeaderOffset = PktSizeOffset + MsgIDOffset // 6 bytes total header
	MTU          = 1500
)

var (
	ErrMessageTooLarge  = errors.New("gonet: message body too large")
	ErrIncompleteHeader = errors.New("gonet: incomplete header")
	ErrIncompletePacket = errors.New("gonet: incomplete packet")
)

// PacketSize returns the total packet size for a given body length.
func PacketSize(bodyLen int) (int, error) {
	bodyFieldLen := MsgIDOffset + bodyLen
	if bodyFieldLen > 0xFFFF {
		return 0, ErrMessageTooLarge
	}
	return HeaderOffset + bodyLen, nil
}

// FillHeader writes the packet header for the given msgID and body length into dst.
// dst must be at least HeaderOffset bytes long.
func FillHeader(dst []byte, msgID uint32, bodyLen int) error {
	if len(dst) < HeaderOffset {
		return errors.New("gonet: header buffer too small")
	}
	bodyFieldLen := MsgIDOffset + bodyLen
	if bodyFieldLen > 0xFFFF {
		return ErrMessageTooLarge
	}
	binary.LittleEndian.PutUint16(dst[:PktSizeOffset], uint16(bodyFieldLen))
	binary.LittleEndian.PutUint32(dst[PktSizeOffset:], msgID)
	return nil
}

// INetPackager encodes and decodes network packets.
type INetPackager interface {
	// Package serializes msgID + v into a wire packet.
	Package(s ISession, msgID uint32, v any) ([]byte, error)
	// UnPackage decodes one message from data.
	// Returns (message, unused_byte_count, error).
	// unused_byte_count > 0 means data contained trailing bytes (multi-packet read).
	UnPackage(s ISession, data []byte) (IMessage, int, error)
}

// DefaultNetPackager is the built-in packet codec.
type DefaultNetPackager struct{}

func (d *DefaultNetPackager) Package(s ISession, msgID uint32, v any) ([]byte, error) {
	body, err := s.GetContext().Marshal(v)
	if err != nil {
		return nil, err
	}
	pktSize, err := PacketSize(len(body))
	if err != nil {
		return nil, err
	}
	pkt := make([]byte, pktSize)
	if err := FillHeader(pkt[:HeaderOffset], msgID, len(body)); err != nil {
		return nil, err
	}
	copy(pkt[HeaderOffset:], body)
	return pkt, nil
}

func (d *DefaultNetPackager) UnPackage(s ISession, data []byte) (IMessage, int, error) {
	if len(data) < HeaderOffset {
		return nil, 0, ErrIncompleteHeader
	}
	bodyFieldLen := int(binary.LittleEndian.Uint16(data[:PktSizeOffset]))
	totalLen := PktSizeOffset + bodyFieldLen
	if len(data) < totalLen {
		return nil, 0, ErrIncompletePacket
	}
	msgID := binary.LittleEndian.Uint32(data[PktSizeOffset : PktSizeOffset+MsgIDOffset])
	// Copy body so the returned message does not retain a reference to the
	// caller's read buffer, which may be reused on the next Read call.
	bodyLen := totalLen - HeaderOffset
	var body []byte
	if bp, ok := s.GetContext().(BodyBufferProvider); ok {
		body = bp.GetBodyBuffer(bodyLen)
	} else {
		body = make([]byte, bodyLen)
	}
	copy(body, data[HeaderOffset:totalLen])
	unused := len(data) - totalLen
	return s.GetContext().NewMsg(msgID, body, s), unused, nil
}
