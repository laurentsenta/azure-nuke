package resources

import (
	"context"
	"github.com/ekristen/azure-nuke/pkg/nuke"
	"github.com/ekristen/cloud-nuke-sdk/pkg/resource"
	"github.com/ekristen/cloud-nuke-sdk/pkg/types"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
)

func init() {
	resource.Register(resource.Registration{
		Name:   "SecurityWorkspace",
		Scope:  nuke.Subscription,
		Lister: SecurityWorkspaceLister{},
	})
}

type SecurityWorkspace struct {
	client security.WorkspaceSettingsClient
	name   string
	scope  string
}

func (r *SecurityWorkspace) Remove() error {
	_, err := r.client.Delete(context.TODO(), r.name)
	return err
}

func (r *SecurityWorkspace) Properties() types.Properties {
	properties := types.NewProperties()

	properties.Set("Name", r.name)
	properties.Set("Scope", r.scope)

	return properties
}

func (r *SecurityWorkspace) String() string {
	return r.name
}

// -------------------------------------------------------------

type SecurityWorkspaceLister struct {
	opts nuke.ListerOpts
}

func (l SecurityWorkspaceLister) SetOptions(opts interface{}) {
	l.opts = opts.(nuke.ListerOpts)
}

func (l SecurityWorkspaceLister) List() ([]resource.Resource, error) {
	log := logrus.
		WithField("resource", "SecurityWorkspace").
		WithField("scope", nuke.Subscription).
		WithField("subscription", l.opts.SubscriptionId)

	log.Trace("creating client")

	client := security.NewWorkspaceSettingsClient(l.opts.SubscriptionId)
	client.Authorizer = l.opts.Authorizers.Management
	client.RetryAttempts = 1
	client.RetryDuration = time.Second * 2

	resources := make([]resource.Resource, 0)

	log.Trace("listing resources")

	ctx := context.TODO()
	list, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	for list.NotDone() {
		log.Trace("listing not done")
		for _, g := range list.Values() {
			resources = append(resources, &SecurityWorkspace{
				client: client,
				name:   *g.Name,
				scope:  *g.Scope,
			})
		}

		if err := list.NextWithContext(ctx); err != nil {
			return nil, err
		}
	}

	return resources, nil
}
