package resources

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-04-01/compute"

	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/azure-nuke/pkg/nuke"
)

const SSHPublicKeyResource = "SSHPublicKey"

func init() {
	resource.Register(resource.Registration{
		Name:   SSHPublicKeyResource,
		Lister: SSHPublicKeyLister{},
		Scope:  nuke.ResourceGroup,
	})
}

type SSHPublicKey struct {
	client compute.SSHPublicKeysClient
	name   *string
	rg     *string
}

func (r *SSHPublicKey) Remove() error {
	_, err := r.client.Delete(context.TODO(), *r.rg, *r.name)
	return err
}

func (r *SSHPublicKey) Properties() types.Properties {
	properties := types.NewProperties()

	properties.Set("Name", *r.name)

	return properties
}

func (r *SSHPublicKey) String() string {
	return *r.name
}

// --------------------------------------

type SSHPublicKeyLister struct {
}

func (l SSHPublicKeyLister) List(o interface{}) ([]resource.Resource, error) {
	opts := o.(nuke.ListerOpts)

	log := logrus.WithField("r", SSHPublicKeyResource).WithField("s", opts.SubscriptionId)

	client := compute.NewSSHPublicKeysClient(opts.SubscriptionId)
	client.Authorizer = opts.Authorizers.Management
	client.RetryAttempts = 1
	client.RetryDuration = time.Second * 2

	resources := make([]resource.Resource, 0)

	log.Trace("attempting to list ssh key")

	ctx := context.TODO()

	list, err := client.ListByResourceGroup(ctx, opts.ResourceGroup)
	if err != nil {
		return nil, err
	}

	log.Trace("listing ....")

	for list.NotDone() {
		log.Trace("list not done")
		for _, g := range list.Values() {
			resources = append(resources, &SSHPublicKey{
				client: client,
				name:   g.Name,
				rg:     &opts.ResourceGroup,
			})
		}

		if err := list.NextWithContext(ctx); err != nil {
			return nil, err
		}
	}

	log.Trace("done")

	return resources, nil
}
