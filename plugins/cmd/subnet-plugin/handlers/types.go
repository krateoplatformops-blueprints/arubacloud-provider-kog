package subnet

type SubnetType string

const (
	SubnetTypeBasic    SubnetType = "Basic"
	SubnetTypeAdvanced SubnetType = "Advanced"
)

type SubnetDto struct {
	Metadata   *MetadataDto         `json:"metadata,omitempty"`
	Properties *SubnetPropertiesDto `json:"properties,omitempty"`
}

type MetadataDto struct {
	Name     string       `json:"name,omitempty"`
	Location *LocationDto `json:"location,omitempty"`
	Tags     []string     `json:"tags,omitempty"`
}

type LocationDto struct {
	Value string `json:"value,omitempty"`
}

type SubnetPropertiesDto struct {
	// "Type of the subnet.\r\nAvailable values:\r\n- Basic\r\n- Advanced\r\n\r\nWith Basic type, every configuration settings of the subnet will be automatically handled by the CMP.\r\nWith Advanced type, configuration settings must be evaluated by the user."
	Type SubnetType `json:"type,omitempty"`
	// "Indicates if the subnet must be a default subnet.\r\nOnly one default subnet for vpc is admissible."
	Default bool        `json:"default,omitempty"`
	Network *NetworkDto `json:"network,omitempty"`
	Dhcp    *DhcpDto    `json:"dhcp,omitempty"`
}

type NetworkDto struct {
	Address string `json:"address,omitempty"`
}

type DhcpDto struct {
	Enabled bool       `json:"enabled,omitempty"`
	Range   *RangeDto  `json:"range,omitempty"`
	Routes  []RouteDto `json:"routes,omitempty"`
	Dns     []string   `json:"dns,omitempty"`
}

type RangeDto struct {
	Start string `json:"start,omitempty"`
	Count int32  `json:"count,omitempty"`
}

type RouteDto struct {
	Address string `json:"address,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

type SubnetUpdateDto struct {
	Metadata   *MetadataDto               `json:"metadata,omitempty"`
	Properties *SubnetUpdatePropertiesDto `json:"properties,omitempty"`
}

type SubnetUpdatePropertiesDto struct {
	Default bool `json:"default,omitempty"`
}

type SubnetResponseDto struct {
	Metadata   *MetadataResponseDto         `json:"metadata,omitempty"`
	Status     *StatusResponseDto           `json:"status,omitempty"`
	Properties *SubnetPropertiesResponseDto `json:"properties,omitempty"`
}

type MetadataResponseDto struct {
	ID           string               `json:"id,omitempty"`
	URI          string               `json:"uri,omitempty"`
	Name         string               `json:"name,omitempty"`
	Location     *LocationResponseDto `json:"location,omitempty"`
	Project      *ProjectResponseDto  `json:"project,omitempty"`
	Tags         []string             `json:"tags,omitempty"`
	Category     *CategoryResponseDto `json:"category,omitempty"`
	CreationDate string               `json:"creationDate,omitempty"`
	CreatedBy    string               `json:"createdBy,omitempty"`
	UpdateDate   string               `json:"updateDate,omitempty"`
	UpdatedBy    string               `json:"updatedBy,omitempty"`
	Version      string               `json:"version,omitempty"`
	CreatedUser  string               `json:"createdUser,omitempty"`
	UpdatedUser  string               `json:"updatedUser,omitempty"`
}

type LocationResponseDto struct {
	Code    string `json:"code,omitempty"`
	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`
	Name    string `json:"name,omitempty"`
	Value   string `json:"value,omitempty"`
}

type ProjectResponseDto struct {
	ID string `json:"id,omitempty"`
}

type CategoryResponseDto struct {
	Name     string               `json:"name,omitempty"`
	Provider string               `json:"provider,omitempty"`
	Typology *TypologyResponseDto `json:"typology,omitempty"`
}

type TypologyResponseDto struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type StatusResponseDto struct {
	State             string                        `json:"state,omitempty"`
	CreationDate      string                        `json:"creationDate,omitempty"`
	DisableStatusInfo *DisableStatusInfoResponseDto `json:"disableStatusInfo,omitempty"`
	FailureReason     string                        `json:"failureReason,omitempty"`
}

type DisableStatusInfoResponseDto struct {
	IsDisabled     bool                       `json:"isDisabled,omitempty"`
	Reasons        []string                   `json:"reasons,omitempty"`
	PreviousStatus *PreviousStatusResponseDto `json:"previousStatus,omitempty"`
}

