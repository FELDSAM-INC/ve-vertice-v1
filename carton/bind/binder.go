/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package bind

import (
	"fmt"
)

// EnvVar represents a environment variable for a carton.
type EnvVar struct {
	Name  string
	Value string
}

func (e *EnvVar) String() string {
	return fmt.Sprintf("%s=%s", e.Name, e.Value)
}

type Binder interface {

	// Bind makes the bind between two boxes.
	//Bind(b *provision.Box) error

	// Unbind makes the unbind between two boxes
	//Unbind(b *provision.Box) error

	// Provides the YetToBeBoud instances for a box.
	//Group() (*[]YBoundBox, error)
}

/*Yet to be bound instance for a box.
Details needed like the Envs (in a map), name.domain
*/
type YBoundBox struct {
	Name string
	Envs map[string]string
}
