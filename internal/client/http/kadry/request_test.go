package kadry

import (
	"reflect"
	"testing"
)

func Test_newRequest(t *testing.T) {
	type args struct {
		persons []string
		attrs   []AttributeName
	}
	tests := []struct {
		name string
		args args
		want *request
	}{
		{
			name: "empty attributes",
			args: args{
				persons: []string{"person_id_1", "person_id_2", "person_id_3"},
			},
			want: &request{
				PersonIDArray: []string{"person_id_1", "person_id_2", "person_id_3"},
			},
		},
		{
			name: "empty attributes",
			args: args{
				persons: []string{"person_id_1", "person_id_2"},
				attrs: []AttributeName{
					PersonID,
					FIOPerson,
					SNILS,
					OrgID,
					InnOrg,
					NameOrg,
					SubdivID,
					NameSubdiv,
					PositionID,
					NamePosition,
					EmploymentType,
					DateRecept,
				},
			},
			want: &request{
				PersonIDArray: []string{"person_id_1", "person_id_2"},
				AttributeList: &attributeList{
					MobileApp: attributes{
						Include: []AttributeName{
							PersonID,
							FIOPerson,
							SNILS,
							OrgID,
							InnOrg,
							NameOrg,
							SubdivID,
							NameSubdiv,
							PositionID,
							NamePosition,
							EmploymentType,
							DateRecept,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRequest(tt.args.persons, tt.args.attrs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
