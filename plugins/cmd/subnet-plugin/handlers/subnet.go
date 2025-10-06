package subnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/krateoplatformops/arubacloud-provider-kog/pkg/handlers"
	"github.com/krateoplatformops/arubacloud-provider-kog/pkg/utils"
)

func GetSubnet(opts handlers.HandlerOptions) handlers.Handler {
	return &getHandler{baseHandler: newBaseHandler(opts)}
}

func PostSubnet(opts handlers.HandlerOptions) handlers.Handler {
	return &postHandler{baseHandler: newBaseHandler(opts)}
}

func PutSubnet(opts handlers.HandlerOptions) handlers.Handler {
	return &putHandler{baseHandler: newBaseHandler(opts)}
}

func ListSubnets(opts handlers.HandlerOptions) handlers.Handler {
	return &listHandler{baseHandler: newBaseHandler(opts)}
}

// Interface compliance verification
var _ handlers.Handler = &getHandler{}
var _ handlers.Handler = &postHandler{}
var _ handlers.Handler = &putHandler{}
var _ handlers.Handler = &listHandler{}

// Base handler with common functionality
type baseHandler struct {
	handlers.HandlerOptions
}

// Constructor for the base handler
func newBaseHandler(opts handlers.HandlerOptions) *baseHandler {
	return &baseHandler{HandlerOptions: opts}
}

// Handler types embedding the base handler
type getHandler struct {
	*baseHandler
}

type postHandler struct {
	*baseHandler
}

type putHandler struct {
	*baseHandler
}

type listHandler struct {
	*baseHandler
}

// Common methods, defined once on baseHandler
func (h *baseHandler) makeArubaCloudRequest(method, url string, authHeader string, body []byte) (*http.Response, error) {
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if authHeader != "" {
		h.Log.Print("Using provided Authorization header for Bearer authentication")
		req.Header.Set("Authorization", authHeader)
	} else {
		h.Log.Print("No Authorization header provided, Bearer authentication required")
		return nil, fmt.Errorf("no Authorization header provided, Bearer authentication required")
	}

	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

func (h *baseHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	h.Log.Print(message)
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func (h *baseHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}

// GET handler implementation
// @Summary Get a Subnet from Aruba Cloud
// @Description Get a Subnet from Aruba Cloud using the provided project, vpc, and subnet details.
// @ID get-subnet
// @Param projectId path string true "Project ID"
// @Param vpcId path string true "VPC ID"
// @Param id path string true "Subnet ID"
// @Param api-version query string true "API version (e.g., 1.0)"
// @Param ignoreDeletedStatus query boolean false "if the resource exists in status 'Deleted', returns NotFound according to the value of this flag"
// @Param Authorization header string true "Bearer Token (Bearer <token>)"
// @Accept json
// @Produce json
// @Success 200 {object} FlattenedSubnetResponseDto "Subnet details"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /projects/{projectId}/providers/Aruba.Network/vpcs/{vpcId}/subnets/{id} [get]
func (h *getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("projectId")
	vpcId := r.PathValue("vpcId")
	id := r.PathValue("id")
	authHeader := r.Header.Get("Authorization")

	// Validate required parameters
	if projectId == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Project ID parameter is required")
		return
	}
	if vpcId == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VPC ID parameter is required")
		return
	}
	if id == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Subnet ID parameter is required")
		return
	}
	queryParams := r.URL.Query()
	if apiVersion := queryParams.Get("api-version"); apiVersion == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "API version parameter is required")
		return
	}
	if authHeader == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	// Construct the URL for the Aruba Cloud API
	baseURL := fmt.Sprintf("https://api.arubacloud.com/projects/%s/providers/Aruba.Network/vpcs/%s/subnets/%s", projectId, vpcId, id)
	url := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	// Make the GET request to Aruba Cloud API
	resp, err := h.makeArubaCloudRequest("GET", url, authHeader, nil)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to make get subnet request: %v", err))
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to read get subnet response")
		return
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		h.Log.Printf("Aruba Cloud API returned non-200 status for get subnet: %d. Body: %s", resp.StatusCode, string(body))
		// Proxy the original error response
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return
	}

	// Unmarshal the response into the Go struct to validate it
	var arubaResponse SubnetResponseDto
	if err := json.Unmarshal(body, &arubaResponse); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal Aruba Cloud response: %v", err))
		return
	}

	// Marshal the validated struct back to JSON to prepare for flattening
	validatedBody, err := json.Marshal(arubaResponse)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to marshal validated response: %v", err))
		return
	}

	// Flatten the validated response
	flattenedBody, err := utils.FlattenObject(validatedBody, "metadata") // Move contents of "metadata" to top level
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to flatten response: %v", err))
		return
	}

	h.writeJSONResponse(w, http.StatusOK, flattenedBody)
	h.Log.Printf("Successfully retrieved and flattened subnet '%s'", id)
}

