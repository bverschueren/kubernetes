/*
Copyright 2020 The Kubernetes Authors.

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

package v1

import (
	v1 "k8s.io/api/networking/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	networking "k8s.io/kubernetes/pkg/apis/networking"
)

func Convert_v1_NetworkPolicyPeer_To_networking_NetworkPolicyPeer(in *v1.NetworkPolicyPeer, out *networking.NetworkPolicyPeer, s conversion.Scope) error {
	if err := autoConvert_v1_NetworkPolicyPeer_To_networking_NetworkPolicyPeer(in, out, s); err != nil {
		return err
	}

	if (len(in.IPBlock.CIDR) > 0 && len(in.IPBlocks) > 0) && (in.IPBlock != in.IPBlocks[0]) {
		out.IPBlocks = []*networking.IPBlock{
			{
				CIDR: in.IPBlock.CIDR,
			},
		}
	}
	// at the this point, autoConvert copied v1.IPBlocks -> networking.IPBlocks
	// if v1.IPBlocks was empty but v1.IPBlock is not, then set networking.IPBlocks[0] with v1.IPBlock
	if len(in.IPBlock.CIDR) > 0 && len(in.IPBlocks) == 0 {
		out.IPBlocks = []*networking.IPBlock{
			{
				CIDR: in.IPBlock.CIDR,
			},
		}
	}
	return nil
}

func Convert_networking_NetworkPolicyPeer_To_v1_NetworkPolicyPeer(in *networking.NetworkPolicyPeer, out *v1.NetworkPolicyPeer, s conversion.Scope) error {
	if err := autoConvert_networking_NetworkPolicyPeer_To_v1_NetworkPolicyPeer(in, out, s); err != nil {
		return err
	}
	// at the this point autoConvert copied networking.IPBlocks -> v1.IPBlocks
	//  v1.IPBlock (singular value field, which does not exist in networking) needs to
	// be set with networking.IPBlocks[0]
	if len(in.IPBlocks) > 0 {
		out.IPBlock.CIDR = in.IPBlocks[0].CIDR
	}
	return nil
}
