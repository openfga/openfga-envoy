package extractor

import (
	"context"
	"testing"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestSpiffeUnmarshal(t *testing.T) {
	c := &SpiffeConfig{}
	yaml.Unmarshal([]byte("type: user"), c)
	require.Equal(t, spiffeTypeUser, c.Type)
}

func TestSpiffeExtractor(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		extractor := NewSpiffe(nil)

		extraction, found, err := extractor(context.Background(), &authv3.CheckRequest{
			Attributes: &authv3.AttributeContext{
				Request: &authv3.AttributeContext_Request{
					Http: &authv3.AttributeContext_HttpRequest{},
				},
			},
		})

		require.NoError(t, err)
		require.False(t, found)
		require.Empty(t, extraction.Value)
	})

	t.Run("success for subject", func(t *testing.T) {
		extractor := NewSpiffe(&SpiffeConfig{
			Type: spiffeTypeUser,
		})

		extraction, found, err := extractor(context.Background(), &authv3.CheckRequest{
			Attributes: &authv3.AttributeContext{
				Request: &authv3.AttributeContext_Request{
					Http: &authv3.AttributeContext_HttpRequest{
						Headers: map[string]string{
							"x-forwarded-client-cert": "Hash=519dbf0d617cd943359dcf71f4d26d35e95347b616e62dc9c5ce4f3a7492ec76;Cert=\"-----BEGIN%20CERTIFICATE-----%0AMIIC3DCCAcQCAQEwDQYJKoZIhvcNAQEFBQAwLTEVMBMGA1UECgwMZXhhbXBsZSBJ%0AbmMuMRQwEgYDVQQDDAtleGFtcGxlLmNvbTAeFw0yMTA2MDcxNDUwMDBaFw0yMjA2%0AMDcxNDUwMDBaMDsxGzAZBgNVBAMMEmNsaWVudC5leGFtcGxlLmNvbTEcMBoGA1UE%0ACgwTY2xpZW50IG9yZ2FuaXphdGlvbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC%0AAQoCggEBAMe%2FVGy3syX%2BkMpqe6MfStPuFlwwEzwqz3lAaMm9YqASOk9uv0Qc%2FZPm%0AwIMrV1dnnLtLo6nZaRfsMgz1XiYBp%2BR2O87dFw2WN5AjD98zUT3XqUzGHF63cZvH%0AmoVkqrHiwh35HCFiwh9KIS3CtIdYe1n1%2FSkJpj0tVszIY%2Bi288Hu5K5fVYyYzyk%2F%0ANYQKG7A1KEgZ39SkRCHIyK%2FWeGhNT1VCfn%2Bx3D62RUVDjK6Au9yu5pJsKAyB76eg%0AZkXdvBv92dIN2iPhe9DzJ99MfpkjG9JfZG43svc2I2BAIj1eLzh5ZUnpILIHP42S%0AQUd6IBX9X%2BvFH%2FsFVnoySKewcMRK%2BkcCAwEAATANBgkqhkiG9w0BAQUFAAOCAQEA%0AswYNs6M4UT8epcW3MOjA%2B5c8EWfMI7SjNuShpRUNaJAhzdvpfj3PFIW%2BROLNs0Tj%0AwGVwkkZW8daxBQw8yC9kEE%2Broj7eJmV9SE%2BozZwa6L4hf18pcNaJlKIyvQUS3mgB%0ApGYO9YvC%2Bsg%2B0gfbSWfbzL17jRS1UI%2BOiW%2BWS5o85SOpusSHDtrG4qcISm7jpgyb%0AudzCZQHOkknO4e%2BrWiGKLpGBE1LkS5Cl%2FJkU1qJWspa4JaFtQxNCdT2Tmo6XDRZ7%0AKfoZiH6c1lI7C07duz9iPkNATc2w%2BNP7bzQgp4BlC0zQ3MwEbcR5uVxvC3vTRsIa%0AznXIRj23jj3NmidA4DTASQ%3D%3D%0A-----END%20CERTIFICATE-----%0A\";Subject=\"O=client organization,CN=client.example.com\";URI=,By=spiffe://cluster.local/ns/default/sa/httpbin;Hash=58531cf54811dc1fd60ee4aaea52866daecb353cb23d2fa237c580cbc217b4be;Subject=\"\";URI=spiffe://cluster.local/ns/istio-system/sa/istio-ingressgateway-service-account",
						},
					},
				},
			},
		})

		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, "spiffe://cluster.local/ns/istio-system/sa/istio-ingressgateway-service-account", extraction.Value)
	})

	t.Run("success for object", func(t *testing.T) {
		extractor := NewSpiffe(&SpiffeConfig{
			Type: spiffeTypeObject,
		})

		extraction, found, err := extractor(context.Background(), &authv3.CheckRequest{
			Attributes: &authv3.AttributeContext{
				Request: &authv3.AttributeContext_Request{
					Http: &authv3.AttributeContext_HttpRequest{
						Headers: map[string]string{
							"x-forwarded-client-cert": "Hash=519dbf0d617cd943359dcf71f4d26d35e95347b616e62dc9c5ce4f3a7492ec76;Cert=\"-----BEGIN%20CERTIFICATE-----%0AMIIC3DCCAcQCAQEwDQYJKoZIhvcNAQEFBQAwLTEVMBMGA1UECgwMZXhhbXBsZSBJ%0AbmMuMRQwEgYDVQQDDAtleGFtcGxlLmNvbTAeFw0yMTA2MDcxNDUwMDBaFw0yMjA2%0AMDcxNDUwMDBaMDsxGzAZBgNVBAMMEmNsaWVudC5leGFtcGxlLmNvbTEcMBoGA1UE%0ACgwTY2xpZW50IG9yZ2FuaXphdGlvbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC%0AAQoCggEBAMe%2FVGy3syX%2BkMpqe6MfStPuFlwwEzwqz3lAaMm9YqASOk9uv0Qc%2FZPm%0AwIMrV1dnnLtLo6nZaRfsMgz1XiYBp%2BR2O87dFw2WN5AjD98zUT3XqUzGHF63cZvH%0AmoVkqrHiwh35HCFiwh9KIS3CtIdYe1n1%2FSkJpj0tVszIY%2Bi288Hu5K5fVYyYzyk%2F%0ANYQKG7A1KEgZ39SkRCHIyK%2FWeGhNT1VCfn%2Bx3D62RUVDjK6Au9yu5pJsKAyB76eg%0AZkXdvBv92dIN2iPhe9DzJ99MfpkjG9JfZG43svc2I2BAIj1eLzh5ZUnpILIHP42S%0AQUd6IBX9X%2BvFH%2FsFVnoySKewcMRK%2BkcCAwEAATANBgkqhkiG9w0BAQUFAAOCAQEA%0AswYNs6M4UT8epcW3MOjA%2B5c8EWfMI7SjNuShpRUNaJAhzdvpfj3PFIW%2BROLNs0Tj%0AwGVwkkZW8daxBQw8yC9kEE%2Broj7eJmV9SE%2BozZwa6L4hf18pcNaJlKIyvQUS3mgB%0ApGYO9YvC%2Bsg%2B0gfbSWfbzL17jRS1UI%2BOiW%2BWS5o85SOpusSHDtrG4qcISm7jpgyb%0AudzCZQHOkknO4e%2BrWiGKLpGBE1LkS5Cl%2FJkU1qJWspa4JaFtQxNCdT2Tmo6XDRZ7%0AKfoZiH6c1lI7C07duz9iPkNATc2w%2BNP7bzQgp4BlC0zQ3MwEbcR5uVxvC3vTRsIa%0AznXIRj23jj3NmidA4DTASQ%3D%3D%0A-----END%20CERTIFICATE-----%0A\";Subject=\"O=client organization,CN=client.example.com\";URI=,By=spiffe://cluster.local/ns/default/sa/httpbin;Hash=58531cf54811dc1fd60ee4aaea52866daecb353cb23d2fa237c580cbc217b4be;Subject=\"\";URI=spiffe://cluster.local/ns/istio-system/sa/istio-ingressgateway-service-account",
						},
					},
				},
			},
		})

		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, "spiffe://cluster.local/ns/default/sa/httpbin", extraction.Value)
	})
}