// POST handler implementation
// @Summary Create a new Subnet on Aruba Cloud
// @Description Create a new Subnet on Aruba Cloud using the provided project and vpc details.
// @ID post-subnet
// @Param projectId path string true "Project ID"
// @Param vpcId path string true "VPC ID"
// @Param api-version query string true "API version (e.g., 1.0)"
// @Param Authorization header string true "Bearer Token (Bearer <token>)"
// @Param subnetCreate body FlattenedCreateSubnetRequestDto true "Subnet creation request body"
// @Accept json
// @Produce json
// @Success 201 {object} FlattenedSubnetResponseDto "Subnet details"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /projects/{projectId}/providers/Aruba.Network/vpcs/{vpcId}/subnets [post]
func (h *postHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("projectId")
	vpcId := r.PathValue("vpcId")
	apiVersion := r.URL.Query().Get("api-version")
	authHeader := r.Header.Get("Authorization")

	// Validate required parameters
	if projectId == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Project ID parameter is required")
		return
	}
	if vpcId == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VPC ID parameter is required")
		return
	}
	if apiVersion == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "API version parameter is required")
		return
	}
	if authHeader == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	// Read and parse the flattened request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var flattenedRequest FlattenedCreateSubnetRequestDto
	if err := json.Unmarshal(body, &flattenedRequest); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON in request body")
		return
	}

	// "Unflatten" the request body: build the nested structure that Aruba Cloud expects
	arubaRequest := SubnetDto{
		Metadata: &MetadataDto{
			Name:     flattenedRequest.Name,
			Location: flattenedRequest.Location,
			Tags:     flattenedRequest.Tags,
		},
		Properties: flattenedRequest.Properties, // Directly use the nested Properties struct
	}

	arubaRequestBody, err := json.Marshal(arubaRequest)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to marshal Aruba Cloud request body")
		return
	}

	// Construct the URL for the Aruba Cloud API
	url := fmt.Sprintf("https://api.arubacloud.com/projects/%s/providers/Aruba.Network/vpcs/%s/subnets?api-version=%s", projectId, vpcId, apiVersion)

	// Make the POST request to Aruba Cloud API
	resp, err := h.makeArubaCloudRequest("POST", url, authHeader, arubaRequestBody)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to make create subnet request: %v", err))
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to read create subnet response")
		return
	}

	// Check for non-201 status codes
	if resp.StatusCode != http.StatusCreated {
		h.Log.Printf("Aruba Cloud API returned non-201 status for create subnet: %d. Body: %s", resp.StatusCode, string(respBody))
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
		return
	}

	// Unmarshal the response into the Go struct to validate it
	var arubaResponse SubnetResponseDto
	if err := json.Unmarshal(respBody, &arubaResponse); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal Aruba Cloud response: %v", err))
		return
	}

	// Marshal the validated struct back to JSON to prepare for flattening
	validatedBody, err := json.Marshal(arubaResponse)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to marshal validated response: %v", err))
		return
	}

	// Flatten the response
	flattenedBody, err := utils.FlattenObject(validatedBody, "metadata")
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to flatten response: %v", err))
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, flattenedBody)
	h.Log.Printf("Successfully created subnet in project '%s', vpc '%s'", projectId, vpcId)
}

// PUT handler implementation
// @Summary Update a Subnet on Aruba Cloud
// @Description Update a Subnet on Aruba Cloud using the provided project, vpc, and subnet details.
// @ID put-subnet
// @Param projectId path string true "Project ID"
// @Param vpcId path string true "VPC ID"
// @Param id path string true "Subnet ID"
// @Param api-version query string true "API version (e.g., 1.0)"
// @Param Authorization header string true "Bearer Token (Bearer <token>)"
// @Param subnetUpdate body FlattenedUpdateSubnetRequestDto true "Subnet update request body"
// @Accept json
// @Produce json
// @Success 200 {object} FlattenedSubnetResponseDto "Subnet details"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Router /projects/{projectId}/providers/Aruba.Network/vpcs/{vpcId}/subnets/{id} [put]
func (h *putHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("projectId")
	vpcId := r.PathValue("vpcId")
	id := r.PathValue("id")
	apiVersion := r.URL.Query().Get("api-version")
	authHeader := r.Header.Get("Authorization")

	// Validate required parameters
	if projectId == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Project ID parameter is required")
		return
	}
	if vpcId == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VPC ID parameter is required")
		return
	}
	if id == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Subnet ID parameter is required")
		return
	}
	if apiVersion == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "API version parameter is required")
		return
	}
	if authHeader == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	// Read and parse the flattened request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var flattenedRequest FlattenedUpdateSubnetRequestDto
	if err := json.Unmarshal(body, &flattenedRequest); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON in request body")
		return
	}

	// "Unflatten" the request
	arubaRequest := SubnetUpdateDto{
		Metadata: &MetadataDto{
			Name:     flattenedRequest.Name,
			Location: flattenedRequest.Location,
			Tags:     flattenedRequest.Tags,
		},
		Properties: flattenedRequest.Properties, // Directly use the nested Properties struct
	}

	arubaRequestBody, err := json.Marshal(arubaRequest)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to marshal ArubaCloud request body")
		return
	}

	h.Log.Printf("Request body to send to Aruba Cloud: %s", string(arubaRequestBody))

	// Construct the URL for the Aruba Cloud API
	url := fmt.Sprintf("https://api.arubacloud.com/projects/%s/providers/Aruba.Network/vpcs/%s/subnets/%s?api-version=%s", projectId, vpcId, id, apiVersion)

	// Make the PUT request to Aruba Cloud API
	resp, err := h.makeArubaCloudRequest("PUT", url, authHeader, arubaRequestBody)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to make update subnet request: %v", err))
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to read update subnet response")
		return
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		h.Log.Printf("Aruba Cloud API returned non-200 status for update subnet: %d. Body: %s", resp.StatusCode, string(respBody))
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
		return
	}

	// Unmarshal the response into the Go struct to validate it
	var arubaResponse SubnetResponseDto
	if err := json.Unmarshal(respBody, &arubaResponse); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal Aruba Cloud response: %v", err))
		return
	}

	// Marshal the validated struct back to JSON to prepare for flattening
	validatedBody, err := json.Marshal(arubaResponse)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to marshal validated response: %v", err))
		return
	}

	// Flatten the response
	flattenedBody, err := utils.FlattenObject(validatedBody, "metadata")
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to flatten response: %v", err))
		return
	}

	h.writeJSONResponse(w, http.StatusOK, flattenedBody)
	h.Log.Printf("Successfully updated subnet '%s'", id)
}

