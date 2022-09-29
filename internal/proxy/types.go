package proxy

type RunConfiguration struct {
	Tunnel   TunnelConfiguration
	Services []ServiceConfiguration
}

type TunnelConfiguration struct {
	UserAndHost         string
	CertificateLocation string
	Destination         string
}

type ServiceConfiguration struct {
	ServiceName string
	Port        int
}
