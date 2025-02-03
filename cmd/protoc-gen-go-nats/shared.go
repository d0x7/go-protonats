package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	protonats "xiam.li/protonats/go"
)

func isUsingBroadcasting(method *protogen.Method) bool {
	if !proto.HasExtension(method.Desc.Options(), protonats.E_Broadcast) {
		return false
	}
	extension := proto.GetExtension(method.Desc.Options(), protonats.E_Broadcast)
	return extension.(bool)
}

func getConsensusTarget(method *protogen.Method) *protonats.ConsensusTarget {
	if !proto.HasExtension(method.Desc.Options(), protonats.E_ConsensusTarget) {
		return nil
	}
	extension := proto.GetExtension(method.Desc.Options(), protonats.E_ConsensusTarget).(protonats.ConsensusTarget)
	return &extension
}

func isConsensusLeader(method *protogen.Method) bool {
	if target := getConsensusTarget(method); target == nil {
		return false
	} else {
		return *target == protonats.ConsensusTarget_LEADER
	}
}

func isConsensusFollower(method *protogen.Method) bool {
	if target := getConsensusTarget(method); target == nil {
		return false
	} else {
		return *target == protonats.ConsensusTarget_FOLLOWER
	}
}
