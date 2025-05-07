package client

import (
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type weightSvcNode struct {
	mu sync.Mutex

	subConn       balancer.SubConn
	weight        uint32
	currentWeight uint32
}

var _ balancer.Picker = (*WeightBalancer)(nil)

type WeightBalancer struct {
	nodes       []*weightSvcNode
	totalWeight uint32
}

func (w *WeightBalancer) Pick(_ balancer.PickInfo) (balancer.PickResult, error) {
	if len(w.nodes) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	var selectedNode *weightSvcNode
	for _, node := range w.nodes {
		node.mu.Lock()
		node.currentWeight = node.currentWeight + node.weight

		if selectedNode == nil || selectedNode.currentWeight < node.currentWeight {
			selectedNode = node
		}
		node.mu.Unlock()
	}

	if selectedNode == nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	selectedNode.mu.Lock()
	selectedNode.currentWeight -= w.totalWeight
	selectedNode.mu.Unlock()

	return balancer.PickResult{
		SubConn: selectedNode.subConn,
		Done:    func(_ balancer.DoneInfo) {},
	}, nil
}

type WeightBalancerBuilder struct {
}

func (w *WeightBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	nodes := make([]*weightSvcNode, 0, len(info.ReadySCs))
	totalWeight := uint32(0)

	for subConn, subConnInfo := range info.ReadySCs {
		weight, _ := subConnInfo.Address.Attributes.Value(attrNameWeight).(uint32)
		totalWeight += weight

		nodes = append(nodes, &weightSvcNode{
			subConn:       subConn,
			weight:        weight,
			currentWeight: weight,
		})
	}

	return &WeightBalancer{
		nodes:       nodes,
		totalWeight: totalWeight,
	}
}
