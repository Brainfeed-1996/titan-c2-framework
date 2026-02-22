package transport

import (
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Connect returns a grpc client connection to the C2 server
func Connect(address string, secure bool) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	if secure {
		// In production, load CA certs properly
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, // For demo/dev purposes
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithTimeout(5*time.Second))

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
