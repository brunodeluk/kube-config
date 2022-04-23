package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"os"
	"path/filepath"
)

type Kubernetes struct {
}

func (k *Kubernetes) Apply(ctx context.Context, path string) error {
	objects, err := filesToObjects(path)
	fmt.Printf("[client][kube-client][INFO] read %d objects\n", len(objects))
	fmt.Printf("[client][kube-client][INFO] configuring rest client...\n")

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("[client][kube-client][ERROR] configuring rest client\n")
		return err
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("[client][kube-client][ERROR] configuring creating rest client\n")
		return err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	groupResource, err := restmapper.GetAPIGroupResources(clientSet)
	if err != nil {
		return err
	}

	rm := restmapper.NewDiscoveryRESTMapper(groupResource)

	fmt.Printf("[client][kube-client][INFO] applying objects...\n")
	for _, obj := range objects {
		applyObject := obj.DeepCopy()

		mapping, _ := rm.RESTMapping(obj.GroupVersionKind().GroupKind(), obj.GroupVersionKind().Version)

		res, err := client.
			Resource(mapping.Resource).
			Namespace(obj.GetNamespace()).
			Create(ctx, applyObject, v1.CreateOptions{})

		if err != nil {
			fmt.Printf("[client][kube-client][ERROR] creating object %v: %s\n", obj.GetKind(), err.Error())
			return err
		}

		fmt.Printf("[client][kube-client][INFO] created %s\n", res.GetName())
	}

	return err
}

func filesToObjects(path string) ([]*unstructured.Unstructured, error) {
	fmt.Printf("[client][kube-client][INFO] transforming yaml files to objects\n")
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("[client][kube-client][ERROR] error with os.Stat\n")
		return nil, err
	}

	objects := make([]*unstructured.Unstructured, 0)

	if info.IsDir() {
		err = filepath.Walk(path, func(filepath string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			objs, err := readObjects(filepath)
			if err != nil {
				fmt.Printf("[client][kube-client][ERROR] reading object\n")
				return err
			}

			for _, obj := range objs {
				objects = append(objects, obj)
			}

			return err
		})
		if err != nil {
			fmt.Printf("[client][kube-client][ERROR] walking through dir\n")
			return nil, err
		}

		return objects, nil
	}

	objects, err = readObjects(path)
	if err != nil {
		fmt.Printf("[client][kube-client][ERROR] reading single object\n")
		return nil, err
	}
	return objects, nil
}

func readObjects(yamlFile string) ([]*unstructured.Unstructured, error) {
	fmt.Printf("[client][kube-client][INFO] reading yaml file %s\n", yamlFile)
	file, err := os.Open(yamlFile)
	if err != nil {
		fmt.Printf("[client][kube-client][ERROR] error opening yaml file %s\n", yamlFile)
		return nil, err
	}

	defer file.Close()

	reader := yaml.NewYAMLOrJSONDecoder(bufio.NewReader(file), 2048)
	objects := make([]*unstructured.Unstructured, 0)

	for {
		obj := &unstructured.Unstructured{}
		err := reader.Decode(obj)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return objects, err
		}

		if obj.IsList() {
			err = obj.EachListItem(func(item runtime.Object) error {
				obj := item.(*unstructured.Unstructured)
				objects = append(objects, obj)
				return nil
			})
			if err != nil {
				fmt.Printf("[client][kube-client][ERROR] applying func to item\n")
				return objects, err
			}
			continue
		}

		if isKubernetesObject(obj) {
			objects = append(objects, obj)
		}
	}

	return objects, nil
}

func isKubernetesObject(obj *unstructured.Unstructured) bool {
	if obj.GetName() == "" || obj.GetKind() == "" || obj.GetAPIVersion() == "" {
		return false
	}
	return true
}
