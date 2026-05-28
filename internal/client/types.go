package client

import "time"

// APIError represents an error response from the Simfra Admin API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}

func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 404
	}
	return false
}

// --- Account types ---

type CreateAccountRequest struct {
	AccountID         string   `json:"accountId,omitempty"`
	Bootstrap         string   `json:"bootstrap,omitempty"`
	Region            string   `json:"region,omitempty"`
	AvailabilityZones []string `json:"availability_zones,omitempty"`
	VpcCIDR           string   `json:"vpc_cidr,omitempty"`
}

type AccountDetail struct {
	AccountID          string `json:"accountId"`
	RootAccessKeyID    string `json:"rootAccessKeyId"`
	RootSecretAccessKey string `json:"rootSecretAccessKey"`
	CreatedAt          string `json:"createdAt"`
}

type AccountSummary struct {
	AccountID string `json:"accountId"`
	Alias     string `json:"alias"`
	CreatedAt string `json:"createdAt"`
}

// --- Port Forward types ---

type CreatePortForwardRequest struct {
	TargetARN  string `json:"targetArn"`
	TargetPort int    `json:"targetPort,omitempty"`
	Service    string `json:"service,omitempty"`
	ResourceID string `json:"resourceId,omitempty"`
	AccountID  string `json:"accountId,omitempty"`
	Region     string `json:"region,omitempty"`
}

type PortForwardSession struct {
	ID          string    `json:"id"`
	TargetARN   string    `json:"targetArn"`
	TargetIP    string    `json:"targetIp"`
	TargetPort  int       `json:"targetPort"`
	LocalPort   int       `json:"localPort"`
	LocalAddress string   `json:"localAddress"`
	VPCNetwork  string    `json:"vpcNetwork"`
	ContainerID string    `json:"containerId"`
	Service     string    `json:"service"`
	ResourceID  string    `json:"resourceId"`
	AccountID   string    `json:"accountId"`
	Region      string    `json:"region"`
	CreatedAt   time.Time `json:"createdAt"`
	Status      string    `json:"status"`
}

type PortForwardTarget struct {
	ARN         string `json:"arn"`
	Service     string `json:"service"`
	ResourceID  string `json:"resourceId"`
	AccountID   string `json:"accountId"`
	Region      string `json:"region"`
	DefaultPort int    `json:"defaultPort"`
	VPCNetwork  string `json:"vpcNetwork"`
}

// --- SSO Session types ---

type CreateSSOSessionRequest struct {
	UserID       string              `json:"user_id"`
	UserName     string              `json:"user_name"`
	Assignments  []SSOAssignment     `json:"assignments"`
	ExpiresInSec int                 `json:"expires_in_sec,omitempty"`
}

type SSOAssignment struct {
	AccountID    string `json:"account_id"`
	AccountName  string `json:"account_name,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	RoleName     string `json:"role_name,omitempty"`
	RoleARN      string `json:"role_arn,omitempty"`
}

type SSOSessionResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SSOSessionInfo struct {
	Token       string          `json:"token"`
	UserID      string          `json:"user_id"`
	UserName    string          `json:"user_name"`
	Assignments []SSOAssignment `json:"assignments"`
	ExpiresAt   time.Time       `json:"expires_at"`
	CreatedAt   time.Time       `json:"created_at"`
	Expired     bool            `json:"expired"`
}

// --- ACM types ---

type ACMPendingValidation struct {
	CertificateARN string                   `json:"certificateArn"`
	DomainName     string                   `json:"domainName"`
	Status         string                   `json:"status"`
	Domains        []ACMDomainValidation    `json:"domains"`
}

type ACMDomainValidation struct {
	DomainName       string                 `json:"domainName"`
	ValidationMethod string                 `json:"validationMethod"`
	ValidationStatus string                 `json:"validationStatus"`
	ResourceRecord   *ACMResourceRecord     `json:"resourceRecord"`
}

type ACMResourceRecord struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// --- Health ---

type HealthResponse struct {
	Status string `json:"status"`
}

// --- Services ---

type ServiceInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Protocols   []string `json:"protocols"`
	Operations  []string `json:"operations"`
}

// --- Storage ---

type StorageSummary struct {
	Enabled  bool                `json:"enabled"`
	DataDir  string              `json:"dataDir"`
	DBPath   string              `json:"dbPath"`
	DBSize   int64               `json:"dbSize"`
	Total    int64               `json:"total"`
	Services []StorageServiceInfo `json:"services"`
}

type StorageServiceInfo struct {
	Service string `json:"service"`
	Size    int64  `json:"size"`
}

// --- Docker ---

type DockerSummary struct {
	Containers DockerContainerCounts `json:"containers"`
	Images     int                   `json:"images"`
	Networks   int                   `json:"networks"`
	Volumes    int                   `json:"volumes"`
}

type DockerContainerCounts struct {
	Total   int `json:"total"`
	Running int `json:"running"`
	Stopped int `json:"stopped"`
}

type DockerContainer struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	State   string            `json:"state"`
	Status  string            `json:"status"`
	Created string            `json:"created"`
	Labels  map[string]string `json:"labels"`
	Ports   []DockerPort      `json:"ports"`
}

type DockerPort struct {
	IP          string `json:"ip"`
	PrivatePort int    `json:"privatePort"`
	PublicPort  int    `json:"publicPort"`
	Type        string `json:"type"`
}

type DockerNetwork struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Driver  string            `json:"driver"`
	Scope   string            `json:"scope"`
	Labels  map[string]string `json:"labels"`
}

// --- CA ---

type CAInfo struct {
	Persistent   bool         `json:"persistent"`
	Directory    string       `json:"directory"`
	Root         CACertInfo   `json:"root"`
	Intermediate CACertInfo   `json:"intermediate"`
}

type CACertInfo struct {
	Subject      string `json:"subject"`
	SerialNumber string `json:"serialNumber"`
	NotAfter     string `json:"notAfter"`
}

// --- DNS / SMTP ---

type PortResponse struct {
	Port int `json:"port"`
}
