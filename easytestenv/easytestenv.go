package easytestenv

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

type EasyTestEnv struct {
	Env           *envtest.Environment
	RestCfg       *rest.Config
	Client        client.Client
	DynamicClient *dynamic.DynamicClient

	Context context.Context
	Cancel  context.CancelFunc
}

type Config struct {
	CRDDirectoryPaths     []string
	ErrorIfCRDPathMissing bool
}

func New(config Config) (result *EasyTestEnv, err error) {
	result = &EasyTestEnv{}
	result.Context, result.Cancel = context.WithCancel(context.TODO())

	result.Env = &envtest.Environment{
		CRDDirectoryPaths:     config.CRDDirectoryPaths,
		ErrorIfCRDPathMissing: config.ErrorIfCRDPathMissing,
	}

	result.RestCfg, err = result.Env.Start()
	if err != nil {
		return result, err
	}

	result.Client, err = client.New(result.RestCfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return result, err
	}

	result.DynamicClient, err = dynamic.NewForConfig(result.RestCfg)
	if err != nil {
		return result, err
	}

	return result, err
}

func (i *EasyTestEnv) Stop() {
	i.Cancel()
}

// Create all resources at the given path, recursively if the path is a directory
func (i *EasyTestEnv) ApplyResources(path string) error {
	var prioritizedResources [3][]client.Object
	err := gatherResourcesAtPath(path, true, prioritizedResources)
	if err != nil {
		return err
	}

	for _, resources := range prioritizedResources {
		for _, resource := range resources {
			err = i.Client.Create(i.Context, resource)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func gatherResourcesAtPath(path string, recursive bool, result [3][]client.Object) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		err := filepath.WalkDir(path, func(subpath string, dirEntry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !dirEntry.IsDir() || recursive {
				err := gatherResourcesAtPath(subpath, recursive, result)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
	} else {
		resource, err := loadResourceAtPath(path)
		if err != nil {
			return err
		}
		gvk := resource.GroupVersionKind()
		if !gvk.Empty() {
			priority := determinePriorityByKind(gvk)
			result[priority] = append(result[priority], resource)
		}
	}
	return nil
}

func loadResourceAtPath(filepath string) (*unstructured.Unstructured, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	resource := &unstructured.Unstructured{}
	err = yaml.Unmarshal(contents, resource)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// Create these resources first
var priority0 = []schema.GroupVersionKind{
	{
		Version: "v1",
		Kind:    "Namespace",
	},
	{
		Version: "v1",
		Kind:    "CustomResourceDefinition",
	},
}

// Create these resources second
var priority1 = []schema.GroupVersionKind{
	{
		Group:   "admissionregistration.k8s.io",
		Version: "v1",
		Kind:    "MutatingWebhookConfiguration",
	},
	{
		Group:   "admissionregistration.k8s.io",
		Version: "v1",
		Kind:    "ValidatingWebhookConfiguration",
	},
}

// Determines the priority order of a resource when creating multiple resources
func determinePriorityByKind(gvk schema.GroupVersionKind) uint {
	if slices.IndexFunc(priority0, func(haystack schema.GroupVersionKind) bool {
		return gvk.Group == haystack.Group && gvk.Version == haystack.Version && gvk.Kind == haystack.Kind
	}) > -1 {
		return 0
	}

	if slices.IndexFunc(priority1, func(haystack schema.GroupVersionKind) bool {
		return gvk.Group == haystack.Group && gvk.Version == haystack.Version && gvk.Kind == haystack.Kind
	}) > -1 {
		return 1
	}

	return 2
}