type PreviousStatusResponseDto struct {
	State        string `json:"state,omitempty"`
	CreationDate string `json:"creationDate,omitempty"`
}

type SubnetPropertiesResponseDto struct {
	LinkedResources []LinkedResourceResponseDto `json:"linkedResources,omitempty"`
	Vpc             *GenericResourceResponseDto `json:"vpc,omitempty"`
	Type            SubnetType                  `json:"type,omitempty"`
	Default         bool                        `json:"default,omitempty"`
	Network         *NetworkResponseDto         `json:"network,omitempty"`
	Dhcp            *DhcpResponseDto            `json:"dhcp,omitempty"`
}

type LinkedResourceResponseDto struct {
	URI               string `json:"uri,omitempty"`
	StrictCorrelation bool   `json:"strictCorrelation,omitempty"`
}

type GenericResourceResponseDto struct {
	URI string `json:"uri,omitempty"`
}

type NetworkResponseDto struct {
	Address string `json:"address,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

type DhcpResponseDto struct {
	Enabled bool               `json:"enabled,omitempty"`
	Range   *RangeResponseDto  `json:"range,omitempty"`
	Routes  []RouteResponseDto `json:"routes,omitempty"`
	Dns     []string           `json:"dns,omitempty"`
}

type RangeResponseDto struct {
	Start string `json:"start,omitempty"`
	Count int32  `json:"count,omitempty"`
	Last  string `json:"last,omitempty"`
}

type RouteResponseDto struct {
	Address string `json:"address,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

type SubnetListResponseDto struct {
	Total  int64               `json:"total,omitempty"`
	Self   string              `json:"self,omitempty"`
	Prev   string              `json:"prev,omitempty"`
	Next   string              `json:"next,omitempty"`
	First  string              `json:"first,omitempty"`
	Last   string              `json:"last,omitempty"`
	Values []SubnetResponseDto `json:"values,omitempty"`
}

type ProblemDetails struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int32  `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// --------------------------------------------------------------------------
// Flattened types
// --------------------------------------------------------------------------

// FlattenedCreateSubnetRequestDto is the flattened request body for creating a subnet.
// The fields from MetadataDto are at the root level, alongside the Properties object.
type FlattenedCreateSubnetRequestDto struct {
	Name       string               `json:"name,omitempty"`
	Location   *LocationDto         `json:"location,omitempty"`
	Tags       []string             `json:"tags,omitempty"`
	Properties *SubnetPropertiesDto `json:"properties,omitempty"`
}

// FlattenedUpdateSubnetRequestDto is the flattened request body for updating a subnet.
// The fields from MetadataDto are at the root level, alongside the Properties object.
type FlattenedUpdateSubnetRequestDto struct {
	Name       string                     `json:"name,omitempty"`
	Location   *LocationDto               `json:"location,omitempty"`
	Tags       []string                   `json:"tags,omitempty"`
	Properties *SubnetUpdatePropertiesDto `json:"properties,omitempty"`
}

// FlattenedSubnetResponseDto is the flattened response body for a single subnet.
// The fields from MetadataResponseDto are at the root level, alongside Status and Properties objects.
type FlattenedSubnetResponseDto struct {
	ID           string                       `json:"id,omitempty"`
	URI          string                       `json:"uri,omitempty"`
	Name         string                       `json:"name,omitempty"`
	Location     *LocationResponseDto         `json:"location,omitempty"`
	Project      *ProjectResponseDto          `json:"project,omitempty"`
	Tags         []string                     `json:"tags,omitempty"`
	Category     *CategoryResponseDto         `json:"category,omitempty"`
	CreationDate string                       `json:"creationDate,omitempty"`
	CreatedBy    string                       `json:"createdBy,omitempty"`
	UpdateDate   string                       `json:"updateDate,omitempty"`
	UpdatedBy    string                       `json:"updatedBy,omitempty"`
	Version      string                       `json:"version,omitempty"`
	CreatedUser  string                       `json:"createdUser,omitempty"`
	UpdatedUser  string                       `json:"updatedUser,omitempty"`
	Status       *StatusResponseDto           `json:"status,omitempty"`
	Properties   *SubnetPropertiesResponseDto `json:"properties,omitempty"`
}

type FlattenedSubnetListResponseDto struct {
	Total  int64                        `json:"total,omitempty"`
	Self   string                       `json:"self,omitempty"`
	Prev   string                       `json:"prev,omitempty"`
	Next   string                       `json:"next,omitempty"`
	First  string                       `json:"first,omitempty"`
	Last   string                       `json:"last,omitempty"`
	Values []FlattenedSubnetResponseDto `json:"values,omitempty"`
}
