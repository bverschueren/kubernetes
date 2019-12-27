/*
Copyright 2019 The Kubernetes Authors.

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

package genericclioptions

import (
	"encoding/json"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"strings"
	"testing"
)

func TestRecordFlags(t *testing.T) {
	tests := []struct {
		name          string
		object        runtime.Object
		record        bool
		update        bool
		changeCommand []string
	}{
		{
			name: "record ChangeCauseAnnotation on appsv1/Deployment change with existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "myobject",
					Annotations: map[string]string{ChangeCauseAnnotation: "create_cmd some_argument --record=true"},
				},
				Spec: appsv1.DeploymentSpec{},
			},
			record:        true,
			update:        true,
			changeCommand: []string{"change_cmd", "some_argument", "--record=true"},
		},
		{
			name: "record ChangeCauseAnnotation on appsv1/Deployment change without existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: "myobject"},
				Spec:       appsv1.DeploymentSpec{},
			},
			record:        true,
			update:        true,
			changeCommand: []string{"change_cmd", "some_argument", "--record=true"},
		},
		{
			name: "update ChangeCauseAnnotation on appsv1/Deployment change with existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "myobject",
					Annotations: map[string]string{ChangeCauseAnnotation: "create_cmd some_argument --record=true"},
				},
				Spec: appsv1.DeploymentSpec{},
			},
			record:        false,
			update:        true,
			changeCommand: []string{"change_cmd", "some_argument"},
		},
		{
			name: "update ChangeCauseAnnotation on appsv1/Deployment change without existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "myobject",
				},
				Spec: appsv1.DeploymentSpec{},
			},
			record:        false,
			update:        true,
			changeCommand: []string{"change_cmd", "some_argument"},
		},
		{
			name: "do not record ChangeCauseAnnotation on appsv1/Deployment change with existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "myobject",
					Annotations: map[string]string{ChangeCauseAnnotation: "create_cmd some_argument --record=true"},
				},
				Spec: appsv1.DeploymentSpec{},
			},
			record:        false,
			update:        false,
			changeCommand: []string{"change_cmd", "some_argument", "--record=false"},
		},
		{
			name: "do not record ChangeCauseAnnotation on appsv1/Deployment change without existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "myobject",
				},
				Spec: appsv1.DeploymentSpec{},
			},
			record:        false,
			update:        false,
			changeCommand: []string{"change_cmd", "some_argument", "--record=false"},
		},
	}
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expected := strings.Join(tc.changeCommand, " ")
			// mock provided command
			os.Args = tc.changeCommand

			accessor, err := meta.Accessor(tc.object)
			if err != nil {
				t.Fatal(err)
			}
			annotations := accessor.GetAnnotations()
			// keep original change-cause for comparison
			original := annotations[ChangeCauseAnnotation]

			rf := &RecordFlags{
				Record: &tc.record,
				Update: &tc.update,
			}
			cmd := &cobra.Command{}
			rf.AddFlags(cmd)
			rf.Complete(cmd)
			rec, err := rf.ToRecorder()
			if err != nil {
				t.Fatal(err)
			}
			rec.Record(tc.object)

			annotations = accessor.GetAnnotations()
			actual := annotations[ChangeCauseAnnotation]
			// verify annotation is recorded:
			if tc.record || (tc.update && original != "") {
				if expected != actual {
					t.Errorf("%v: expected '%v', got '%v'", tc.name, expected, actual)
				}
			} else {
				if original != actual {
					t.Errorf("%v: expected '%v', got '%v'", tc.name, original, actual)
				}
			}
		})
	}
}

func TestMakeRecordMergePatch(t *testing.T) {
	tests := []struct {
		name          string
		object        runtime.Object
		record        bool
		update        bool
		changeCommand []string
	}{
		{
			name: "record ChangeCauseAnnotation on appsv1/Deployment change with existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "myobject",
					Annotations: map[string]string{ChangeCauseAnnotation: "create_cmd some_argument --record=true"},
				},
				Spec: appsv1.DeploymentSpec{},
			},
			record:        true,
			update:        true,
			changeCommand: []string{"change_cmd", "some_argument", "--record=true"},
		},
		{
			name: "record ChangeCauseAnnotation on appsv1/Deployment change without existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: "myobject"},
				Spec:       appsv1.DeploymentSpec{},
			},
			record:        true,
			update:        true,
			changeCommand: []string{"change_cmd", "some_argument", "--record=true"},
		},
		{
			name: "update ChangeCauseAnnotation on appsv1/Deployment change with existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "myobject",
					Annotations: map[string]string{ChangeCauseAnnotation: "create_cmd some_argument --record=true"},
				},
				Spec: appsv1.DeploymentSpec{},
			},
			record:        false,
			update:        true,
			changeCommand: []string{"change_cmd", "some_argument"},
		},
		{
			name: "update ChangeCauseAnnotation on appsv1/Deployment change without existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "myobject",
				},
				Spec: appsv1.DeploymentSpec{},
			},
			record:        false,
			update:        true,
			changeCommand: []string{"change_cmd", "some_argument"},
		},
		{
			name: "do not record ChangeCauseAnnotation on appsv1/Deployment change with existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "myobject",
					Annotations: map[string]string{ChangeCauseAnnotation: "create_cmd some_argument --record=true"},
				},
				Spec: appsv1.DeploymentSpec{},
			},
			record:        false,
			update:        false,
			changeCommand: []string{"change_cmd", "some_argument", "--record=false"},
		},
		{
			name: "do not record ChangeCauseAnnotation on appsv1/Deployment change without existing ChangeCauseAnnotation",
			object: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "myobject",
				},
				Spec: appsv1.DeploymentSpec{},
			},
			record:        false,
			update:        false,
			changeCommand: []string{"change_cmd", "some_argument", "--record=false"},
		},
	}
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expected := strings.Join(tc.changeCommand, " ")
			// mock provided command
			os.Args = tc.changeCommand

			accessor, err := meta.Accessor(tc.object)
			if err != nil {
				t.Fatal(err)
			}
			annotations := accessor.GetAnnotations()
			// keep original change-cause for comparison
			original := annotations[ChangeCauseAnnotation]

			rf := &RecordFlags{
				Record: &tc.record,
				Update: &tc.update,
			}
			cmd := &cobra.Command{}
			rf.AddFlags(cmd)
			rf.Complete(cmd)
			rec, err := rf.ToRecorder()
			if err != nil {
				t.Fatal(err)
			}
			patch, _ := rec.MakeRecordMergePatch(tc.object)

			patchObj := appsv1.Deployment{}
			if err := json.Unmarshal(patch, &patchObj); err != nil {
				patchObj = appsv1.Deployment{}
			}

			annotations = patchObj.GetAnnotations()
			actual := annotations[ChangeCauseAnnotation]
			// verify annotation is in the patch
			if tc.record || (tc.update && original != "") {
				if expected != actual {
					t.Errorf("%v: expected '%v', got '%v'", tc.name, expected, actual)
				}
			} else {
				if actual != "" {
					t.Errorf("%v: expected '', got '%v'", tc.name, actual)
				}
			}
		})
	}
}
