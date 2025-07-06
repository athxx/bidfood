package rpc

import (
	pb "bidone/bidrpc/bidrpcproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProductClient wraps the gRPC client for product service
type ProductClient struct {
	Clt  pb.ProductServiceClient
	conn *grpc.ClientConn
}

var RpcClientProduct *ProductClient

func InitProductGrpcClient(addr string) (err error) {
	RpcClientProduct, err = NewProductClient(addr)
	return err
}

// NewProductClient creates a new product gRPC client
func NewProductClient(addr string) (*ProductClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ProductClient{
		Clt:  pb.NewProductServiceClient(conn),
		conn: conn,
	}, nil
}

// Close closes the gRPC connection
func (c *ProductClient) Close() error {
	return c.conn.Close()
}
