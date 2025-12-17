package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/edvirons/ssp/ims/internal/audit"
	"github.com/edvirons/ssp/ims/internal/middleware"
	"github.com/edvirons/ssp/ims/internal/models"
	"github.com/edvirons/ssp/ims/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type DeviceInventoryHandler struct {
	log   *zap.Logger
	pg    *store.Postgres
	audit audit.AuditLogger
}

func NewDeviceInventoryHandler(log *zap.Logger, pg *store.Postgres, auditLogger audit.AuditLogger) *DeviceInventoryHandler {
	return &DeviceInventoryHandler{log: log, pg: pg, audit: auditLogger}
}

// ---------- School Inventory ----------

// GetSchoolInventory returns the device inventory for a school
// GET /v1/schools/{schoolId}/inventory
func (h *DeviceInventoryHandler) GetSchoolInventory(w http.ResponseWriter, r *http.Request) {
	schoolID := chi.URLParam(r, "schoolId")
	if schoolID == "" {
		http.Error(w, "schoolId required", http.StatusBadRequest)
		return
	}
	tenant := middleware.TenantID(r.Context())

	// Get school info from snapshot
	school, err := h.pg.SchoolsSnapshot().Get(r.Context(), tenant, schoolID)
	if err != nil {
		h.log.Error("failed to get school", zap.Error(err))
		http.Error(w, "school not found", http.StatusNotFound)
		return
	}

	// Get devices for this school from snapshot
	devices, err := h.pg.DevicesSnapshot().ListBySchool(r.Context(), tenant, schoolID)
	if err != nil {
		h.log.Error("failed to get devices", zap.Error(err))
		http.Error(w, "failed to get devices", http.StatusInternalServerError)
		return
	}

	// Get locations for this school
	locations, err := h.pg.Locations().ListBySchool(r.Context(), tenant, schoolID)
	if err != nil {
		h.log.Error("failed to get locations", zap.Error(err))
		// Continue without locations
		locations = []models.Location{}
	}

	// Get device counts by location
	locDeviceCounts, _ := h.pg.Locations().CountDevicesByLocation(r.Context(), tenant, schoolID)

	// Build location map for quick lookup
	locationMap := make(map[string]models.Location)
	for i := range locations {
		locations[i].DeviceCount = locDeviceCounts[locations[i].ID]
		locationMap[locations[i].ID] = locations[i]
	}

	// Get device IDs for batch operations
	deviceIDs := make([]string, len(devices))
	for i, d := range devices {
		deviceIDs[i] = d.DeviceID
	}

	// Get current assignments for all devices
	assignmentMap, _ := h.pg.Assignments().GetCurrentAssignmentMap(r.Context(), tenant, deviceIDs)

	// Get MAC addresses for all devices
	macMap, _ := h.pg.NetworkSnapshot().GetMACAddressesForDevices(r.Context(), tenant, deviceIDs)

	// Build inventory devices with joined data
	inventoryDevices := make([]models.InventoryDevice, 0, len(devices))
	summary := models.InventorySummary{
		ByStatus:   make(map[string]int),
		ByLocation: map[string]int{"assigned": 0, "unassigned": 0},
		ByModel:    make(map[string]int),
	}

	for _, d := range devices {
		// Map DeviceSnapshot to InventoryDevice
		inv := models.InventoryDevice{
			ID:        d.DeviceID,
			TenantID:  d.TenantID,
			Serial:    d.Serial,
			AssetTag:  d.AssetTag,
			Model:     d.Model,
			SchoolID:  d.SchoolID,
			Lifecycle: d.Status, // Status maps to lifecycle
			UpdatedAt: d.UpdatedAt,
		}

		// Extract make from model (first word)
		if d.Model != "" {
			parts := strings.SplitN(d.Model, " ", 2)
			if len(parts) > 0 {
				inv.Make = parts[0]
			}
		}

		// Add assignment/location info
		if assignment, ok := assignmentMap[d.DeviceID]; ok && assignment.LocationID != nil {
			if loc, ok := locationMap[*assignment.LocationID]; ok {
				inv.Location = &loc
				// Get path
				if path, err := h.pg.Locations().GetPath(r.Context(), tenant, loc.ID); err == nil {
					inv.LocationPath = path
				}
			}
			summary.ByLocation["assigned"]++
		} else {
			summary.ByLocation["unassigned"]++
		}

		// Add MAC addresses
		if macs, ok := macMap[d.DeviceID]; ok {
			inv.MACAddresses = macs
		}

		// Update summary
		status := d.Status
		if status == "" {
			status = "unknown"
		}
		summary.ByStatus[status]++
		if inv.Model != "" {
			summary.ByModel[inv.Model]++
		}

		inventoryDevices = append(inventoryDevices, inv)
	}

	summary.TotalDevices = len(inventoryDevices)

	// Build response
	resp := models.SchoolInventoryResponse{
		Summary:   summary,
		Devices:   inventoryDevices,
		Locations: locations,
	}
	resp.School.ID = school.SchoolID
	resp.School.Name = school.Name

	writeJSON(w, http.StatusOK, resp)
}

