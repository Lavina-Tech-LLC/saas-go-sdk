package auth

import "testing"

func TestMemberHasAnyRole(t *testing.T) {
	cases := []struct {
		name    string
		member  Member
		allowed []string
		want    bool
	}{
		{
			name: "custom role via Roles array when primary is system role",
			member: Member{
				Role:  "member",
				Roles: []MemberRole{{Key: "editor", ID: "r-1"}},
			},
			allowed: []string{"editor"},
			want:    true,
		},
		{
			name: "match by role UUID from Roles array",
			member: Member{
				Role:  "member",
				Roles: []MemberRole{{Key: "editor", ID: "r-1"}},
			},
			allowed: []string{"r-1"},
			want:    true,
		},
		{
			name: "owner plus custom — custom matches",
			member: Member{
				Role: "owner",
				Roles: []MemberRole{
					{Key: "owner", ID: "r-0"},
					{Key: "editor", ID: "r-1"},
				},
			},
			allowed: []string{"editor"},
			want:    true,
		},
		{
			name:    "older backend without Roles array falls back to primary",
			member:  Member{Role: "member"},
			allowed: []string{"member"},
			want:    true,
		},
		{
			name:    "negative control — primary and Roles both miss",
			member:  Member{Role: "member"},
			allowed: []string{"editor"},
			want:    false,
		},
		{
			name:    "empty primary role never matches",
			member:  Member{Role: ""},
			allowed: []string{"member"},
			want:    false,
		},
		{
			name:    "empty allowed identifier never matches",
			member:  Member{Role: "admin"},
			allowed: []string{""},
			want:    false,
		},
		{
			name:    "primary RoleID matches",
			member:  Member{Role: "member", RoleID: "r-m"},
			allowed: []string{"r-m"},
			want:    true,
		},
		{
			name: "multiple allowed, second matches via Roles",
			member: Member{
				Role:  "member",
				Roles: []MemberRole{{Key: "editor", ID: "r-1"}},
			},
			allowed: []string{"owner", "admin", "editor"},
			want:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.member.HasAnyRole(tc.allowed...)
			if got != tc.want {
				t.Errorf("HasAnyRole(%v) = %v, want %v", tc.allowed, got, tc.want)
			}
		})
	}
}
