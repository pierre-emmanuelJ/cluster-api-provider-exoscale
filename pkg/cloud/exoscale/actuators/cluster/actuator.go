/*
Copyright 2019 The Kubernetes authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/exoscale/egoscale"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/klog"
	exoscaleconfigv1 "sigs.k8s.io/cluster-api-provider-exoscale/pkg/apis/exoscale/v1alpha1"
	exoclient "sigs.k8s.io/cluster-api-provider-exoscale/pkg/cloud/exoscale/client"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
)

const (
	ExoscaleIPAnnotationKey = "exoscale-ip"
)

// Actuator is responsible for performing cluster reconciliation
type Actuator struct {
	clustersGetter client.ClustersGetter
}

// ActuatorParams holds parameter information for Actuator
type ActuatorParams struct {
	ClustersGetter client.ClustersGetter
}

// NewActuator creates a new Actuator
func NewActuator(params ActuatorParams) (*Actuator, error) {
	return &Actuator{
		clustersGetter: params.ClustersGetter,
	}, nil
}

// Reconcile reconciles a cluster and is invoked by the Cluster Controller
func (a *Actuator) Reconcile(cluster *clusterv1.Cluster) error {
	log.Printf("Reconciling cluster %v.", cluster.Name)

	clusterConfig, err := clusterProviderFromProviderConfig(cluster.Spec.ProviderSpec)
	if err != nil {
		return fmt.Errorf("error loading cluster provider config: %v", err)
	}

	exoClient, err := exoclient.Client()
	if err != nil {
		return err
	}

	req := egoscale.CreateSecurityGroup{
		Name: clusterConfig.Spec.SecurityGroup,
	}

	_, err = exoClient.RequestWithContext(context.TODO(), req)
	if err != nil {
		return fmt.Errorf("error creating or updating network security group: %v", err)
	}

	return nil
}

// Delete deletes a cluster and is invoked by the Cluster Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster) error {
	log.Printf("Deleting cluster %v.", cluster.Name)
	return fmt.Errorf("TODO: Not yet implemented")
}

// The Machine Actuator interface must implement GetIP and GetKubeConfig functions as a workaround for issues
// cluster-api#158 (https://github.com/kubernetes-sigs/cluster-api/issues/158) and cluster-api#160
// (https://github.com/kubernetes-sigs/cluster-api/issues/160).

// GetIP returns IP address of the machine in the cluster.
func (*Actuator) GetIP(cluster *clusterv1.Cluster, machine *clusterv1.Machine) (string, error) {
	log.Printf("Getting IP of machine %v for cluster %v.", machine.Name, cluster.Name)
	if machine.ObjectMeta.Annotations != nil {
		if ip, ok := machine.ObjectMeta.Annotations[ExoscaleIPAnnotationKey]; ok {
			klog.Infof("Returning IP from machine annotation %s", ip)
			return ip, nil
		}
	}

	return "", errors.New("could not get IP")
}

// GetKubeConfig gets a kubeconfig from the master.
func (*Actuator) GetKubeConfig(cluster *clusterv1.Cluster, master *clusterv1.Machine) (string, error) {
	log.Printf("Getting IP of machine %v for cluster %v.", master.Name, cluster.Name)
	return "", fmt.Errorf("Provisionner exoscale GetKubeConfig() not yet implemented")
}

func clusterProviderFromProviderConfig(providerConfig clusterv1.ProviderSpec) (*exoscaleconfigv1.ExoscaleClusterProviderConfig, error) {
	var config exoscaleconfigv1.ExoscaleClusterProviderConfig
	if err := yaml.Unmarshal(providerConfig.Value.Raw, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
