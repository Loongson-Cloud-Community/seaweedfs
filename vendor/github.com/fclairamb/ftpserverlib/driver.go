// Package ftpserver provides all the tools to build your own FTP server: The core library and the driver.
package ftpserver

import (
	"crypto/tls"
	"net"

	"github.com/spf13/afero"
)

// This file is the driver part of the server. It must be implemented by anyone wanting to use the server.

// MainDriver handles the authentication and ClientHandlingDriver selection
type MainDriver interface {
	// GetSettings returns some general settings around the server setup
	GetSettings() (*Settings, error)

	// ClientConnected is called to send the very first welcome message
	ClientConnected(cc ClientContext) (string, error)

	// ClientDisconnected is called when the user disconnects, even if he never authenticated
	ClientDisconnected(cc ClientContext)

	// AuthUser authenticates the user and selects an handling driver
	AuthUser(cc ClientContext, user, pass string) (ClientDriver, error)

	// GetTLSConfig returns a TLS Certificate to use
	// The certificate could frequently change if we use something like "let's encrypt"
	GetTLSConfig() (*tls.Config, error)
}

// ClientDriver is the base FS implementation that allows to manipulate files
type ClientDriver interface {
	afero.Fs
}

// ClientDriverExtensionAllocate is an extension to allow to support the allocation command
type ClientDriverExtensionAllocate interface {

	// AllocateSpace reserves the space necessary to upload files
	AllocateSpace(size int) error
}

// ClientDriverExtensionChown is an extension to allow to support the chown command
type ClientDriverExtensionChown interface {

	// Chown changes the owner of a file
	Chown(name string, user string, group string) error
}

// ClientContext is implemented on the server side to provide some access to few data around the client
type ClientContext interface {
	// Path provides the path of the current connection
	Path() string

	// SetDebug activates the debugging of this connection commands
	SetDebug(debug bool)

	// Debug returns the current debugging status of this connection commands
	Debug() bool

	// Client's ID on the server
	ID() uint32

	// Client's address
	RemoteAddr() net.Addr

	// Servers's address
	LocalAddr() net.Addr

	// Client's version can be empty
	GetClientVersion() string
}

// PortRange is a range of ports
type PortRange struct {
	Start int // Range start
	End   int // Range end
}

// PublicIPResolver takes a ClientContext for a connection and returns the public IP
// to use in the response to the PASV command, or an error if a public IP cannot be determined.
type PublicIPResolver func(ClientContext) (string, error)

// Settings defines all the server settings
// nolint: maligned
type Settings struct {
	Listener                 net.Listener     // (Optional) To provide an already initialized listener
	ListenAddr               string           // Listening address
	PublicHost               string           // Public IP to expose (only an IP address is accepted at this stage)
	PublicIPResolver         PublicIPResolver // (Optional) To fetch a public IP lookup
	PassiveTransferPortRange *PortRange       // (Optional) Port Range for data connections. Random if not specified
	ActiveTransferPortNon20  bool             // Do not impose the port 20 for active data transfer (#88, RFC 1579)
	IdleTimeout              int              // Maximum inactivity time before disconnecting (#58)
	ConnectionTimeout        int              // Maximum time to establish passive or active transfer connections
	DisableMLSD              bool             // Disable MLSD support
	DisableMLST              bool             // Disable MLST support
	DisableMFMT              bool             // Disable MFMT support (modify file mtime)
}
