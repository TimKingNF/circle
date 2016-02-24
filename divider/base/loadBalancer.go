package base

type LoadBalancer interface {
	Clusters(clusters []uint32)
	GenLoadBalancing() uint32
}

type myLoadBalancer struct {
	clusters []uint32 // clusters list (uint32)

	currentClustersKey int
}

func NewLoadBalancer() LoadBalancer {
	return &myLoadBalancer{
		currentClustersKey: -1,
	}
}

func (lb *myLoadBalancer) Clusters(clusters []uint32) {
	lb.clusters = clusters
}

func (lb *myLoadBalancer) GenLoadBalancing() uint32 {
	if lb.clusters == nil {
		return 0
	}
	for k, _ := range lb.clusters {
		if lb.currentClustersKey == len(lb.clusters)-1 {
			lb.currentClustersKey = -1
		}
		if k == lb.currentClustersKey+1 {
			lb.currentClustersKey = k
			return lb.clusters[k]
		}
	}
	return 0
}
