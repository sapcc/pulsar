/*******************************************************************************
*
* Copyright 2019 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package auth

// UserRole ...
type UserRole string

// UserRoles enumerates available UserRole.
var UserRoles = struct {

	// Base role required for any interaction with the bot.
	Base,

	// KubernetesAdmin is required for admin operations in Kubernetes clusters via the bot.
	KubernetesAdmin,

	// KubernetesUser is required for reading operations in Kubernetes clusters via the bot.
	KubernetesUser UserRole
}{
	"Base",
	"KubernetesAdmin",
	"KubernetesUser",
}
