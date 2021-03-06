package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// KubernetesInstance represents a single node/master within a Kubernetes cluster
type KubernetesInstance struct {
	Hostname   string    `json:"hostname"`
	Size       string    `json:"size"`
	Region     string    `json:"region"`
	CreatedAt  time.Time `json:"created_at"`
	Status     string    `json:"status"`
	FirewallID string    `json:"firewall_id"`
	PublicIP   string    `json:"public_ip"`
	Tags       []string  `json:"tags"`
}

// KubernetesInstalledApplication is an application within our marketplace available for
// installation
type KubernetesInstalledApplication struct {
	Application   string            `json:"application"`
	Title         string            `json:"title,omitempty"`
	Version       string            `json:"version"`
	Dependencies  []string          `json:"dependencies,omitempty"`
	Maintainer    string            `json:"maintainer"`
	Description   string            `json:"description"`
	PostInstall   string            `json:"post_install"`
	Installed     bool              `json:"installed"`
	URL           string            `json:"url"`
	Category      string            `json:"category"`
	UpdatedAt     time.Time         `json:"updated_at"`
	ImageURL      string            `json:"image_url"`
	Plan          string            `json:"plan,omitempty"`
	Configuration map[string]string `json:"configuration,omitempty"`
}

// KubernetesCluster is a Kubernetes k3s cluster
type KubernetesCluster struct {
	ID                    string                           `json:"id"`
	Name                  string                           `json:"name"`
	Version               string                           `json:"version"`
	Status                string                           `json:"status"`
	Ready                 bool                             `json:"ready"`
	NumTargetNode         int                              `json:"num_target_nodes"`
	TargetNodeSize        string                           `json:"target_nodes_size"`
	BuiltAt               time.Time                        `json:"built_at"`
	KubeConfig            string                           `json:"kubeconfig"`
	KubernetesVersion     string                           `json:"kubernetes_version"`
	APIEndPoint           string                           `json:"api_endpoint"`
	DNSEntry              string                           `json:"dns_entry"`
	Tags                  []string                         `json:"tags"`
	CreatedAt             time.Time                        `json:"created_at"`
	Instances             []KubernetesInstance             `json:"instances"`
	InstalledApplications []KubernetesInstalledApplication `json:"installed_applications"`
}

// KubernetesClusterConfig is used to create a new cluster
type KubernetesClusterConfig struct {
	Name              string `form:"name"`
	NumTargetNodes    int    `form:"num_target_nodes"`
	TargetNodesSize   string `form:"target_nodes_size"`
	KubernetesVersion string `form:"kubernetes_version"`
	Tags              string `form:"tags"`
	Applications      string `form:"applications"`
}

// KubernetesPlanConfiguration is a value within a configuration for
// an application's plan
type KubernetesPlanConfiguration struct {
	Value string `json:"value"`
}

// KubernetesMarketplacePlan is a plan for
type KubernetesMarketplacePlan struct {
	Label         string                                 `json:"label"`
	Configuration map[string]KubernetesPlanConfiguration `json:"configuration"`
}

// KubernetesMarketplaceApplication is an application within our marketplace
// available for installation
type KubernetesMarketplaceApplication struct {
	Name         string                      `json:"name"`
	Title        string                      `json:"title,omitempty"`
	Version      string                      `json:"version"`
	Default      string                      `json:"default,omitempty"`
	Dependencies []string                    `json:"dependencies,omitempty"`
	Maintainer   string                      `json:"maintainer"`
	Description  string                      `json:"description"`
	PostInstall  string                      `json:"post_install"`
	URL          string                      `json:"url"`
	Category     string                      `json:"category"`
	Plans        []KubernetesMarketplacePlan `json:"plans"`
}

// KubernetesVersion represents an available version of k3s to install
type KubernetesVersion struct {
	Version string `json:"version"`
	Type    string `json:"type"`
	Default bool   `json:"default,omitempty"`
}

