package utils

import (
	"github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/hyperledger/fabric/common/policies"
	"github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
)

const (
	// ConsensusTypeSolo identifies the solo consensus implementation.
	ConsensusTypeSolo = "solo"
	// ConsensusTypeKafka identifies the Kafka-based consensus implementation.
	ConsensusTypeKafka = "kafka"

	// BlockValidationPolicyKey TODO
	BlockValidationPolicyKey = "BlockValidation"

	// OrdererAdminsPolicy is the absolute path to the orderer admins policy
	OrdererAdminsPolicy = "/Channel/Orderer/Admins"

	// SignaturePolicyType is the 'Type' string for signature policies
	SignaturePolicyType = "Signature"

	// ImplicitMetaPolicyType is the 'Type' string for implicit meta policies
	ImplicitMetaPolicyType = "ImplicitMeta"
)

const (
	// TestChainID is the channel name used for testing purposes when one is
	// not given
	TestChainID = "testchainid"

	// SampleInsecureSoloProfile references the sample profile which does not
	// include any MSPs and uses solo for ordering.
	SampleInsecureSoloProfile = "SampleInsecureSolo"
	// SampleDevModeSoloProfile references the sample profile which requires
	// only basic membership for admin privileges and uses solo for ordering.
	SampleDevModeSoloProfile = "SampleDevModeSolo"
	// SampleSingleMSPSoloProfile references the sample profile which includes
	// only the sample MSP and uses solo for ordering.
	SampleSingleMSPSoloProfile = "SampleSingleMSPSolo"

	// SampleInsecureKafkaProfile references the sample profile which does not
	// include any MSPs and uses Kafka for ordering.
	SampleInsecureKafkaProfile = "SampleInsecureKafka"
	// SampleDevModeKafkaProfile references the sample profile which requires only
	// basic membership for admin privileges and uses Kafka for ordering.
	SampleDevModeKafkaProfile = "SampleDevModeKafka"
	// SampleSingleMSPKafkaProfile references the sample profile which includes
	// only the sample MSP and uses Kafka for ordering.
	SampleSingleMSPKafkaProfile = "SampleSingleMSPKafka"

	// SampleDevModeEtcdRaftProfile references the sample profile used for testing
	// the etcd/raft-based ordering service.
	SampleDevModeEtcdRaftProfile = "SampleDevModeEtcdRaft"

	// SampleSingleMSPChannelProfile references the sample profile which
	// includes only the sample MSP and is used to create a channel
	SampleSingleMSPChannelProfile = "SampleSingleMSPChannel"

	// SampleConsortiumName is the sample consortium from the
	// sample configtx.yaml
	SampleConsortiumName = "SampleConsortium"
	// SampleOrgName is the name of the sample org in the sample profiles
	SampleOrgName = "SampleOrg"

	// AdminRoleAdminPrincipal is set as AdminRole to cause the MSP role of
	// type Admin to be used as the admin principal default
	AdminRoleAdminPrincipal = "Role.ADMIN"
	// MemberRoleAdminPrincipal is set as AdminRole to cause the MSP role of
	// type Member to be used as the admin principal default
	MemberRoleAdminPrincipal = "Role.MEMBER"
)

func NewOrgGroup(conf *localconfig.Organization)(*common.ConfigGroup,error){
	mspConfig, err := msp.GetVerifyingMspConfig(conf.MSPDir, conf.ID, conf.MSPType)
	if err != nil {
		return nil, errors.Wrapf(err, "1 - Error loading MSP configuration for org: %s", conf.Name)
	}

	ordererOrgGroup := common.NewConfigGroup()
	if len(conf.Policies) == 0 {
		addSignaturePolicyDefaults(ordererOrgGroup, conf.ID, conf.AdminPrincipal != AdminRoleAdminPrincipal)
	} else {
		if err := addPolicies(ordererOrgGroup, conf.Policies, channelconfig.AdminsPolicyKey); err != nil {
			return nil, errors.Wrapf(err, "error adding policies to orderer org group '%s'", conf.Name)
		}
	}

	addValue(ordererOrgGroup, channelconfig.MSPValue(mspConfig), channelconfig.AdminsPolicyKey)

	ordererOrgGroup.ModPolicy = channelconfig.AdminsPolicyKey
	return ordererOrgGroup, nil
}


