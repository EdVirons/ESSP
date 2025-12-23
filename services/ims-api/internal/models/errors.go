package models

import "errors"

// Bulk operation errors
var (
	ErrBulkNoIDs           = errors.New("no work order IDs provided")
	ErrBulkTooManyIDs      = errors.New("too many work order IDs in batch")
	ErrBulkNoStatus        = errors.New("status is required for bulk status update")
	ErrBulkNoAssignment    = errors.New("at least one of assignedStaffId or serviceShopId is required")
	ErrBulkInvalidDecision = errors.New("decision must be 'approved' or 'rejected'")
)

// Rework errors
var (
	ErrReworkInvalidTransition = errors.New("invalid rework transition")
	ErrReworkMaxExceeded       = errors.New("maximum rework count exceeded")
	ErrReworkReasonRequired    = errors.New("rejection reason is required")
	ErrReworkInvalidCategory   = errors.New("invalid rejection category")
)

// Feature flag errors
var (
	ErrFeatureDisabled = errors.New("feature is disabled")
	ErrFeatureNotFound = errors.New("feature configuration not found")
)

// Work order update errors
var (
	ErrUpdateNoFields = errors.New("no fields to update")
	ErrUpdateNotFound = errors.New("work order not found")
)
