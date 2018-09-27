package cluster

import (
	"context"

	"yunion.io/x/yke/pkg/hosts"
	"yunion.io/x/yke/pkg/pki"
	"yunion.io/x/yke/pkg/services"
	"yunion.io/x/yke/pkg/types"
)

func (c *Cluster) ClusterRemove(ctx context.Context) error {
	externalEtcd := false
	if len(c.Services.Etcd.ExternalURLs) > 0 {
		externalEtcd = true
	}
	// Remove Worker Plane
	if err := services.RemoveWorkerPlane(ctx, c.WorkerHosts, true); err != nil {
		return err
	}

	// Remove Contol Plane
	if err := services.RemoveControlPlane(ctx, c.ControlPlaneHosts, true); err != nil {
		return err
	}

	// Remove Etcd Plane
	if err := services.RemoveEtcdPlane(ctx, c.EtcdHosts, true); err != nil {
		return err
	}

	// Clean up all hosts
	if err := cleanUpHosts(ctx, c.ControlPlaneHosts, c.WorkerHosts, c.EtcdHosts, c.SystemImages.Alpine, c.PrivateRegistriesMap, externalEtcd); err != nil {
		return err
	}

	pki.RemoveAdminConfig(ctx, c.LocalKubeConfigPath)
	return nil
}

func cleanUpHosts(ctx context.Context, cpHosts, workerHosts, etcdHosts []*hosts.Host, cleanerImage string, prsMap map[string]types.PrivateRegistry, externalEtcd bool) error {
	allHosts := []*hosts.Host{}
	allHosts = append(allHosts, cpHosts...)
	allHosts = append(allHosts, workerHosts...)
	allHosts = append(allHosts, etcdHosts...)

	for _, host := range allHosts {
		if err := host.CleanUpAll(ctx, cleanerImage, prsMap, externalEtcd); err != nil {
			return err
		}
	}
	return nil
}
