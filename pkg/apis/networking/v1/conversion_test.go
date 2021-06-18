/*
Copyright 2015 The Kubernetes Authors.

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

package v1_test

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/networking/v1"
	networking "k8s.io/kubernetes/pkg/apis/networking"
	networkingv1 "k8s.io/kubernetes/pkg/apis/networking/v1"
)

func Test_networking_NetworkingPolicyPeer_to_v1_NetworkingPolicyPeer(t *testing.T) {
	// networking to v1
	testInputs := []networking.NetworkPolicyPeer{
		{
			// one IPBlock
			IPBlocks: []*networking.IPBlock{
				{
					CIDR: "192.168.1.0/24",
				},
			},
		},
		{
			// no IPBlock
			IPBlocks: nil,
		},
		{
			// list of IPBlocks
			IPBlocks: []*networking.IPBlock{
				{
					CIDR: "192.168.1.0/24",
				},
				{
					CIDR: "192.168.2.0/24",
				},
			},
		},
	}
	for i, input := range testInputs {
		v1NetworkingPolicyPeer := v1.NetworkPolicyPeer{}
		if err := networkingv1.Convert_networking_NetworkPolicyPeer_To_v1_NetworkPolicyPeer(&input, &v1NetworkingPolicyPeer, nil); nil != err {
			t.Errorf("%v: Convert networking.NetworkPolicyPeer to v1.NetworkPolicyPeer failed with error %v", i, err.Error())
		}

		if len(input.IPBlocks) == 0 {
			// no more work needed
			continue
		}
		// Primary IP was not set..
		if len(v1NetworkingPolicyPeer.IPBlock.CIDR) == 0 {
			t.Errorf("%v: Convert networking.NetworkPolicyPeer to v1.NetworkPolicyPeer failed out.IPBlock is empty, should be %v", i, v1NetworkingPolicyPeer.IPBlock)
		}

		// Primary should always == in.IPBlocks[0].IP
		if len(input.IPBlocks) > 0 && v1NetworkingPolicyPeer.IPBlock.CIDR != input.IPBlocks[0].CIDR {
			t.Errorf("%v: Convert networking.NetworkPolicyPeer to v1.NetworkPolicyPeer failed out.IPBlock != in.IPBlock[0] expected %v found %v", i, input.IPBlocks[0].CIDR, v1NetworkingPolicyPeer.IPBlock)
		}
		// match v1.IPBlocks to networking.IPBlocks
		for idx := range input.IPBlocks {
			if v1NetworkingPolicyPeer.IPBlocks[idx].CIDR != input.IPBlocks[idx].CIDR {
				t.Errorf("%v: Convert networking.NetworkPolicyPeer to v1.NetworkPolicyPeer failed. Expected  v1.NetworkPolicyPeer[%v]=%v but found %v", i, idx, input.IPBlocks[idx].CIDR, v1NetworkingPolicyPeer.IPBlocks[idx].CIDR)
			}
		}
	}
}
func Test_v1_NetworkingPolicyPeer_to_networking_NetworkingPolicyPeer(t *testing.T) {
	asymmetricInputs := []struct {
		name string
		in   v1.NetworkPolicyPeer
		out  networking.NetworkPolicyPeer
	}{
		{
			name: "mismatched IPBlock",
			in: v1.NetworkPolicyPeer{
				IPBlock: &v1.IPBlock{CIDR: "1.1.2.1"}, // Older field takes precedence for compatibility with patch by older clients
				IPBlocks: []*v1.IPBlock{
					{CIDR: "1.1.1.1"},
					{CIDR: "2.2.2.2"},
				},
			},
			out: networking.NetworkPolicyPeer{
				IPBlocks: []*networking.IPBlock{
					{CIDR: "1.1.2.1"},
				},
			},
		},
		{
			name: "matching IPBlock",
			in: v1.NetworkPolicyPeer{
				IPBlock: &v1.IPBlock{CIDR: "1.1.1.1"},
				IPBlocks: []*v1.IPBlock{
					{CIDR: "1.1.1.1"},
					{CIDR: "2.2.2.2"},
				},
			},
			out: networking.NetworkPolicyPeer{
				IPBlocks: []*networking.IPBlock{
					{CIDR: "1.1.1.1"},
					{CIDR: "2.2.2.2"},
				},
			},
		},
		{
			name: "empty IPBlock",
			in: v1.NetworkPolicyPeer{
				IPBlock: &v1.IPBlock{CIDR: ""},
				IPBlocks: []*v1.IPBlock{
					{CIDR: "1.1.1.1"},
					{CIDR: "2.2.2.2"},
				},
			},
			out: networking.NetworkPolicyPeer{
				IPBlocks: []*networking.IPBlock{
					{CIDR: "1.1.1.1"},
					{CIDR: "2.2.2.2"},
				},
			},
		},
	}

	// success
	v1TestInputs := []v1.NetworkPolicyPeer{
		// only Primary IP Provided
		{
			IPBlock: &v1.IPBlock{CIDR: "1.1.1.1"},
		},
		{
			// both are not provided
			IPBlock:  &v1.IPBlock{CIDR: ""},
			IPBlocks: nil,
		},
		// only list of IPs
		{
			IPBlocks: []*v1.IPBlock{
				{CIDR: "1.1.1.1"},
				{CIDR: "2.2.2.2"},
			},
		},
		// Both
		{
			IPBlock: &v1.IPBlock{CIDR: "1.1.1.1"},
			IPBlocks: []*v1.IPBlock{
				{CIDR: "1.1.1.1"},
				{CIDR: "2.2.2.2"},
			},
		},
		// v4 and v6
		{
			IPBlock: &v1.IPBlock{CIDR: "1.1.1.1"},
			IPBlocks: []*v1.IPBlock{
				{CIDR: "1.1.1.1"},
				{CIDR: "::1"},
			},
		},
		// v6 and v4
		{
			IPBlock: &v1.IPBlock{CIDR: "::1"},
			IPBlocks: []*v1.IPBlock{
				{CIDR: "::1"},
				{CIDR: "1.1.1.1"},
			},
		},
	}

	// run asymmetric cases
	for _, tc := range asymmetricInputs {
		testInput := tc.in

		networkingNetworkPolicyPeer := networking.NetworkPolicyPeer{}
		// convert..
		if err := networkingv1.Convert_v1_NetworkPolicyPeer_To_networking_NetworkPolicyPeer(&testInput, &networkingNetworkPolicyPeer, nil); err != nil {
			t.Errorf("%s: Convert v1.NetworkPolicyPeer to networking.NetworkPolicyPeer failed with error:%v for input %+v", tc.name, err.Error(), testInput)
		}
		if !reflect.DeepEqual(networkingNetworkPolicyPeer, tc.out) {
			t.Errorf("%s: expected %#v, got %#v", tc.name, tc.out.IPBlocks, networkingNetworkPolicyPeer.IPBlocks)
		}
	}

	// run ok cases
	for i, testInput := range v1TestInputs {
		networkingNetworkPolicyPeer := networking.NetworkPolicyPeer{}
		// convert..
		if err := networkingv1.Convert_v1_NetworkPolicyPeer_To_networking_NetworkPolicyPeer(&testInput, &networkingNetworkPolicyPeer, nil); err != nil {
			t.Errorf("%v: Convert v1.NetworkPolicyPeer to networking.NetworkPolicyPeer failed with error:%v for input %+v", i, err.Error(), testInput)
		}

		if len(testInput.IPBlock.CIDR) == 0 && len(testInput.IPBlocks) == 0 {
			continue //no more work needed
		}

		// List should have at least 1 IP == v1.IPBlock || v1.IPBlocks[0] (whichever provided)
		if len(testInput.IPBlock.CIDR) > 0 && networkingNetworkPolicyPeer.IPBlocks[0].CIDR != testInput.IPBlock.CIDR {
			t.Errorf("%v: Convert v1.NetworkPolicyPeer to networking.NetworkPolicyPeer failed. expected networkingNetworkPolicyPeer.IPBlocks[0].ip=%v found %v", i, networkingNetworkPolicyPeer.IPBlocks[0].CIDR, networkingNetworkPolicyPeer.IPBlocks[0].CIDR)
		}

		// walk the list
		for idx := range testInput.IPBlocks {
			if networkingNetworkPolicyPeer.IPBlocks[idx].CIDR != testInput.IPBlocks[idx].CIDR {
				t.Errorf("%v: Convert v1.NetworkPolicyPeer to networking.NetworkPolicyPeer failed networking.IPBlocks[%v]=%v expected %v", i, idx, networkingNetworkPolicyPeer.IPBlocks[idx].CIDR, testInput.IPBlocks[idx].CIDR)
			}
		}

		// if input has a list of IPs
		// then out put should have the same length
		if len(testInput.IPBlocks) > 0 && len(testInput.IPBlocks) != len(networkingNetworkPolicyPeer.IPBlocks) {
			t.Errorf("%v: Convert  v1.NetworkPolicyPeer to networking.NetworkPolicyPeer failed len(networkingNetworkPolicyPeer.IPBlocks) != len(v1.NetworkPolicyPeer.IPBlocks) [%v]=[%v]", i, len(networkingNetworkPolicyPeer.IPBlocks), len(testInput.IPBlocks))
		}
	}
}
