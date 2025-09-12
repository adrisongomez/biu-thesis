package models

type ServiceNode struct {
	Name string `json:"name"`
}

func (s *ServiceNode) Raw() map[string]any {
	return map[string]any{"name": s.Name}
}

func NewServiceNode(name string) ServiceNode {
	return ServiceNode{
		Name: name,
	}
}
