// Copyright 2022 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package common

import "github.com/cloudbase/garm/params"

// Each struct is a wrapper for the actual parameters struct for a specific version.
// Version 0.1.0 doesn't have any specific parameters, so there is no need for a struct for it.
type CreateInstanceParams struct {
	CreateInstanceV011 CreateInstanceV011Params
}

type DeleteInstanceParams struct {
	DeleteInstanceV011 DeleteInstanceV011Params
}

type GetInstanceParams struct {
	GetInstanceV011 GetInstanceV011Params
}

type ListInstancesParams struct {
	ListInstancesV011 ListInstancesV011Params
}

type RemoveAllInstancesParams struct {
	RemoveAllInstancesV011 RemoveAllInstancesV011Params
}

type StopParams struct {
	StopV011 StopV011Params
}

type StartParams struct {
	StartV011 StartV011Params
}

// Structs for version v0.1.1.
type CreateInstanceV011Params struct {
	PoolInfo       params.Pool
	ControllerInfo params.ControllerInfo
}

type DeleteInstanceV011Params struct {
	PoolInfo       params.Pool
	ControllerInfo params.ControllerInfo
}

type GetInstanceV011Params struct {
	PoolInfo       params.Pool
	ControllerInfo params.ControllerInfo
}

type ListInstancesV011Params struct {
	PoolInfo       params.Pool
	ControllerInfo params.ControllerInfo
}

type RemoveAllInstancesV011Params struct {
	PoolInfo       params.Pool
	ControllerInfo params.ControllerInfo
}

type StopV011Params struct {
	PoolInfo       params.Pool
	ControllerInfo params.ControllerInfo
}

type StartV011Params struct {
	PoolInfo       params.Pool
	ControllerInfo params.ControllerInfo
}