func addValue(cg *common.ConfigGroup, value channelconfig.ConfigValue, modPolicy string) {
	cg.Values[value.Key()] = &common.ConfigValue{
		Value:     utils.MarshalOrPanic(value.Value()),
		ModPolicy: modPolicy,
	}
}

func addPolicy(cg *common.ConfigGroup, policy policies.ConfigPolicy, modPolicy string) {
	cg.Policies[policy.Key()] = &common.ConfigPolicy{
		Policy:    policy.Value(),
		ModPolicy: modPolicy,
	}
}

func addPolicies(cg *common.ConfigGroup, policyMap map[string]*localconfig.Policy, modPolicy string) error {
	for policyName, policy := range policyMap {
		switch policy.Type {
		case ImplicitMetaPolicyType:
			imp, err := policies.ImplicitMetaFromString(policy.Rule)
			if err != nil {
				return errors.Wrapf(err, "invalid implicit meta policy rule '%s'", policy.Rule)
			}
			cg.Policies[policyName] = &common.ConfigPolicy{
				ModPolicy: modPolicy,
				Policy: &common.Policy{
					Type:  int32(common.Policy_IMPLICIT_META),
					Value: utils.MarshalOrPanic(imp),
				},
			}
		case SignaturePolicyType:
			sp, err := cauthdsl.FromString(policy.Rule)
			if err != nil {
				return errors.Wrapf(err, "invalid signature policy rule '%s'", policy.Rule)
			}
			cg.Policies[policyName] = &common.ConfigPolicy{
				ModPolicy: modPolicy,
				Policy: &common.Policy{
					Type:  int32(common.Policy_SIGNATURE),
					Value: utils.MarshalOrPanic(sp),
				},
			}
		default:
			return errors.Errorf("unknown policy type: %s", policy.Type)
		}
	}
	return nil
}

// addImplicitMetaPolicyDefaults adds the Readers/Writers/Admins policies, with Any/Any/Majority rules respectively.
func addImplicitMetaPolicyDefaults(cg *common.ConfigGroup) {
	addPolicy(cg, policies.ImplicitMetaMajorityPolicy(channelconfig.AdminsPolicyKey), channelconfig.AdminsPolicyKey)
	addPolicy(cg, policies.ImplicitMetaAnyPolicy(channelconfig.ReadersPolicyKey), channelconfig.AdminsPolicyKey)
	addPolicy(cg, policies.ImplicitMetaAnyPolicy(channelconfig.WritersPolicyKey), channelconfig.AdminsPolicyKey)
}

// addSignaturePolicyDefaults adds the Readers/Writers/Admins policies as signature policies requiring one signature from the given mspID.
// If devMode is set to true, the Admins policy will accept arbitrary user certs for admin functions, otherwise it requires the cert satisfies
// the admin role principal.
func addSignaturePolicyDefaults(cg *common.ConfigGroup, mspID string, devMode bool) {
	if devMode {
		addPolicy(cg, policies.SignaturePolicy(channelconfig.AdminsPolicyKey, cauthdsl.SignedByMspMember(mspID)), channelconfig.AdminsPolicyKey)
	} else {
		addPolicy(cg, policies.SignaturePolicy(channelconfig.AdminsPolicyKey, cauthdsl.SignedByMspAdmin(mspID)), channelconfig.AdminsPolicyKey)
	}
	addPolicy(cg, policies.SignaturePolicy(channelconfig.ReadersPolicyKey, cauthdsl.SignedByMspMember(mspID)), channelconfig.AdminsPolicyKey)
	addPolicy(cg, policies.SignaturePolicy(channelconfig.WritersPolicyKey, cauthdsl.SignedByMspMember(mspID)), channelconfig.AdminsPolicyKey)
}