/*
* Copyright 2018 - Present Okta, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

// AUTO-GENERATED!  DO NOT EDIT FILE DIRECTLY

package okta

import (
	"time"
)

type BookmarkApplication struct {
	Embedded      interface{}                  `json:"_embedded,omitempty"`
	Links         interface{}                  `json:"_links,omitempty"`
	Accessibility *ApplicationAccessibility    `json:"accessibility,omitempty"`
	Created       *time.Time                   `json:"created,omitempty"`
	Credentials   *ApplicationCredentials      `json:"credentials,omitempty"`
	Features      []string                     `json:"features,omitempty"`
	Id            string                       `json:"id,omitempty"`
	Label         string                       `json:"label,omitempty"`
	LastUpdated   *time.Time                   `json:"lastUpdated,omitempty"`
	Licensing     *ApplicationLicensing        `json:"licensing,omitempty"`
	Name          string                       `json:"name,omitempty"`
	Profile       interface{}                  `json:"profile,omitempty"`
	Settings      *BookmarkApplicationSettings `json:"settings,omitempty"`
	SignOnMode    string                       `json:"signOnMode,omitempty"`
	Status        string                       `json:"status,omitempty"`
	Visibility    *ApplicationVisibility       `json:"visibility,omitempty"`
}

func NewBookmarkApplication() *BookmarkApplication {
	return &BookmarkApplication{
		Name:       "bookmark",
		SignOnMode: "BOOKMARK",
	}
}

func (a *BookmarkApplication) IsApplicationInstance() bool {
	return true
}
