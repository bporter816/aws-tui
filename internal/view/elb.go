package view

type ELB struct {
}

func (e ELB) GetService() string {
	return "ELB"
}

type ELBLoadBalancers struct {
	ELB
}

func (e ELBLoadBalancers) GetName() string {
	return "Load Balancers"
}

type ELBListeners struct {
	ELB
}

func (e ELBListeners) GetName() string {
	return "Listeners"
}

type ELBTargetGroups struct {
	ELB
}

func (e ELBTargetGroups) GetName() string {
	return "Target Groups"
}

type ELBTrustStores struct {
	ELB
}

func (e ELBTrustStores) GetName() string {
	return "Trust Stores"
}

type ELBTags struct {
	ELB
}

func (e ELBTags) GetName() string {
	return "Tags"
}
