package dto

import (
	"github.com/quiqxiq/roslib-mikhmon/domain"
	"github.com/quiqxiq/roslib-mikhmon/mikrotik/ppp"
)

// ── PPP Secret ─────────────────────────────────────────────────────────

type PPPSecretResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Service          string `json:"service,omitempty"`
	Profile          string `json:"profile,omitempty"`
	LocalAddr        string `json:"local_address,omitempty"`
	RemoteAddr       string `json:"remote_address,omitempty"`
	CallerID         string `json:"caller_id,omitempty"`
	Routes           string `json:"routes,omitempty"`
	IPv6Routes       string `json:"ipv6_routes,omitempty"`
	RemoteIPv6Prefix string `json:"remote_ipv6_prefix,omitempty"`
	LimitBytesIn     int64  `json:"limit_bytes_in"`
	LimitBytesOut    int64  `json:"limit_bytes_out"`
	LastLoggedOut    string `json:"last_logged_out,omitempty"`
	LastCallerID     string `json:"last_caller_id,omitempty"`
	Disabled         bool   `json:"disabled"`
	Comment          string `json:"comment,omitempty"`
	// Password sengaja tidak di-expose.
}

func FromDomainPPPSecret(s domain.PPPSecret) PPPSecretResponse {
	return PPPSecretResponse{
		ID: s.ID, Name: s.Name, Service: s.Service, Profile: s.Profile,
		LocalAddr: s.LocalAddr, RemoteAddr: s.RemoteAddr, CallerID: s.CallerID,
		Routes: s.Routes, IPv6Routes: s.IPv6Routes, RemoteIPv6Prefix: s.RemoteIPv6Prefix,
		LimitBytesIn: s.LimitBytesIn, LimitBytesOut: s.LimitBytesOut,
		LastLoggedOut: s.LastLoggedOut, LastCallerID: s.LastCallerID,
		Disabled: s.Disabled, Comment: s.Comment,
	}
}

func FromDomainPPPSecrets(ss []domain.PPPSecret) []PPPSecretResponse {
	out := make([]PPPSecretResponse, len(ss))
	for i, s := range ss {
		out[i] = FromDomainPPPSecret(s)
	}
	return out
}