// LIST handler implementation
// @Summary List Subnets on Aruba Cloud
// @Description List Subnets on Aruba Cloud using the provided project and vpc details.
// @ID list-subnets
// @Param projectId path string true "Project ID"
// @Param vpcId path string true "VPC ID"
// @Param api-version query string true "API version (e.g., 1.0)"
// @Param filter query string false "Filter expression"
// @Param sort query string false "Sort expression"
// @Param projection query string false "Projection expression"
// @Param offset query integer false "Offset for pagination"
// @Param limit query integer false "Limit for pagination"
// @Param Authorization header string true "Bearer Token (Bearer <token>)"
// @Accept json
// @Produce json
// @Success 200 {object} FlattenedSubnetListResponseDto "A list of subnets"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /projects/{projectId}/providers/Aruba.Network/vpcs/{vpcId}/subnets [get]
func (h *listHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("projectId")
	vpcId := r.PathValue("vpcId")
	authHeader := r.Header.Get("Authorization")

	// Validate required parameters
	if projectId == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Project ID parameter is required")
		return
	}
	if vpcId == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "VPC ID parameter is required")
		return
	}
	queryParams := r.URL.Query()
	if apiVersion := queryParams.Get("api-version"); apiVersion == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "API version parameter is required")
		return
	}
	if authHeader == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	// Construct the URL for the Aruba Cloud API
	baseURL := fmt.Sprintf("https://api.arubacloud.com/projects/%s/providers/Aruba.Network/vpcs/%s/subnets", projectId, vpcId)
	url := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	// Make the GET request to Aruba Cloud API
	resp, err := h.makeArubaCloudRequest("GET", url, authHeader, nil)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to make list subnets request: %v", err))
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to read list subnets response")
		return
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		h.Log.Printf("Aruba Cloud API returned non-200 status for list subnets: %d. Body: %s", resp.StatusCode, string(body))
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
		return
	}

	// Unmarshal the response into the Go struct to validate it
	var arubaResponse SubnetListResponseDto
	if err := json.Unmarshal(body, &arubaResponse); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal Aruba Cloud response: %v", err))
		return
	}

	// Flatten each subnet in the response
	flattenedValues := make([]FlattenedSubnetResponseDto, len(arubaResponse.Values))
	for i, subnet := range arubaResponse.Values {
		subnetBody, err := json.Marshal(subnet)
		if err != nil {
			h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to marshal subnet for flattening: %v", err))
			return
		}

		flattenedSubnetBody, err := utils.FlattenObject(subnetBody, "metadata")
		if err != nil {
			h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to flatten subnet: %v", err))
			return
		}

		var flattenedSubnet FlattenedSubnetResponseDto
		if err := json.Unmarshal(flattenedSubnetBody, &flattenedSubnet); err != nil {
			h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal flattened subnet: %v", err))
			return
		}
		flattenedValues[i] = flattenedSubnet
	}

	// Construct the flattened list response
	flattenedResponse := FlattenedSubnetListResponseDto{
		Total:  arubaResponse.Total,
		Self:   arubaResponse.Self,
		Prev:   arubaResponse.Prev,
		Next:   arubaResponse.Next,
		First:  arubaResponse.First,
		Last:   arubaResponse.Last,
		Values: flattenedValues,
	}

	finalBody, err := json.Marshal(flattenedResponse)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to marshal flattened list response: %v", err))
		return
	}

	h.writeJSONResponse(w, http.StatusOK, finalBody)
	h.Log.Printf("Successfully listed subnets for project '%s', vpc '%s'", projectId, vpcId)
}
