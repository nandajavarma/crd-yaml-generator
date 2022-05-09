package generator

import (
	"testing"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenerator(t *testing.T) {

	var i interface{}
	tests := []struct {
		name    string
		arg     apiext.CustomResourceDefinition
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Test function returns empty string, if input is empty",
			arg:     apiext.CustomResourceDefinition{},
			want:    i,
			wantErr: false,
		},
		{
			name: "Test function returns the YAML resource file",
			arg: apiext.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{},
				TypeMeta:   metav1.TypeMeta{},
				Spec:       apiext.CustomResourceDefinitionSpec{},
				Status:     apiext.CustomResourceDefinitionStatus{}},

			// 	`apiVersion: apiextensions.k8s.io/v1beta1
			// 		kind: CustomResourceDefinition
			// 		metadata:
			// 		name: stars.example.crd.com
			// 		spec:
			// 		group: example.crd.com
			// 		scope: Namespaced
			// 		names:
			// 			kind: Star
			// 			listKind: StarList
			// 			plural: stars
			// 			singular: star
			// 		subresources:
			// 			status: {}
			// 		validation:
			// 			openAPIV3Schema:
			// 			required: ["spec"]
			// 			properties:
			// 				spec:
			// 				required: ["type","location"]
			// 				properties:
			// 					type:
			// 					type: "string"
			// 					minimum: 1
			// 					location:
			// 					type: "string"
			// 					minimum: 1
			// 		versions:
			// 			- name: v1alpha1
			// 			served: true
			// 			storage: true`,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		got, err := generateYaml(tt.arg)
		if (err != nil) != tt.wantErr {
			t.Errorf("`%s` failed: got error %s", tt.name, err)
		}

		if got != tt.want {
			t.Errorf("`%s` failed: got %v want %v", tt.name, got, tt.want)
		}
	}

}
