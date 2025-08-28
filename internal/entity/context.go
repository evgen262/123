package entity

import (
	"context"
	"errors"
	"net"

	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

//go:generate ditgen -source=context.go

const (
	sessionContextKey  ContextKey = "session"
	deviceContextKey   ContextKey = "device"
	clientIPContextKey ContextKey = "client-ip"
)

type UserAgentContextKey struct{}
type UserAgentContextData map[string]string
type ContextKey string

func (k ContextKey) String() string {
	return string(k)
}

type ContextStringValue struct {
	s string
}

func (s *ContextStringValue) String() string {
	return s.s
}

func MakeContextStringValue(s string) ContextStringValue {
	return ContextStringValue{s: s}
}

type SessionContext struct {
	context.Context
	_session *entityAuth.Session
}

func NewSessionContext(session *entityAuth.Session) *SessionContext {
	return &SessionContext{_session: session}
}

func WithSession(parent context.Context, value *entityAuth.Session) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithValue(parent, sessionContextKey, NewSessionContext(value))
}

func (sc *SessionContext) Session() *entityAuth.Session {
	if sc == nil || sc._session == nil {
		return nil
	}
	return sc._session
}

func SessionFromContext(ctx context.Context) (*entityAuth.Session, error) {
	sessionCtx := ctx.Value(sessionContextKey)
	if sessionCtx == nil {
		return nil, errors.New("session context not found")
	}
	if sessionCtx, ok := sessionCtx.(*SessionContext); ok {
		if session := sessionCtx.Session(); session != nil {
			return session, nil
		}
	}
	return nil, errors.New("session not found")
}

type DeviceContext struct {
	context.Context
	_device *entityAuth.Device
}

func NewDeviceContext(device *entityAuth.Device) *DeviceContext {
	return &DeviceContext{_device: device}
}

func WithDevice(parent context.Context, value *entityAuth.Device) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithValue(parent, deviceContextKey, NewDeviceContext(value))
}

func (dc *DeviceContext) Device() *entityAuth.Device {
	if dc == nil || dc._device == nil {
		return nil
	}
	return dc._device
}

func DeviceFromContext(ctx context.Context) (*entityAuth.Device, error) {
	deviceCtx := ctx.Value(deviceContextKey)
	if deviceCtx == nil {
		return nil, errors.New("device context not found")
	}
	if deviceCtx, ok := deviceCtx.(*DeviceContext); ok {
		if session := deviceCtx.Device(); session != nil {
			return session, nil
		}
	}
	return nil, errors.New("device not found")
}

type ClientIPContext struct {
	context.Context
	_clientIP net.IP
}

func NewClientIPContext(clientIP net.IP) *ClientIPContext {
	return &ClientIPContext{_clientIP: clientIP}
}

func WithClientIP(parent context.Context, value net.IP) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithValue(parent, clientIPContextKey, NewClientIPContext(value))
}

func (cip *ClientIPContext) ClientIP() net.IP {
	if cip == nil || cip._clientIP == nil {
		return nil
	}
	return cip._clientIP
}

func ClientIPFromContext(ctx context.Context) (net.IP, error) {
	clientIPCtx := ctx.Value(clientIPContextKey)
	if clientIPCtx == nil {
		return nil, errors.New("client ip context not found")
	}
	if clientIPCtx, ok := clientIPCtx.(*ClientIPContext); ok {
		if clientIP := clientIPCtx.ClientIP(); clientIP != nil {
			return clientIP, nil
		}
	}
	return nil, errors.New("client ip not found")
}
