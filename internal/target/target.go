// Package target names the fixed runtime Docker Zero manages: one WSL2
// distro, one Linux user, and the Portainer container inside it.
package target

const (
	// Distro is the WSL2 distro Docker Zero installs and operates.
	Distro = "Ubuntu-24.04"
	// User is the default Linux user inside Distro.
	User = "dev"
	// PortainerContainer is the name of the Portainer CE container.
	PortainerContainer = "portainer"
	// PortainerURL is the page Open Portainer shows. Portainer is published
	// on 127.0.0.1:9443 only; the engine has no TCP port.
	PortainerURL = "https://localhost:9443/"
	// PortainerStatusURL is Portainer's unauthenticated status endpoint.
	PortainerStatusURL = "https://127.0.0.1:9443/api/system/status"
)
