package v1

import (
	"path/filepath"
)

// RemoteConnector specifies the namespace in the remote cluster
type RemoteConnector struct {
	GKE *GKEConnector `json:"gcp,omitempty" protobuf:"bytes,1,opt,name=gcp"`
}

// Path returns the relative path used for caching connectors
func (c *RemoteConnector) Path() string {
	if c.GKE != nil {
		return c.GKE.Path()
	}
	return "unknown"
}

// GKEConnector the connection details for using Google Cloud
type GKEConnector struct {
	Project string `json:"project,omitempty" protobuf:"bytes,1,opt,name=project"`
	Cluster string `json:"cluster,omitempty" protobuf:"bytes,2,opt,name=cluster"`
	Region  string `json:"region,omitempty" protobuf:"bytes,3,opt,name=region"`
	Zone    string `json:"zone,omitempty" protobuf:"bytes,4,opt,name=zone"`
}

// Path returns the relative path for the dir for any configuration files needed for this connector
func (c *GKEConnector) Path() string {
	if c.Region != "" {
		return filepath.Join("gcp", c.Project, c.Cluster, "region", c.Region)
	} else {
		return filepath.Join("gcp", c.Project, c.Cluster, "zone", c.Zone)
	}
}
