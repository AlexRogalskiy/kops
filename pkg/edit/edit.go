/*
Copyright 2017 The Kubernetes Authors.

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

package edit

import (
	"bytes"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/pkg/diff"
	"k8s.io/kops/pkg/kopscodecs"
	"k8s.io/kops/upup/pkg/fi/utils"
)

// HasExtraFields checks if the yaml has fields that were not mapped to the object
// (for example due to a typo in the field name)
// If there are extra fields it returns a string with a description of the diffs
// If there are no extra fields it returns an empty string
func HasExtraFields(yamlString string, object runtime.Object) (string, error) {
	switch object.(type) {
	case *kops.Cluster:
		editedObj := kops.Cluster{}
		err := yaml.UnmarshalStrict([]byte(yamlString), &editedObj)
		if err == nil {
			return "", nil
		}
	case *kops.InstanceGroup:
		editedObj := kops.InstanceGroup{}
		err := yaml.UnmarshalStrict([]byte(yamlString), &editedObj)
		if err == nil {
			return "", nil
		}
	default:
		panic("unknown object")
	}

	// Convert the cluster back to YAML for comparison purposes
	newYaml, err := kopscodecs.ToVersionedYaml(object)
	if err != nil {
		return "", err
	}

	// Marshal the edited YAML to a map; this will prevent bad diffs due to sorting
	var editedYamlObj map[string]interface{}
	err = yaml.UnmarshalStrict([]byte(yamlString), &editedYamlObj)
	if err != nil {
		return "", err
	}

	// Convert the object back to YAML so that we can compare it to the cluster YAML
	editedYaml, err := utils.YamlMarshal(editedYamlObj)
	if err != nil {
		return "", err
	}

	if !bytes.Equal(editedYaml, newYaml) {
		discardedChanges := diff.FormatDiff(string(newYaml), string(editedYaml))
		return discardedChanges, nil
	}

	return "", nil
}
