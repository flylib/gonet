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
	PktSizeOffset = 4                           // uint16: stores len(msgID field + body)
	MsgIDOffset   = 4                           // uint32: message ID
	HeaderOffset  = PktSizeOffset + MsgIDOffset // 6 bytes total header
	MTU           = 1500
)

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
	bodyFieldLen := MsgIDOffset + len(body) // what PktSizeOffset field stores
	if bodyFieldLen > 0xFFFF {
		return nil, errors.New("gonet: message body too large")
	}
	pkt := make([]byte, HeaderOffset+len(body))
	binary.LittleEndian.PutUint16(pkt, uint16(bodyFieldLen))
	binary.LittleEndian.PutUint32(pkt[PktSizeOffset:], msgID)
	copy(pkt[HeaderOffset:], body)
	return pkt, nil
}

func (d *DefaultNetPackager) UnPackage(s ISession, data []byte) (IMessage, int, error) {
	if len(data) < HeaderOffset {
		return nil, 0, errors.New("gonet: incomplete header")
	}
	bodyFieldLen := int(binary.LittleEndian.Uint16(data[:PktSizeOffset]))
	totalLen := PktSizeOffset + bodyFieldLen
	if len(data) < totalLen {
		return nil, 0, errors.New("gonet: incomplete packet")
	}
	msgID := binary.LittleEndian.Uint32(data[PktSizeOffset : PktSizeOffset+MsgIDOffset])
	// Copy body so the returned message does not retain a reference to the
	// caller's read buffer, which may be reused on the next Read call.
	body := make([]byte, totalLen-HeaderOffset)
	copy(body, data[HeaderOffset:totalLen])
	unused := len(data) - totalLen
	return s.GetContext().NewMsg(msgID, body, s), unused, nil
}
