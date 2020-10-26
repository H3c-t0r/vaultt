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

package v1beta1

import (
	"context"
	"time"

	v1beta1 "k8s.io/api/storage/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

// CSINodesGetter has a method to return a CSINodeInterface.
// A group's client should implement this interface.
type CSINodesGetter interface {
	CSINodes() CSINodeInterface
}

// CSINodeInterface has methods to work with CSINode resources.
type CSINodeInterface interface {
	Create(ctx context.Context, cSINode *v1beta1.CSINode, opts v1.CreateOptions) (*v1beta1.CSINode, error)
	Update(ctx context.Context, cSINode *v1beta1.CSINode, opts v1.UpdateOptions) (*v1beta1.CSINode, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.CSINode, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.CSINodeList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.CSINode, err error)
	CSINodeExpansion
}

// cSINodes implements CSINodeInterface
type cSINodes struct {
	client rest.Interface
}

// newCSINodes returns a CSINodes
func newCSINodes(c *StorageV1beta1Client) *cSINodes {
	return &cSINodes{
		client: c.RESTClient(),
	}
}

// Get takes name of the cSINode, and returns the corresponding cSINode object, and an error if there is any.
func (c *cSINodes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.CSINode, err error) {
	result = &v1beta1.CSINode{}
	err = c.client.Get().
		Resource("csinodes").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CSINodes that match those selectors.
func (c *cSINodes) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.CSINodeList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.CSINodeList{}
	err = c.client.Get().
		Resource("csinodes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested cSINodes.
func (c *cSINodes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("csinodes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a cSINode and creates it.  Returns the server's representation of the cSINode, and an error, if there is any.
func (c *cSINodes) Create(ctx context.Context, cSINode *v1beta1.CSINode, opts v1.CreateOptions) (result *v1beta1.CSINode, err error) {
	result = &v1beta1.CSINode{}
	err = c.client.Post().
		Resource("csinodes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(cSINode).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a cSINode and updates it. Returns the server's representation of the cSINode, and an error, if there is any.
func (c *cSINodes) Update(ctx context.Context, cSINode *v1beta1.CSINode, opts v1.UpdateOptions) (result *v1beta1.CSINode, err error) {
	result = &v1beta1.CSINode{}
	err = c.client.Put().
		Resource("csinodes").
		Name(cSINode.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(cSINode).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the cSINode and deletes it. Returns an error if one occurs.
func (c *cSINodes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("csinodes").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *cSINodes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("csinodes").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched cSINode.
func (c *cSINodes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.CSINode, err error) {
	result = &v1beta1.CSINode{}
	err = c.client.Patch(pt).
		Resource("csinodes").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