type PPPSecretCreateRequest struct {
	Name       string `json:"name"                     binding:"required,min=1,max=128"`
	Password   string `json:"password,omitempty"`
	Service    string `json:"service,omitempty"`
	Profile    string `json:"profile,omitempty"`
	LocalAddr  string `json:"local_address,omitempty"`
	RemoteAddr string `json:"remote_address,omitempty"`
	CallerID   string `json:"caller_id,omitempty"`
	Disabled   *bool  `json:"disabled,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

func (r PPPSecretCreateRequest) ToArgs() ppp.SecretAddArgs {
	return ppp.SecretAddArgs{
		Name: r.Name, Password: r.Password, Service: r.Service, Profile: r.Profile,
		LocalAddr: r.LocalAddr, RemoteAddr: r.RemoteAddr,
		CallerID: r.CallerID, Disabled: r.Disabled, Comment: r.Comment,
	}
}

type PPPSecretUpdateRequest struct {
	Name       string  `json:"name,omitempty"`
	Password   string  `json:"password,omitempty"`
	Service    string  `json:"service,omitempty"`
	Profile    string  `json:"profile,omitempty"`
	LocalAddr  string  `json:"local_address,omitempty"`
	RemoteAddr string  `json:"remote_address,omitempty"`
	CallerID   *string `json:"caller_id,omitempty"`
	Disabled   *bool   `json:"disabled,omitempty"`
	Comment    *string `json:"comment,omitempty"`
}

func (r PPPSecretUpdateRequest) ToArgs(id string) ppp.SecretSetArgs {
	return ppp.SecretSetArgs{
		ID: id, Name: r.Name, Password: r.Password, Service: r.Service, Profile: r.Profile,
		LocalAddr: r.LocalAddr, RemoteAddr: r.RemoteAddr,
		CallerID: r.CallerID, Disabled: r.Disabled, Comment: r.Comment,
	}
}

// ── PPP Profile ────────────────────────────────────────────────────────

type PPPProfileResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	LocalAddr      string `json:"local_address,omitempty"`
	RemoteAddr     string `json:"remote_address,omitempty"`
	RateLimit      string `json:"rate_limit,omitempty"`
	DNSServer      string `json:"dns_server,omitempty"`
	Bridge         string `json:"bridge,omitempty"`
	ParentQueue    string `json:"parent_queue,omitempty"`
	IdleTimeout    string `json:"idle_timeout,omitempty"`
	SessionTimeout string `json:"session_timeout,omitempty"`
	OnUp           string `json:"on_up,omitempty"`
	OnDown         string `json:"on_down,omitempty"`
	OnlyOne        string `json:"only_one,omitempty"`
	UseCompression string `json:"use_compression,omitempty"`
	UseEncryption  string `json:"use_encryption,omitempty"`
	ChangeTCPMSS   string `json:"change_tcp_mss,omitempty"`
	Comment        string `json:"comment,omitempty"`
}

func FromDomainPPPProfile(p domain.PPPProfile) PPPProfileResponse {
	return PPPProfileResponse{
		ID: p.ID, Name: p.Name, LocalAddr: p.LocalAddr, RemoteAddr: p.RemoteAddr,
		RateLimit: p.RateLimit, DNSServer: p.DNSServer, Bridge: p.Bridge,
		ParentQueue: p.ParentQueue, IdleTimeout: p.IdleTimeout, SessionTimeout: p.SessionTimeout,
		OnUp: p.OnUp, OnDown: p.OnDown, OnlyOne: p.OnlyOne,
		UseCompression: p.UseCompression, UseEncryption: p.UseEncryption,
		ChangeTCPMSS: p.ChangeTCPMSS, Comment: p.Comment,
	}
}

func FromDomainPPPProfiles(ps []domain.PPPProfile) []PPPProfileResponse {
	out := make([]PPPProfileResponse, len(ps))
	for i, p := range ps {
		out[i] = FromDomainPPPProfile(p)
	}
	return out
}

type PPPProfileCreateRequest struct {
	Name       string `json:"name"      binding:"required,min=1,max=128"`
	LocalAddr  string `json:"local_address,omitempty"`
	RemoteAddr string `json:"remote_address,omitempty"`
	RateLimit  string `json:"rate_limit,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

func (r PPPProfileCreateRequest) ToArgs() ppp.ProfileAddArgs {
	return ppp.ProfileAddArgs{
		Name: r.Name, LocalAddr: r.LocalAddr, RemoteAddr: r.RemoteAddr,
		RateLimit: r.RateLimit, Comment: r.Comment,
	}
}

type PPPProfileUpdateRequest struct {
	Name       string  `json:"name,omitempty"`
	LocalAddr  string  `json:"local_address,omitempty"`
	RemoteAddr string  `json:"remote_address,omitempty"`
	RateLimit  string  `json:"rate_limit,omitempty"`
	Comment    *string `json:"comment,omitempty"`
}

func (r PPPProfileUpdateRequest) ToArgs(id string) ppp.ProfileSetArgs {
	return ppp.ProfileSetArgs{
		ID: id, Name: r.Name, LocalAddr: r.LocalAddr, RemoteAddr: r.RemoteAddr,
		RateLimit: r.RateLimit, Comment: r.Comment,
	}
}

// ── PPP Active ─────────────────────────────────────────────────────────

type PPPActiveResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Service       string `json:"service,omitempty"`
	CallerID      string `json:"caller_id,omitempty"`
	Address       string `json:"address,omitempty"`
	Uptime        string `json:"uptime,omitempty"`
	Encoding      string `json:"encoding,omitempty"`
	SessionID     string `json:"session_id,omitempty"`
	LimitBytesIn  int64  `json:"limit_bytes_in"`
	LimitBytesOut int64  `json:"limit_bytes_out"`
	Comment       string `json:"comment,omitempty"`
}

func FromDomainPPPActive(a domain.PPPActive) PPPActiveResponse {
	return PPPActiveResponse{
		ID: a.ID, Name: a.Name, Service: a.Service,
		CallerID: a.CallerID, Address: a.Address, Uptime: a.Uptime,
		Encoding: a.Encoding, SessionID: a.SessionID,
		LimitBytesIn: a.LimitBytesIn, LimitBytesOut: a.LimitBytesOut,
		Comment: a.Comment,
	}
}

func FromDomainPPPActives(as []domain.PPPActive) []PPPActiveResponse {
	out := make([]PPPActiveResponse, len(as))
	for i, a := range as {
		out[i] = FromDomainPPPActive(a)
	}
	return out
}
