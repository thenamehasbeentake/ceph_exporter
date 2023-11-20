//   Copyright 2022 DigitalOcean
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package ceph

import (
	"reflect"
	"testing"
)

func TestParseCephVersion(t *testing.T) {
	type args struct {
		cephVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    *Version
		wantErr bool
	}{
		{
			name:    "invalid version 1",
			args:    args{cephVersion: "totally real version"},
			want:    nil,
			wantErr: true,
		},
		// {
		// 	name:    "invalid version 2",
		// 	args:    args{cephVersion: "ceph version 14.2.18-97"},
		// 	want:    nil,
		// 	wantErr: true,
		// },
		{
			name:    "nautilus",
			args:    args{cephVersion: "ceph version 14.2.18-97-gcc1e126 (cc1e1267bc7afc8288c718fc3e59c5a6735f6f4a) nautilus (stable)"},
			want:    &Version{Major: 14, Minor: 2, Patch: 18, Revision: 97, Commit: "gcc1e126"},
			wantErr: false,
		},
		{
			name:    "nautilus-ceph-ansible",
			args:    args{cephVersion: "ceph version 14.2.11-184.el8cp (44441323476fee97be0ff7a92c6065958c77f1b9) nautilus (stable)"},
			want:    &Version{Major: 14, Minor: 2, Patch: 11, Revision: 184, Commit: "el8cp"},
			wantErr: false,
		},
		{
			name:    "octopus",
			args:    args{cephVersion: "ceph version 15.2.0-1-gcc1e126"},
			want:    &Version{Major: 15, Minor: 2, Patch: 0, Revision: 1, Commit: "gcc1e126"},
			wantErr: false,
		},
		{
			name:    "pacific",
			args:    args{cephVersion: "ceph version 16.2.3-33-gcc1e126"},
			want:    &Version{Major: 16, Minor: 2, Patch: 3, Revision: 33, Commit: "gcc1e126"},
			wantErr: false,
		},
		{
			name:    "real pacific",
			args:    args{cephVersion: "ceph version 16.2.7 (dd0603118f56ab514f133c8d2e3adfc983942503) pacific (stable)"},
			want:    &Version{Major: 16, Minor: 2, Patch: 7, Revision: 0, Commit: ""},
			wantErr: false,
		},
		{
			name:    "uos",
			args:    args{cephVersion: "ceph version 14.2.11-2023071915 (e4f959eb759b9bd0153ef6bb6be2291dcd8dbc8c) nautilus (stable)"},
			want:    &Version{Major: 14, Minor: 2, Patch: 11, Revision: 2023071915, Commit: ""},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCephVersion(tt.args.cephVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCephVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCephVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_IsAtLeast(t *testing.T) {
	type fields struct {
		Major    int
		Minor    int
		Patch    int
		Revision int
		Commit   string
	}
	type args struct {
		constraint *Version
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "equal versions",
			fields: fields{Major: Nautilus.Major, Minor: Nautilus.Minor, Patch: Nautilus.Patch, Revision: Nautilus.Revision, Commit: Nautilus.Commit},
			args:   args{constraint: Nautilus},
			want:   true,
		},
		{
			name:   "slightly older",
			fields: fields{Major: Pacific.Major, Minor: Pacific.Minor - 1, Patch: Pacific.Patch, Revision: Pacific.Revision, Commit: Pacific.Commit},
			args:   args{constraint: Pacific},
			want:   false,
		},
		{
			name:   "significantly newer",
			fields: fields{Major: Pacific.Major, Minor: Pacific.Minor - 1, Patch: Pacific.Patch, Revision: Pacific.Revision, Commit: Pacific.Commit},
			args:   args{constraint: Octopus},
			want:   true,
		},
		{
			name:   "older revision",
			fields: fields{Major: 16, Minor: 2, Patch: 0, Revision: 1},
			args:   args{constraint: &Version{Major: 16, Minor: 2, Patch: 0, Revision: 2}},
			want:   false,
		},
		{
			name:   "newer revision",
			fields: fields{Major: 16, Minor: 2, Patch: 0, Revision: 1},
			args:   args{constraint: Pacific},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version := &Version{
				Major:    tt.fields.Major,
				Minor:    tt.fields.Minor,
				Patch:    tt.fields.Patch,
				Revision: tt.fields.Revision,
				Commit:   tt.fields.Commit,
			}
			if got := version.IsAtLeast(tt.args.constraint); got != tt.want {
				t.Errorf("IsAtLeast() = %v, want %v", got, tt.want)
			}
		})
	}
}
