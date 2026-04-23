package entities

import (
	"testing"
)

func TestNewRepository(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		dir     string
		wantErr bool
	}{
		{"https URL", "https://github.com/fabianoflorentino/repo.git", "repo", false},
		{"SSH URL", "git@github.com:fabianoflorentino/repo.git", "repo", false},
		{"empty URL", "", "repo", true},
		{"empty dir", "https://github.com/fabianoflorentino/repo.git", "", true},
		{"invalid URL", "invalid", "repo", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRepository(tt.url, tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Directory != tt.dir {
				t.Errorf("Directory = %v, want %v", got.Directory, tt.dir)
			}
		})
	}
}

func TestRepository_SetBranch(t *testing.T) {
	tests := []struct {
		name    string
		repo    *Repository
		branch  string
		wantErr bool
	}{
		{"valid branch", &Repository{
			URL:        "https://github.com/fabiano/repo.git",
			Directory:  "repo",
			Branch:     "main",
			AuthMethod: AuthenticationMethod{Type: AuthTypeSSH, SSHKey: SSHKeyConfig{Name: "id_rsa"}},
		}, "develop", false},
		{"empty branch", &Repository{
			URL:        "https://github.com/fabiano/repo.git",
			Directory:  "repo",
			Branch:     "main",
			AuthMethod: AuthenticationMethod{Type: AuthTypeSSH, SSHKey: SSHKeyConfig{Name: "id_rsa"}},
		}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repo.SetBranch(tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_IsSSH(t *testing.T) {
	repo := &Repository{
		AuthMethod: AuthenticationMethod{Type: AuthTypeSSH},
	}

	if !repo.IsSSH() {
		t.Error("IsSSH() should return true")
	}

	repo2 := &Repository{
		AuthMethod: AuthenticationMethod{Type: AuthTypeHTTPS},
	}

	if repo2.IsSSH() {
		t.Error("IsSSH() should return false")
	}
}

func TestRepository_IsHTTPS(t *testing.T) {
	repo := &Repository{
		AuthMethod: AuthenticationMethod{Type: AuthTypeHTTPS},
	}

	if !repo.IsHTTPS() {
		t.Error("IsHTTPS() should return true")
	}

	repo2 := &Repository{
		AuthMethod: AuthenticationMethod{Type: AuthTypeSSH},
	}

	if repo2.IsHTTPS() {
		t.Error("IsHTTPS() should return false")
	}
}

func TestRepository_Validate(t *testing.T) {
	tests := []struct {
		name    string
		repo    *Repository
		wantErr bool
	}{
		{
			"valid",
			&Repository{URL: "https://repo.git", Directory: "repo", Branch: "main", AuthMethod: AuthenticationMethod{Type: AuthTypeHTTPS, Username: "user", Token: "token"}},
			false,
		},
		{
			"empty URL",
			&Repository{URL: "", Directory: "repo", Branch: "main"},
			true,
		},
		{
			"empty directory",
			&Repository{URL: "https://repo.git", Directory: "", Branch: "main"},
			true,
		},
		{
			"empty branch",
			&Repository{URL: "https://repo.git", Directory: "repo", Branch: ""},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repo.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthenticationMethod_Validate(t *testing.T) {
	tests := []struct {
		name    string
		auth    AuthenticationMethod
		wantErr bool
	}{
		{
			"SSH valid",
			AuthenticationMethod{Type: AuthTypeSSH, SSHKey: SSHKeyConfig{Name: "id_rsa"}},
			false,
		},
		{
			"SSH invalid",
			AuthenticationMethod{Type: AuthTypeSSH},
			true,
		},
		{
			"HTTPS valid",
			AuthenticationMethod{Type: AuthTypeHTTPS, Username: "user", Token: "token"},
			false,
		},
		{
			"HTTPS invalid",
			AuthenticationMethod{Type: AuthTypeHTTPS},
			true,
		},
		{
			"invalid type",
			AuthenticationMethod{Type: "invalid"},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.auth.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtractRepositoryName(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"https", "https://github.com/fabianoflorentino/repo.git", "repo"},
		{"ssh", "git@github.com:fabianoflorentino/repo.git", "repo"},
		{"no git suffix", "https://github.com/org/repo", "repo"},
		{"invalid", "invalid", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractRepositoryName(tt.url)
			if got != tt.want {
				t.Errorf("extractRepositoryName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetermineAuthMethod(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want AuthType
	}{
		{"ssh", "git@github.com:org/repo.git", AuthTypeSSH},
		{"https", "https://github.com/org/repo.git", AuthTypeHTTPS},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineAuthMethod(tt.url)
			if got.Type != tt.want {
				t.Errorf("determineAuthMethod() = %v, want %v", got.Type, tt.want)
			}
		})
	}
}