// ListKubernetesClusters returns all cluster of kubernetes in the account
func (c *Client) ListKubernetesClusters() ([]KubernetesCluster, error) {
	resp, err := c.SendGetRequest("/v2/kubernetes/clusters")
	if err != nil {
		return nil, err
	}

	kubernetes := make([]KubernetesCluster, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&kubernetes); err != nil {
		return nil, err
	}

	return kubernetes, nil
}

// FindKubernetesCluster finds a Kubernetes cluster by either part of the ID or part of the name
func (c *Client) FindKubernetesCluster(search string) (*KubernetesCluster, error) {
	clusters, err := c.ListKubernetesClusters()
	if err != nil {
		return nil, err
	}

	found := -1

	for i, cluster := range clusters {
		if strings.Contains(cluster.ID, search) || strings.Contains(cluster.Name, search) {
			if found != -1 {
				return nil, fmt.Errorf("unable to find %s because there were multiple matches", search)
			}
			found = i
		}
	}

	if found == -1 {
		return nil, fmt.Errorf("unable to find %s, zero matches", search)
	}

	return &clusters[found], nil
}

// NewKubernetesClusters create a new cluster of kubernetes
func (c *Client) NewKubernetesClusters(kc *KubernetesClusterConfig) (*KubernetesCluster, error) {
	body, err := c.SendPostRequest("/v2/kubernetes/clusters", kc)
	if err != nil {
		return nil, err
	}

	kubernetes := &KubernetesCluster{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(kubernetes); err != nil {
		return nil, err
	}

	return kubernetes, nil
}

// GetKubernetesClusters returns a single kubernetes cluster by its full ID
func (c *Client) GetKubernetesClusters(id string) (*KubernetesCluster, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/kubernetes/clusters/%s", id))
	if err != nil {
		return nil, err
	}

	kubernetes := &KubernetesCluster{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(kubernetes)
	return kubernetes, nil
}

// UpdateKubernetesCluster update a single kubernetes cluster by its full ID
func (c *Client) UpdateKubernetesCluster(id string, i *KubernetesClusterConfig) (*KubernetesCluster, error) {
	params := map[string]interface{}{
		"name":             i.Name,
		"num_target_nodes": i.NumTargetNodes,
		"version":          i.KubernetesVersion,
		"applications":     i.Applications,
	}

	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/kubernetes/clusters/%s", id), params)
	if err != nil {
		return nil, err
	}

	kubernetes := &KubernetesCluster{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(kubernetes)
	return kubernetes, nil
}

// ListKubernetesMarketplaceApplications returns all application inside marketplace
func (c *Client) ListKubernetesMarketplaceApplications() ([]KubernetesMarketplaceApplication, error) {
	resp, err := c.SendGetRequest("/v2/kubernetes/applications")
	if err != nil {
		return nil, err
	}

	kubernetes := make([]KubernetesMarketplaceApplication, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&kubernetes); err != nil {
		return nil, err
	}

	return kubernetes, nil
}

// DeleteKubernetesCluster deletes a cluster
func (c *Client) DeleteKubernetesCluster(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/kubernetes/clusters/%s", id))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}

// RecycleKubernetesCluster create a new cluster of kubernetes
func (c *Client) RecycleKubernetesCluster(id string, hostname string) (*SimpleResponse, error) {
	body, err := c.SendPostRequest(fmt.Sprintf("/v2/kubernetes/clusters/%s", id), map[string]string{
		"hostname": hostname,
	})
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(body)
}

// ListAvailableKubernetesVersions returns all version of kubernetes available
func (c *Client) ListAvailableKubernetesVersions() ([]KubernetesVersion, error) {
	resp, err := c.SendGetRequest("/v2/kubernetes/versions")
	if err != nil {
		return nil, err
	}

	kubernetes := make([]KubernetesVersion, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&kubernetes); err != nil {
		return nil, err
	}

	return kubernetes, nil
}
