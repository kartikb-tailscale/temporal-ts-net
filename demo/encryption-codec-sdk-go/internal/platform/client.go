package platform

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

func Dial(address, namespace, codecEndpoint string) (client.Client, error) {
	dc := converter.GetDefaultDataConverter()
	if codecEndpoint != "" {
		dc = converter.NewRemoteDataConverter(dc, converter.RemoteDataConverterOptions{Endpoint: codecEndpoint})
	}

	return client.Dial(client.Options{
		HostPort:      address,
		Namespace:     namespace,
		DataConverter: dc,
	})
}
