package hub

import "github.com/aws/aws-sdk-go/aws/client"

type Hub struct {
	clients []client.Client
}
