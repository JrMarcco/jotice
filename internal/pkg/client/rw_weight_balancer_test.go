package client

import (
	"testing"

	"google.golang.org/grpc/balancer"
)

type mockSubConn struct {
	balancer.SubConn
	name string
}

func (m *mockSubConn) Name() string {
	return m.name
}

func createTestSvcNode(name string, readWeight int32, writeWeight int32) *rwWeightSvcNode {
	return &rwWeightSvcNode{
		subConn:              &mockSubConn{name: name},
		readWeight:           readWeight,
		currentReadWeight:    readWeight,
		efficientReadWeight:  readWeight,
		writeWeight:          writeWeight,
		currentWriteWeight:   writeWeight,
		efficientWriteWeight: writeWeight,
	}
}

type reqType int

const (
	reqTypeRead reqType = iota
	reqTypeWrite
)

func TestRWWeightBalancer_Pick(t *testing.T) {
	t.Parallel()

	_ = []*rwWeightSvcNode{
		createTestSvcNode("read-1-write-4", 1, 4),
		createTestSvcNode("read-2-write-3", 2, 3),
		createTestSvcNode("read-3-write-2", 3, 2),
		createTestSvcNode("read-4-write-1", 4, 1),
	}

	_ = []struct {
		requestType reqType
		name        string
		wantNode    string
		wantErr     error
	}{{}}

}
