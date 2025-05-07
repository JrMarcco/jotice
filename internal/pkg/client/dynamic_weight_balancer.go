package client

import (
	"context"
	"errors"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type dynamicSvcNode struct {
	mu sync.Mutex

	subConn         balancer.SubConn
	weight          uint32
	currentWeight   uint32
	efficientWeight uint32
}

var _ balancer.Picker = (*DynamicWeightBalancer)(nil)

type DynamicWeightBalancer struct {
	nodes []*dynamicSvcNode
}

func (d DynamicWeightBalancer) Pick(_ balancer.PickInfo) (balancer.PickResult, error) {
	if len(d.nodes) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	var totalWeight uint32
	var selectedNode *dynamicSvcNode

	for _, node := range d.nodes {
		node.mu.Lock()
		totalWeight += node.efficientWeight
		node.currentWeight += node.efficientWeight

		if selectedNode == nil || selectedNode.currentWeight < node.currentWeight {
			selectedNode = node
		}
		node.mu.Unlock()
	}

	if selectedNode == nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	selectedNode.mu.Lock()
	selectedNode.currentWeight -= totalWeight
	selectedNode.mu.Unlock()

	return balancer.PickResult{
		SubConn: selectedNode.subConn,
		Done: func(info balancer.DoneInfo) {
			selectedNode.mu.Lock()
			defer selectedNode.mu.Unlock()

			if info.Err == nil {
				selectedNode.efficientWeight++

				const twice = 2
				selectedNode.efficientWeight = max(selectedNode.efficientWeight, selectedNode.weight*twice)
				return
			}

			if errors.Is(info.Err, context.DeadlineExceeded) {
				selectedNode.efficientWeight = 1
				return
			}

			statusCode, _ := status.FromError(info.Err)
			switch statusCode.Code() {
			case codes.Unavailable:
				selectedNode.efficientWeight = 1
				return
			default:
				if selectedNode.efficientWeight > 1 {
					selectedNode.efficientWeight--
				}
			}
		},
	}, nil

}

type DynamicWeightBalancerBuilder struct{}

func (d *DynamicWeightBalancerBuilder) Builder(info base.PickerBuildInfo) balancer.Picker {
	nodes := make([]*dynamicSvcNode, 0, len(info.ReadySCs))

	for subConn, subConnInfo := range info.ReadySCs {
		weight, _ := subConnInfo.Address.Attributes.Value(attrNameWeight).(uint32)

		nodes = append(nodes, &dynamicSvcNode{
			subConn:         subConn,
			weight:          weight,
			currentWeight:   weight,
			efficientWeight: weight,
		})
	}

	return &DynamicWeightBalancer{
		nodes: nodes,
	}
}
