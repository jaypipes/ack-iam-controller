// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package role

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/iam"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/iam-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.IAM{}
	_ = &svcapitypes.Role{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.GetRoleOutput
	resp, err = rm.sdkapi.GetRoleWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "GetRole", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "NoSuchEntity" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.Role.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.Role.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.Role.AssumeRolePolicyDocument != nil {
		ko.Spec.AssumeRolePolicyDocument = resp.Role.AssumeRolePolicyDocument
	} else {
		ko.Spec.AssumeRolePolicyDocument = nil
	}
	if resp.Role.CreateDate != nil {
		ko.Status.CreateDate = &metav1.Time{*resp.Role.CreateDate}
	} else {
		ko.Status.CreateDate = nil
	}
	if resp.Role.Description != nil {
		ko.Spec.Description = resp.Role.Description
	} else {
		ko.Spec.Description = nil
	}
	if resp.Role.MaxSessionDuration != nil {
		ko.Spec.MaxSessionDuration = resp.Role.MaxSessionDuration
	} else {
		ko.Spec.MaxSessionDuration = nil
	}
	if resp.Role.Path != nil {
		ko.Spec.Path = resp.Role.Path
	} else {
		ko.Spec.Path = nil
	}
	if resp.Role.PermissionsBoundary != nil {
		ko.Spec.PermissionsBoundary = resp.Role.PermissionsBoundary.PermissionsBoundaryArn
	} else {
		ko.Spec.PermissionsBoundary = nil
	}
	if resp.Role.RoleId != nil {
		ko.Status.RoleID = resp.Role.RoleId
	} else {
		ko.Status.RoleID = nil
	}
	if resp.Role.RoleLastUsed != nil {
		f8 := &svcapitypes.RoleLastUsed{}
		if resp.Role.RoleLastUsed.LastUsedDate != nil {
			f8.LastUsedDate = &metav1.Time{*resp.Role.RoleLastUsed.LastUsedDate}
		}
		if resp.Role.RoleLastUsed.Region != nil {
			f8.Region = resp.Role.RoleLastUsed.Region
		}
		ko.Status.RoleLastUsed = f8
	} else {
		ko.Status.RoleLastUsed = nil
	}
	if resp.Role.RoleName != nil {
		ko.Spec.Name = resp.Role.RoleName
	} else {
		ko.Spec.Name = nil
	}
	if resp.Role.Tags != nil {
		f10 := []*svcapitypes.Tag{}
		for _, f10iter := range resp.Role.Tags {
			f10elem := &svcapitypes.Tag{}
			if f10iter.Key != nil {
				f10elem.Key = f10iter.Key
			}
			if f10iter.Value != nil {
				f10elem.Value = f10iter.Value
			}
			f10 = append(f10, f10elem)
		}
		ko.Spec.Tags = f10
	} else {
		ko.Spec.Tags = nil
	}

	rm.setStatusDefaults(ko)
	if ko.Spec.AssumeRolePolicyDocument != nil {
		if doc, err := decodeDocument(*ko.Spec.AssumeRolePolicyDocument); err != nil {
			return nil, err
		} else {
			ko.Spec.AssumeRolePolicyDocument = &doc
		}
	}
	ko.Spec.Policies, err = rm.getManagedPolicies(ctx, &resource{ko})
	if err != nil {
		return nil, err
	}
	ko.Spec.InlinePolicies, err = rm.getInlinePolicies(ctx, &resource{ko})
	if err != nil {
		return nil, err
	}
	ko.Spec.Tags, err = rm.getTags(ctx, &resource{ko})
	if err != nil {
		return nil, err
	}

	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Spec.Name == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.GetRoleInput, error) {
	res := &svcsdk.GetRoleInput{}

	if r.ko.Spec.Name != nil {
		res.SetRoleName(*r.ko.Spec.Name)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateRoleOutput
	_ = resp
	resp, err = rm.sdkapi.CreateRoleWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateRole", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.Role.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.Role.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.Role.AssumeRolePolicyDocument != nil {
		ko.Spec.AssumeRolePolicyDocument = resp.Role.AssumeRolePolicyDocument
	} else {
		ko.Spec.AssumeRolePolicyDocument = nil
	}
	if resp.Role.CreateDate != nil {
		ko.Status.CreateDate = &metav1.Time{*resp.Role.CreateDate}
	} else {
		ko.Status.CreateDate = nil
	}
	if resp.Role.Description != nil {
		ko.Spec.Description = resp.Role.Description
	} else {
		ko.Spec.Description = nil
	}
	if resp.Role.MaxSessionDuration != nil {
		ko.Spec.MaxSessionDuration = resp.Role.MaxSessionDuration
	} else {
		ko.Spec.MaxSessionDuration = nil
	}
	if resp.Role.Path != nil {
		ko.Spec.Path = resp.Role.Path
	} else {
		ko.Spec.Path = nil
	}
	if resp.Role.PermissionsBoundary != nil {
		ko.Spec.PermissionsBoundary = resp.Role.PermissionsBoundary.PermissionsBoundaryArn
	} else {
		ko.Spec.PermissionsBoundary = nil
	}
	if resp.Role.RoleId != nil {
		ko.Status.RoleID = resp.Role.RoleId
	} else {
		ko.Status.RoleID = nil
	}
	if resp.Role.RoleLastUsed != nil {
		f8 := &svcapitypes.RoleLastUsed{}
		if resp.Role.RoleLastUsed.LastUsedDate != nil {
			f8.LastUsedDate = &metav1.Time{*resp.Role.RoleLastUsed.LastUsedDate}
		}
		if resp.Role.RoleLastUsed.Region != nil {
			f8.Region = resp.Role.RoleLastUsed.Region
		}
		ko.Status.RoleLastUsed = f8
	} else {
		ko.Status.RoleLastUsed = nil
	}
	if resp.Role.RoleName != nil {
		ko.Spec.Name = resp.Role.RoleName
	} else {
		ko.Spec.Name = nil
	}
	if resp.Role.Tags != nil {
		f10 := []*svcapitypes.Tag{}
		for _, f10iter := range resp.Role.Tags {
			f10elem := &svcapitypes.Tag{}
			if f10iter.Key != nil {
				f10elem.Key = f10iter.Key
			}
			if f10iter.Value != nil {
				f10elem.Value = f10iter.Value
			}
			f10 = append(f10, f10elem)
		}
		ko.Spec.Tags = f10
	} else {
		ko.Spec.Tags = nil
	}

	rm.setStatusDefaults(ko)
	if ko.Spec.AssumeRolePolicyDocument != nil {
		if doc, err := decodeDocument(*ko.Spec.AssumeRolePolicyDocument); err != nil {
			return nil, err
		} else {
			ko.Spec.AssumeRolePolicyDocument = &doc
		}
	}
	ackcondition.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, nil)

	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateRoleInput, error) {
	res := &svcsdk.CreateRoleInput{}

	if r.ko.Spec.AssumeRolePolicyDocument != nil {
		res.SetAssumeRolePolicyDocument(*r.ko.Spec.AssumeRolePolicyDocument)
	}
	if r.ko.Spec.Description != nil {
		res.SetDescription(*r.ko.Spec.Description)
	}
	if r.ko.Spec.MaxSessionDuration != nil {
		res.SetMaxSessionDuration(*r.ko.Spec.MaxSessionDuration)
	}
	if r.ko.Spec.Path != nil {
		res.SetPath(*r.ko.Spec.Path)
	}
	if r.ko.Spec.PermissionsBoundary != nil {
		res.SetPermissionsBoundary(*r.ko.Spec.PermissionsBoundary)
	}
	if r.ko.Spec.Name != nil {
		res.SetRoleName(*r.ko.Spec.Name)
	}
	if r.ko.Spec.Tags != nil {
		f6 := []*svcsdk.Tag{}
		for _, f6iter := range r.ko.Spec.Tags {
			f6elem := &svcsdk.Tag{}
			if f6iter.Key != nil {
				f6elem.SetKey(*f6iter.Key)
			}
			if f6iter.Value != nil {
				f6elem.SetValue(*f6iter.Value)
			}
			f6 = append(f6, f6elem)
		}
		res.SetTags(f6)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer func() {
		exit(err)
	}()
	if delta.DifferentAt("Spec.Policies") {
		err = rm.syncManagedPolicies(ctx, desired, latest)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.InlinePolicies") {
		err = rm.syncInlinePolicies(ctx, desired, latest)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.Tags") {
		err = rm.syncTags(ctx, desired, latest)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.PermissionsBoundary") {
		err = rm.syncRolePermissionsBoundary(ctx, desired)
		if err != nil {
			return nil, err
		}
	}
	if !delta.DifferentExcept("Spec.Tags", "Spec.Policies", "Spec.InlinePolicies", "Spec.PermissionsBoundary") {
		return desired, nil
	}

	input, err := rm.newUpdateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.UpdateRoleOutput
	_ = resp
	resp, err = rm.sdkapi.UpdateRoleWithContext(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateRole", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newUpdateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
func (rm *resourceManager) newUpdateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.UpdateRoleInput, error) {
	res := &svcsdk.UpdateRoleInput{}

	if r.ko.Spec.Description != nil {
		res.SetDescription(*r.ko.Spec.Description)
	}
	if r.ko.Spec.MaxSessionDuration != nil {
		res.SetMaxSessionDuration(*r.ko.Spec.MaxSessionDuration)
	}
	if r.ko.Spec.Name != nil {
		res.SetRoleName(*r.ko.Spec.Name)
	}

	return res, nil
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	// This deletes all associated managed and inline policies from the role
	roleCpy := r.ko.DeepCopy()
	roleCpy.Spec.Policies = nil
	if err := rm.syncManagedPolicies(ctx, &resource{ko: roleCpy}, r); err != nil {
		return nil, err
	}
	roleCpy.Spec.InlinePolicies = map[string]*string{}
	if err := rm.syncInlinePolicies(ctx, &resource{ko: roleCpy}, r); err != nil {
		return nil, err
	}

	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteRoleOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteRoleWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteRole", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteRoleInput, error) {
	res := &svcsdk.DeleteRoleInput{}

	if r.ko.Spec.Name != nil {
		res.SetRoleName(*r.ko.Spec.Name)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Role,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	if err == nil {
		return false
	}
	awsErr, ok := ackerr.AWSError(err)
	if !ok {
		return false
	}
	switch awsErr.Code() {
	case "InvalidInput",
		"MalformedPolicyDocument":
		return true
	default:
		return false
	}
}
