package resources

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage" //nolint:staticcheck

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/azure-nuke/pkg/nuke"
)

const StorageAccountResource = "StorageAccount"

func init() {
	registry.Register(&registry.Registration{
		Name:   StorageAccountResource,
		Scope:  nuke.ResourceGroup,
		Lister: &StorageAccountLister{},
		DependsOn: []string{
			VirtualMachineResource,
		},
	})
}

type StorageAccount struct {
	client storage.AccountsClient
	name   *string
	rg     string
	region *string
	tags   map[string]*string
}

func (r *StorageAccount) Remove(ctx context.Context) error {
	_, err := r.client.Delete(ctx, r.rg, *r.name)
	return err
}

func (r *StorageAccount) Properties() types.Properties {
	properties := types.NewProperties()

	properties.Set("Name", r.name)
	properties.Set("ResourceGroup", r.rg)
	properties.Set("Region", r.region)

	for k, v := range r.tags {
		properties.SetTag(&k, v)
	}

	return properties
}

func (r *StorageAccount) String() string {
	return *r.name
}

// --------------------------------------

type StorageAccountLister struct {
}

func (l StorageAccountLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	log := logrus.WithField("r", StorageAccountResource).WithField("s", opts.SubscriptionID)

	client := storage.NewAccountsClient(opts.SubscriptionID)
	client.Authorizer = opts.Authorizers.Management
	client.RetryAttempts = 1
	client.RetryDuration = time.Second * 2

	resources := make([]resource.Resource, 0)

	log.Trace("attempting to list ssh key")

	list, err := client.ListByResourceGroup(ctx, opts.ResourceGroup)
	if err != nil {
		return nil, err
	}

	log.Trace("listing storage accounts")

	for list.NotDone() {
		log.Trace("list not done")
		for _, g := range list.Values() {
			resources = append(resources, &StorageAccount{
				client: client,
				name:   g.Name,
				rg:     opts.ResourceGroup,
				region: g.Location,
			})
		}

		if err := list.NextWithContext(ctx); err != nil {
			return nil, err
		}
	}

	log.Trace("done")

	return resources, nil
}
