/*
Copyright 2023 Pepov.

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
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1beta1

// ConfigSpecApplyConfiguration represents an declarative configuration of the ConfigSpec type for use
// with apply.
type ConfigSpecApplyConfiguration struct {
	Foo   *string `json:"foo,omitempty"`
	Other *string `json:"other,omitempty"`
}

// ConfigSpecApplyConfiguration constructs an declarative configuration of the ConfigSpec type for use with
// apply.
func ConfigSpec() *ConfigSpecApplyConfiguration {
	return &ConfigSpecApplyConfiguration{}
}

// WithFoo sets the Foo field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Foo field is set to the value of the last call.
func (b *ConfigSpecApplyConfiguration) WithFoo(value string) *ConfigSpecApplyConfiguration {
	b.Foo = &value
	return b
}

// WithOther sets the Other field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Other field is set to the value of the last call.
func (b *ConfigSpecApplyConfiguration) WithOther(value string) *ConfigSpecApplyConfiguration {
	b.Other = &value
	return b
}
