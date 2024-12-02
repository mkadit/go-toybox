package models

type (
	Configuration struct {
		Db   DbConfiguration
		Auth struct {
			JwtSecret string
		}
		Server struct {
			CorsAllowOrigin string
			Port            int
		}
		Email EmailConfiguration
		Tcp   TCPConfiguration
		Rpc   RPCConfiguration
	}

	DbConfiguration struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
		Driver   string
	}

	EmailConfiguration struct {
		Address  string
		Password string
		Server   string
		Port     int
	}

	TCPConfiguration struct {
		Role       string
		HostIp     string
		HostPort   int
		ClientIp   string
		ClientPort int
		Timeout    int
		PoolSize   int
	}

	RPCConfiguration struct {
		Port int
	}
)
