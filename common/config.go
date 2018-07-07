package common

// Config is the struct of application config.
type Config struct {
	ListenAddress            string
	MetricsPath              string
	Namespace                string
	NginxUrls                []string
	NginxPlusUrls            []string
	ExcludeUpstreamAddresses []string
}

// NewConfig creates new application config.
func NewConfig(
	listenAddress string,
	metricsPath string,
	namespace string,
	nginxUrls []string,
	nginxPlusUrls []string,
	excludeUpstreamAddresses []string,
) *Config {
	return &Config{
		ListenAddress:            listenAddress,
		MetricsPath:              metricsPath,
		Namespace:                namespace,
		NginxUrls:                nginxUrls,
		NginxPlusUrls:            nginxPlusUrls,
		ExcludeUpstreamAddresses: excludeUpstreamAddresses,
	}
}
