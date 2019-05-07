/*
Portions Copyright 2017 The Kubernetes Authors.
Portions Copyright 2018 Aspen Mesh Authors.
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

package v1alpha3

import (
	"bufio"
	"bytes"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"

	istiov1alpha3 "istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	log "github.com/sirupsen/logrus"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DestinationRule is an Istio DestinationRule resource
type DestinationRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DestinationRuleSpec `json:"spec"`
}

func (dr *DestinationRule) GetSpecMessage() proto.Message {
	return &dr.Spec.DestinationRule
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DestinationRuleList is a list of DestinationRule resources
type DestinationRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []DestinationRule `json:"items"`
}

// DestinationRuleSpec is a wrapper around Istio DestinationRule
type DestinationRuleSpec struct {
	istiov1alpha3.DestinationRule
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
// Based of https://github.com/istio/istio/blob/release-0.8/pilot/pkg/config/kube/crd/types.go#L450
func (drs *DestinationRuleSpec) DeepCopyInto(out *DestinationRuleSpec) {
	*out = *drs
}

func (drs *DestinationRuleSpec) MarshalJSON() ([]byte, error) {
	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)
	marshaler := jsonpb.Marshaler{}
	err := marshaler.Marshal(writer, &drs.DestinationRule)
	if err != nil {
		log.WithField("error", err).Error("Could not marshal DestinationRuleSpec")
		return nil, err
	}

	writer.Flush()
	return buffer.Bytes(), nil
}

func (drs *DestinationRuleSpec) UnmarshalJSON(b []byte) error {
	reader := bytes.NewReader(b)
	unmarshaler := jsonpb.Unmarshaler{}
	err := unmarshaler.Unmarshal(reader, &drs.DestinationRule)
	if err != nil {
		log.WithField("error", err).Error("Could not unmarshal DestinationRuleSpec")
		return err
	}
	return nil
}