// ---------- Locations ----------

type createLocationReq struct {
	ParentID     *string `json:"parentId"`
	Name         string  `json:"name"`
	LocationType string  `json:"locationType"`
	Code         string  `json:"code"`
	Capacity     int     `json:"capacity"`
}

// CreateLocation creates a new location
// POST /v1/schools/{schoolId}/locations
func (h *DeviceInventoryHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	schoolID := chi.URLParam(r, "schoolId")
	tenant := middleware.TenantID(r.Context())

	var req createLocationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	locType := models.LocationType(strings.TrimSpace(req.LocationType))
	if locType == "" {
		locType = models.LocationTypeRoom
	}

	loc := models.Location{
		ID:           store.NewID("loc"),
		TenantID:     tenant,
		SchoolID:     schoolID,
		ParentID:     req.ParentID,
		Name:         strings.TrimSpace(req.Name),
		LocationType: locType,
		Code:         strings.TrimSpace(req.Code),
		Capacity:     req.Capacity,
		Active:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.pg.Locations().Create(r.Context(), loc); err != nil {
		h.log.Error("failed to create location", zap.Error(err))
		http.Error(w, "failed to create location", http.StatusInternalServerError)
		return
	}

	// Audit log
	if err := h.audit.LogCreate(r.Context(), "location", loc.ID, loc); err != nil {
		h.log.Error("failed to log location creation audit", zap.Error(err))
	}

	writeJSON(w, http.StatusCreated, loc)
}

// ListLocations returns locations for a school
// GET /v1/schools/{schoolId}/locations
func (h *DeviceInventoryHandler) ListLocations(w http.ResponseWriter, r *http.Request) {
	schoolID := chi.URLParam(r, "schoolId")
	tenant := middleware.TenantID(r.Context())

	tree := r.URL.Query().Get("tree") == "true"

	if tree {
		// Return as tree structure
		nodes, err := h.pg.Locations().GetTree(r.Context(), tenant, schoolID)
		if err != nil {
			h.log.Error("failed to get location tree", zap.Error(err))
			http.Error(w, "failed to get locations", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"tree": nodes})
		return
	}

	// Return flat list
	locations, err := h.pg.Locations().ListBySchool(r.Context(), tenant, schoolID)
	if err != nil {
		h.log.Error("failed to get locations", zap.Error(err))
		http.Error(w, "failed to get locations", http.StatusInternalServerError)
		return
	}

	// Get device counts
	counts, _ := h.pg.Locations().CountDevicesByLocation(r.Context(), tenant, schoolID)
	for i := range locations {
		locations[i].DeviceCount = counts[locations[i].ID]
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": locations})
}

type updateLocationReq struct {
	ParentID     *string `json:"parentId"`
	Name         string  `json:"name"`
	LocationType string  `json:"locationType"`
	Code         string  `json:"code"`
	Capacity     int     `json:"capacity"`
	Active       *bool   `json:"active"`
}

// UpdateLocation updates a location
// PUT /v1/schools/{schoolId}/locations/{id}
func (h *DeviceInventoryHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	schoolID := chi.URLParam(r, "schoolId")
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	// Get existing
	loc, err := h.pg.Locations().Get(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if loc.SchoolID != schoolID {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Capture before state for audit
	before := loc

	var req updateLocationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Update fields
	if req.Name != "" {
		loc.Name = strings.TrimSpace(req.Name)
	}
	if req.LocationType != "" {
		loc.LocationType = models.LocationType(strings.TrimSpace(req.LocationType))
	}
	if req.Code != "" {
		loc.Code = strings.TrimSpace(req.Code)
	}
	if req.Capacity > 0 {
		loc.Capacity = req.Capacity
	}
	if req.ParentID != nil {
		loc.ParentID = req.ParentID
	}
	if req.Active != nil {
		loc.Active = *req.Active
	}
	loc.UpdatedAt = time.Now().UTC()

	if err := h.pg.Locations().Update(r.Context(), loc); err != nil {
		h.log.Error("failed to update location", zap.Error(err))
		http.Error(w, "failed to update location", http.StatusInternalServerError)
		return
	}

	// Audit log
	if err := h.audit.LogUpdate(r.Context(), "location", id, before, loc); err != nil {
		h.log.Error("failed to log location update audit", zap.Error(err))
	}

	writeJSON(w, http.StatusOK, loc)
}

// DeleteLocation deletes a location
// DELETE /v1/schools/{schoolId}/locations/{id}
func (h *DeviceInventoryHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	// Get existing for audit log
	loc, err := h.pg.Locations().Get(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err := h.pg.Locations().Delete(r.Context(), tenant, id); err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Audit log
	if err := h.audit.LogDelete(r.Context(), "location", id, loc); err != nil {
		h.log.Error("failed to log location deletion audit", zap.Error(err))
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// ---------- Device Assignments ----------

type assignDeviceReq struct {
	LocationID     string `json:"locationId"`
	AssignedToUser string `json:"assignedToUser"`
	AssignmentType string `json:"assignmentType"`
	Notes          string `json:"notes"`
}

// AssignDevice assigns a device to a location
// POST /v1/devices/{deviceId}/assign
func (h *DeviceInventoryHandler) AssignDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceId")
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	var req assignDeviceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	asnType := models.AssignmentType(strings.TrimSpace(req.AssignmentType))
	if asnType == "" {
		asnType = models.AssignmentTypePermanent
	}

	var locID *string
	if req.LocationID != "" {
		locID = &req.LocationID
	}

	assignment := models.DeviceAssignment{
		ID:             store.NewID("asn"),
		TenantID:       tenant,
		DeviceID:       deviceID,
		LocationID:     locID,
		AssignedToUser: strings.TrimSpace(req.AssignedToUser),
		AssignmentType: asnType,
		Notes:          strings.TrimSpace(req.Notes),
		CreatedBy:      userID,
	}

	if err := h.pg.Assignments().AssignDevice(r.Context(), assignment); err != nil {
		h.log.Error("failed to assign device", zap.Error(err))
		http.Error(w, "failed to assign device", http.StatusInternalServerError)
		return
	}

	// Audit log
	if err := h.audit.LogCreate(r.Context(), "device_assignment", assignment.ID, assignment); err != nil {
		h.log.Error("failed to log device assignment audit", zap.Error(err))
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "assignment": assignment})
}

// UnassignDevice removes the current assignment from a device
// DELETE /v1/devices/{deviceId}/assign
func (h *DeviceInventoryHandler) UnassignDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceId")
	tenant := middleware.TenantID(r.Context())

	// Get current assignment for audit log
	current, currentErr := h.pg.Assignments().GetCurrent(r.Context(), tenant, deviceID)

	if err := h.pg.Assignments().UnassignDevice(r.Context(), tenant, deviceID); err != nil {
		h.log.Error("failed to unassign device", zap.Error(err))
		http.Error(w, "no current assignment", http.StatusNotFound)
		return
	}

	// Audit log
	if currentErr == nil && current.ID != "" {
		if err := h.audit.LogUpdate(r.Context(), "device_assignment", current.ID, current, map[string]any{
			"effective_to": time.Now().UTC(),
			"action":       "unassign",
		}); err != nil {
			h.log.Error("failed to log device unassignment audit", zap.Error(err))
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// GetDeviceAssignments returns assignment history for a device
// GET /v1/devices/{deviceId}/assignments
func (h *DeviceInventoryHandler) GetDeviceAssignments(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceId")
	tenant := middleware.TenantID(r.Context())

	items, err := h.pg.Assignments().GetHistory(r.Context(), tenant, deviceID, 20)
	if err != nil {
		h.log.Error("failed to get assignments", zap.Error(err))
		http.Error(w, "failed to get assignments", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetLocationDevices returns devices at a location
// GET /v1/locations/{locationId}/devices
func (h *DeviceInventoryHandler) GetLocationDevices(w http.ResponseWriter, r *http.Request) {
	locationID := chi.URLParam(r, "locationId")
	tenant := middleware.TenantID(r.Context())

	assignments, err := h.pg.Assignments().ListByLocation(r.Context(), tenant, locationID)
	if err != nil {
		h.log.Error("failed to get assignments", zap.Error(err))
		http.Error(w, "failed to get assignments", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": assignments})
}

// ---------- Device Groups ----------

type createGroupReq struct {
	SchoolID    *string               `json:"schoolId"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	GroupType   string                `json:"groupType"`
	LocationID  *string               `json:"locationId"`
	Selector    *models.GroupSelector `json:"selector"`
}

// CreateGroup creates a new device group
// POST /v1/groups
func (h *DeviceInventoryHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	tenant := middleware.TenantID(r.Context())

	var req createGroupReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	gType := models.GroupType(strings.TrimSpace(req.GroupType))
	if gType == "" {
		gType = models.GroupTypeManual
	}

	now := time.Now().UTC()
	group := models.DeviceGroup{
		ID:          store.NewID("grp"),
		TenantID:    tenant,
		SchoolID:    req.SchoolID,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		GroupType:   gType,
		LocationID:  req.LocationID,
		Selector:    req.Selector,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.pg.Groups().Create(r.Context(), group); err != nil {
		h.log.Error("failed to create group", zap.Error(err))
		http.Error(w, "failed to create group", http.StatusInternalServerError)
		return
	}

	// Audit log
	if err := h.audit.LogCreate(r.Context(), "device_group", group.ID, group); err != nil {
		h.log.Error("failed to log group creation audit", zap.Error(err))
	}

	writeJSON(w, http.StatusCreated, group)
}

// GetGroup returns a group by ID
// GET /v1/groups/{id}
func (h *DeviceInventoryHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	group, err := h.pg.Groups().Get(r.Context(), tenant, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Get member count
	group.MemberCount, _ = h.pg.Groups().GetMemberCount(r.Context(), tenant, id)

	writeJSON(w, http.StatusOK, group)
}

// ListGroups returns groups for a school
// GET /v1/schools/{schoolId}/groups
func (h *DeviceInventoryHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	schoolID := chi.URLParam(r, "schoolId")
	tenant := middleware.TenantID(r.Context())

	groups, err := h.pg.Groups().ListBySchool(r.Context(), tenant, schoolID)
	if err != nil {
		h.log.Error("failed to get groups", zap.Error(err))
		http.Error(w, "failed to get groups", http.StatusInternalServerError)
		return
	}

	// Get member counts
	for i := range groups {
		groups[i].MemberCount, _ = h.pg.Groups().GetMemberCount(r.Context(), tenant, groups[i].ID)
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": groups})
}

type addGroupMembersReq struct {
	DeviceIDs []string `json:"deviceIds"`
}

// AddGroupMembers adds devices to a group
// POST /v1/groups/{id}/members
func (h *DeviceInventoryHandler) AddGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	var req addGroupMembersReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if len(req.DeviceIDs) == 0 {
		http.Error(w, "deviceIds required", http.StatusBadRequest)
		return
	}

	count, err := h.pg.Groups().BulkAddMembers(r.Context(), tenant, groupID, req.DeviceIDs, userID)
	if err != nil {
		h.log.Error("failed to add members", zap.Error(err))
		http.Error(w, "failed to add members", http.StatusInternalServerError)
		return
	}

	// Audit log
	if err := h.audit.LogUpdate(r.Context(), "device_group", groupID, nil, map[string]any{
		"action":    "add_members",
		"deviceIds": req.DeviceIDs,
		"added":     count,
	}); err != nil {
		h.log.Error("failed to log group members addition audit", zap.Error(err))
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "added": count})
}

type removeGroupMembersReq struct {
	DeviceIDs []string `json:"deviceIds"`
}

// RemoveGroupMembers removes devices from a group
// DELETE /v1/groups/{id}/members
func (h *DeviceInventoryHandler) RemoveGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	var req removeGroupMembersReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if len(req.DeviceIDs) == 0 {
		http.Error(w, "deviceIds required", http.StatusBadRequest)
		return
	}

	count, err := h.pg.Groups().BulkRemoveMembers(r.Context(), tenant, groupID, req.DeviceIDs)
	if err != nil {
		h.log.Error("failed to remove members", zap.Error(err))
		http.Error(w, "failed to remove members", http.StatusInternalServerError)
		return
	}

	// Audit log
	if err := h.audit.LogUpdate(r.Context(), "device_group", groupID, map[string]any{
		"deviceIds": req.DeviceIDs,
	}, map[string]any{
		"action":  "remove_members",
		"removed": count,
	}); err != nil {
		h.log.Error("failed to log group members removal audit", zap.Error(err))
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "removed": count})
}

// GetGroupDevices returns devices in a group
// GET /v1/groups/{id}/devices
func (h *DeviceInventoryHandler) GetGroupDevices(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	tenant := middleware.TenantID(r.Context())

	members, err := h.pg.Groups().ListMembers(r.Context(), tenant, groupID)
	if err != nil {
		h.log.Error("failed to get members", zap.Error(err))
		http.Error(w, "failed to get members", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": members})
}

// ---------- Device Registration ----------

type registerDeviceReq struct {
	Serial     string  `json:"serial"`
	AssetTag   string  `json:"assetTag"`
	Model      string  `json:"model"`
	Make       string  `json:"make"`
	Notes      string  `json:"notes"`
	LocationID *string `json:"locationId"`
}

// RegisterDevice registers a new device for a school
// POST /v1/schools/{schoolId}/devices
func (h *DeviceInventoryHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	schoolID := chi.URLParam(r, "schoolId")
	tenant := middleware.TenantID(r.Context())
	userID := middleware.UserID(r.Context())

	var req registerDeviceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Validate required fields
	serial := strings.TrimSpace(req.Serial)
	model := strings.TrimSpace(req.Model)
	if serial == "" {
		http.Error(w, "serial is required", http.StatusBadRequest)
		return
	}
	if model == "" {
		http.Error(w, "model is required", http.StatusBadRequest)
		return
	}

	// Create device ID
	deviceID := store.NewID("dev")
	now := time.Now().UTC()

	// Create device in snapshot (school-registered device)
	device := models.DeviceSnapshot{
		DeviceID:  deviceID,
		TenantID:  tenant,
		SchoolID:  schoolID,
		Serial:    serial,
		AssetTag:  strings.TrimSpace(req.AssetTag),
		Model:     model,
		Status:    "active",
		UpdatedAt: now,
	}

	if err := h.pg.DevicesSnapshot().Upsert(r.Context(), device); err != nil {
		h.log.Error("failed to register device", zap.Error(err))
		http.Error(w, "failed to register device", http.StatusInternalServerError)
		return
	}

	// Audit log device creation
	if err := h.audit.LogCreate(r.Context(), "device", deviceID, map[string]any{
		"serial":    serial,
		"assetTag":  req.AssetTag,
		"model":     model,
		"make":      req.Make,
		"schoolId":  schoolID,
		"source":    "school_registered",
		"createdBy": userID,
	}); err != nil {
		h.log.Error("failed to log device creation audit", zap.Error(err))
	}

	// If location provided, create assignment
	var assignment *models.DeviceAssignment
	if req.LocationID != nil && *req.LocationID != "" {
		assignment = &models.DeviceAssignment{
			ID:             store.NewID("asn"),
			TenantID:       tenant,
			DeviceID:       deviceID,
			LocationID:     req.LocationID,
			AssignmentType: models.AssignmentTypePermanent,
			Notes:          strings.TrimSpace(req.Notes),
			CreatedBy:      userID,
		}
		if err := h.pg.Assignments().AssignDevice(r.Context(), *assignment); err != nil {
			h.log.Error("failed to create initial assignment", zap.Error(err))
			// Don't fail - device was created successfully
		} else {
			// Audit log assignment
			if err := h.audit.LogCreate(r.Context(), "device_assignment", assignment.ID, assignment); err != nil {
				h.log.Error("failed to log assignment audit", zap.Error(err))
			}
		}
	}

	// Build response
	resp := models.InventoryDevice{
		ID:        deviceID,
		TenantID:  tenant,
		Serial:    serial,
		AssetTag:  strings.TrimSpace(req.AssetTag),
		Model:     model,
		Make:      strings.TrimSpace(req.Make),
		SchoolID:  schoolID,
		Lifecycle: "active",
		UpdatedAt: now,
	}

	// Add location if assigned
	if assignment != nil && req.LocationID != nil {
		loc, err := h.pg.Locations().Get(r.Context(), tenant, *req.LocationID)
		if err == nil {
			resp.Location = &loc
			if path, err := h.pg.Locations().GetPath(r.Context(), tenant, loc.ID); err == nil {
				resp.LocationPath = path
			}
		}
	}

	writeJSON(w, http.StatusCreated, resp)
}
