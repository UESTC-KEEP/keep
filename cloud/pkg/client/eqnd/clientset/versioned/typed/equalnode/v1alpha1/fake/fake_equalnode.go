/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/UESTC-KEEP/keep/cloud/pkg/apis/keepedge/equalnode/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeEqualNodes implements EqualNodeInterface
type FakeEqualNodes struct {
	Fake *FakeKeepedgeV1alpha1
	ns   string
}

var equalnodesResource = schema.GroupVersionResource{Group: "keepedge.k8s.io", Version: "v1alpha1", Resource: "equalnodes"}

var equalnodesKind = schema.GroupVersionKind{Group: "keepedge.k8s.io", Version: "v1alpha1", Kind: "EqualNode"}

// Get takes name of the equalNode, and returns the corresponding equalNode object, and an error if there is any.
func (c *FakeEqualNodes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.EqualNode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(equalnodesResource, c.ns, name), &v1alpha1.EqualNode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EqualNode), err
}

// List takes label and field selectors, and returns the list of EqualNodes that match those selectors.
func (c *FakeEqualNodes) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.EqualNodeList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(equalnodesResource, equalnodesKind, c.ns, opts), &v1alpha1.EqualNodeList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.EqualNodeList{ListMeta: obj.(*v1alpha1.EqualNodeList).ListMeta}
	for _, item := range obj.(*v1alpha1.EqualNodeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested equalNodes.
func (c *FakeEqualNodes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(equalnodesResource, c.ns, opts))

}

// Create takes the representation of a equalNode and creates it.  Returns the server's representation of the equalNode, and an error, if there is any.
func (c *FakeEqualNodes) Create(ctx context.Context, equalNode *v1alpha1.EqualNode, opts v1.CreateOptions) (result *v1alpha1.EqualNode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(equalnodesResource, c.ns, equalNode), &v1alpha1.EqualNode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EqualNode), err
}

// Update takes the representation of a equalNode and updates it. Returns the server's representation of the equalNode, and an error, if there is any.
func (c *FakeEqualNodes) Update(ctx context.Context, equalNode *v1alpha1.EqualNode, opts v1.UpdateOptions) (result *v1alpha1.EqualNode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(equalnodesResource, c.ns, equalNode), &v1alpha1.EqualNode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EqualNode), err
}

// Delete takes name of the equalNode and deletes it. Returns an error if one occurs.
func (c *FakeEqualNodes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(equalnodesResource, c.ns, name), &v1alpha1.EqualNode{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEqualNodes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(equalnodesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.EqualNodeList{})
	return err
}

// Patch applies the patch and returns the patched equalNode.
func (c *FakeEqualNodes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.EqualNode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(equalnodesResource, c.ns, name, pt, data, subresources...), &v1alpha1.EqualNode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EqualNode), err
}
