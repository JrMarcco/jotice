package client

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/JrMarcco/easy-kit/slice"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type rwWeightSvcNode struct {
	mu sync.RWMutex

	subConn              balancer.SubConn
	readWeight           int32
	currentReadWeight    int32
	efficientReadWeight  int32
	writeWeight          int32
	currentWriteWeight   int32
	efficientWriteWeight int32
	group                string
}

func newRWWeightSvcNode(subConn balancer.SubConn, readWeight int32, writeWeight int32, group string) *rwWeightSvcNode {
	return &rwWeightSvcNode{
		subConn:              subConn,
		readWeight:           readWeight,
		currentReadWeight:    readWeight,
		efficientReadWeight:  readWeight,
		writeWeight:          writeWeight,
		currentWriteWeight:   writeWeight,
		efficientWriteWeight: writeWeight,
		group:                group,
	}
}

var _ balancer.Picker = (*RWWeightBalancer)(nil)

type RWWeightBalancer struct {
	nodes []*rwWeightSvcNode
}

func (r *RWWeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(r.nodes) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	// select candidate nodes
	candidates := slice.FilterMap(r.nodes, func(_ int, src *rwWeightSvcNode) (*rwWeightSvcNode, bool) {
		src.mu.RLock()
		nodeGroup := src.group
		src.mu.RUnlock()

		return src, r.getGroup(info.Ctx) == nodeGroup
	})

	var totalWeight int32
	var selectedNode *rwWeightSvcNode

	ctx := info.Ctx
	isWrite := r.isWrite(ctx)

	for _, node := range candidates {
		node.mu.Lock()
		if isWrite {
			totalWeight += node.efficientWriteWeight
			node.currentWriteWeight += node.efficientWriteWeight
			if selectedNode == nil || selectedNode.currentWriteWeight < node.currentWriteWeight {
				selectedNode = node
			}
		} else {
			totalWeight += node.efficientReadWeight
			node.currentReadWeight += node.efficientReadWeight
			if selectedNode == nil || selectedNode.currentReadWeight < node.currentReadWeight {
				selectedNode = node
			}
		}
		node.mu.Unlock()
	}

	if selectedNode == nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	selectedNode.mu.Lock()
	if isWrite {
		selectedNode.currentWriteWeight -= totalWeight
	} else {
		selectedNode.currentReadWeight -= totalWeight
	}
	selectedNode.mu.Unlock()

	return balancer.PickResult{
		SubConn: selectedNode.subConn,
		Done: func(info balancer.DoneInfo) {
			selectedNode.mu.Lock()
			defer selectedNode.mu.Unlock()

			isDecrementErr := info.Err == nil && (errors.Is(info.Err, context.DeadlineExceeded) || errors.Is(info.Err, io.EOF))
			const twice = 2
			if isWrite {
				if isDecrementErr && selectedNode.efficientWriteWeight > 0 {
					if selectedNode.efficientWriteWeight > 0 {
						selectedNode.efficientWriteWeight--
					}
				} else if info.Err == nil {
					selectedNode.efficientWriteWeight++
					selectedNode.efficientWriteWeight = max(selectedNode.efficientWriteWeight, selectedNode.writeWeight*twice)
				}
			} else {
				if isDecrementErr && selectedNode.efficientReadWeight > 0 {
					if selectedNode.efficientReadWeight > 0 {
						selectedNode.efficientReadWeight--
					}
				} else if info.Err == nil {
					selectedNode.efficientReadWeight++
					selectedNode.efficientReadWeight = max(selectedNode.efficientReadWeight, selectedNode.readWeight*twice)
				}
			}
		},
	}, nil
}

func (r *RWWeightBalancer) getGroup(ctx context.Context) string {
	val := ctx.Value(attrNameGroup)
	if val == nil {
		return ""
	}

	if valString, ok := val.(string); ok {
		return valString
	}

	return ""
}

func (r *RWWeightBalancer) isWrite(ctx context.Context) bool {
	val := ctx.Value(KeyRequestType)
	if val == nil {
		return false
	}

	if valInt, ok := val.(int); ok {
		return valInt == 1
	}

	return false
}

type RWWeightBalancerBuilder struct {
	mu        sync.RWMutex
	nodeCache map[string]*rwWeightSvcNode
}

func (r *RWWeightBalancerBuilder) Builder(info base.PickerBuildInfo) balancer.Picker {
	nodes := make([]*rwWeightSvcNode, 0, len(info.ReadySCs))

	subConnMap := make(map[string]struct{})

	r.mu.Lock()
	defer r.mu.Unlock()

	for subConn, subConnInfo := range info.ReadySCs {
		readWeight, ok := subConnInfo.Address.Attributes.Value(attrNameReadWeight).(int32)
		if !ok {
			continue
		}
		writeWeight, ok := subConnInfo.Address.Attributes.Value(attrNameWriteWeight).(int32)
		if !ok {
			continue
		}
		group, ok := subConnInfo.Address.Attributes.Value(attrNameGroup).(string)
		if !ok {
			continue
		}
		nodeName, ok := subConnInfo.Address.Attributes.Value(attrNameNode).(string)
		if !ok {
			continue
		}

		subConnMap[nodeName] = struct{}{}

		// check if the node exists in the cache
		if cachedNode, ok := r.nodeCache[nodeName]; ok {
			// exists in the cache
			// update connection info and group name
			// but retains the weight
			cachedNode.mu.Lock()
			cachedNode.group = group
			cachedNode.mu.Unlock()

			if cachedNode.readWeight != readWeight || cachedNode.writeWeight != writeWeight {
				cachedNode = newRWWeightSvcNode(subConn, readWeight, writeWeight, group)
				r.nodeCache[nodeName] = cachedNode
			}

			nodes = append(nodes, cachedNode)
		} else {
			// not exists in the cache
			// create a new node
			newNode := newRWWeightSvcNode(subConn, readWeight, writeWeight, group)

			// cache the new node
			r.nodeCache[nodeName] = newNode
			nodes = append(nodes, newNode)
		}
	}

	// clean the node not exists anymore.
	for key := range r.nodeCache {
		if _, ok := subConnMap[key]; !ok {
			delete(r.nodeCache, key)
		}
	}

	return &RWWeightBalancer{
		nodes: nodes,
	}
}

func NewRWWeightBalancerBuilder() *RWWeightBalancerBuilder {
	return &RWWeightBalancerBuilder{
		nodeCache: make(map[string]*rwWeightSvcNode),
	}
}

func WithGroup(ctx context.Context, group string) context.Context {
	return context.WithValue(ctx, groupKey{}, group)
}
