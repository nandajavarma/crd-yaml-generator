package generator

import (
	"fmt"

	"gopkg.in/yaml.v2"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type crdMetadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace,omitempty"`
}

type crdYAML struct {
	ApiVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   crdMetadata `yaml:"metadata"`
	Spec       interface{}
}

func NewCrdYAML() crdYAML {
	return crdYAML{}
}

func GenerateYaml(crdDef extv1.CustomResourceDefinition) ([]byte, error) {
	crdyaml := NewCrdYAML()
	group := crdDef.Spec.Group
	isNamespaced := crdDef.Spec.Scope == "Namespaced"
	versions := crdDef.Spec.Versions

	// FIXME
	version := versions[0].Name

	crdyaml.Kind = crdDef.Spec.Names.Kind
	crdyaml.ApiVersion = fmt.Sprintf("%s/%s", group, version)

	crdyaml.Metadata = crdMetadata{}
	crdyaml.Metadata.Name = "Auto-generated YAML"
	if isNamespaced {
		crdyaml.Metadata.Namespace = "default"
	}

	result, err := yaml.Marshal(crdyaml)
	if err != nil {
		return nil, err
	}

	return result, nil
}
