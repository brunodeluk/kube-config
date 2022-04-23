package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/pkg/apis/clientauthentication/v1alpha1"
	"k8s.io/client-go/rest"
	"os"
	"path/filepath"
)

type Kubernetes struct {
}

func (k *Kubernetes) Apply(ctx context.Context, path string) error {
	objects, err := filesToObjects(path)
	fmt.Printf("read %d objects", len(objects))

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.APIPath = "/apis"
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	config.ContentConfig.GroupVersion = &schema.GroupVersion{
		Group:   v1alpha1.GroupName,
		Version: v1alpha1.SchemeGroupVersion.Version,
	}

	client, err := rest.RESTClientFor(config)
	if err != nil {
		return err
	}

	fmt.Println("starting applying process")
	for i, obj := range objects {
		fmt.Printf("Applying manifest %d...\n", i+1)
		applyObject := obj.DeepCopy()
		b, err := applyObject.MarshalJSON()
		if err != nil {
			return err
		}

		req := client.
			Patch(types.ApplyPatchType).
			Body(b).
			Do(ctx)

		if req.Error() != nil {
			return err
		}
	}

	return err
}

func filesToObjects(path string) ([]*unstructured.Unstructured, error) {
	info, err := os.Stat(path)
	if err != nil {
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
				return err
			}

			for _, obj := range objs {
				objects = append(objects, obj)
			}

			return err
		})
		if err != nil {
			return nil, err
		}

		return objects, nil
	}

	objects, err = readObjects(path)
	if err != nil {
		return nil, err
	}
	return objects, nil
}

func readObjects(yamlFile string) ([]*unstructured.Unstructured, error) {
	file, err := os.Open(yamlFile)
	if err != nil {
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
