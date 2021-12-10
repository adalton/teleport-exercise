/*
Copyright 2021 Andy Dalton
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

// Package cgroup provides an simple abstraction over the cgroup v1 interface.
// The rationale to implement was driven by the fact that on my system, both
// the v1 and v2 interfaces are mounted and the v1 interface is in use.
//
// I've versioned the package structure so that we could easily build a v2
// implementation.  We could also add some code to the parent package to
// enable us to examine the system and automatically pick the most suitable
// implementation.
package cgroup
