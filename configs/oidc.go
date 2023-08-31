package configs

import "github.com/namsral/flag"

var (
	OIDCClientId, OIDCClientSecret, OIDCAuthority, OIDCScope, OIDCAudience, OIDCIntrospectEndpoint string
)

func ConfigIdentityServer(commandSet *flag.FlagSet) {
	if command := commandSet.Lookup("oidc-client-id"); command != nil {
		OIDCClientId = command.Value.String()
	}
	if command := commandSet.Lookup("oidc-client-secret"); command != nil {
		OIDCClientSecret = command.Value.String()
	}
	if command := commandSet.Lookup("oidc-authority"); command != nil {
		OIDCAuthority = command.Value.String()
	}
	if command := commandSet.Lookup("oidc-scope"); command != nil {
		OIDCScope = command.Value.String()
	}
	if command := commandSet.Lookup("oidc-audience"); command != nil {
		OIDCAudience = command.Value.String()
	}
	if command := commandSet.Lookup("oidc-introspect-endpoint"); command != nil {
		OIDCIntrospectEndpoint = command.Value.String()
	}
}
