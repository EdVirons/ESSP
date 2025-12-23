package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool

	incidents                *IncidentRepo
	workOrders               *WorkOrderRepo
	attachments              *AttachmentRepo
	schools                  *SchoolRepo
	shops                    *ServiceShopRepo
	staff                    *ServiceStaffRepo
	parts                    *PartRepo
	inventory                *InventoryRepo
	workOrderParts           *WorkOrderPartRepo
	schoolsSnap              *SchoolsSnapshotRepo
	devicesSnap              *DevicesSnapshotRepo
	partsSnap                *PartsSnapshotRepo
	ssotState                *SSOTStateRepo
	projectsRepo             *ProjectsRepo
	phasesRepo               *PhasesRepo
	surveysRepo              *SurveysRepo
	surveyRoomsRepo          *SurveyRoomsRepo
	surveyPhotosRepo         *SurveyPhotosRepo
	boqRepo                  *BOQRepo
	contactsRepo             *SchoolContactsRepo
	scheduleRepo             *WorkOrderScheduleRepo
	deliverablesRepo         *WorkOrderDeliverablesRepo
	approvalsRepo            *WorkOrderApprovalsRepo
	phaseChecklistsRepo      *PhaseChecklistsRepo
	auditStore               *AuditStoreRef
	messagingRepo            *MessagingRepo
	chatSessionsRepo         *ChatSessionsRepo
	projectTeamRepo          *ProjectTeamRepo
	projectActivitiesRepo    *ProjectActivitiesRepo
	userNotificationsRepo    *UserNotificationsRepo
	workOrderReworkRepo      *WorkOrderReworkRepo
	bulkOperationRepo        *BulkOperationRepo
	featureConfigRepo        *FeatureConfigRepo
	notificationPrefsRepo    *NotificationPrefsRepo
	edtechProfilesRepo       *EdTechProfilesRepo
	edtechProfileHistoryRepo *EdTechProfileHistoryRepo
	demoLeadsRepo            *DemoLeadsRepo
	demoLeadActivitiesRepo   *DemoLeadActivitiesRepo
	demoSchedulesRepo        *DemoSchedulesRepo
	presentationsRepo        *PresentationsRepo
	presentationViewsRepo    *PresentationViewsRepo
	salesMetricsDailyRepo    *SalesMetricsDailyRepo
	ssotLocationsRepo        *SSOTLocationsRepo
	kbArticlesRepo           *KBArticleRepo
	marketingKBRepo          *MarketingKBRepo

	// Device inventory
	locationsRepo   *LocationsRepo
	assignmentsRepo *AssignmentsRepo
	groupsRepo      *GroupsRepo
	networkSnapRepo *NetworkSnapshotRepo

	// HR SSOT snapshots
	peopleSnap          *PeopleSnapshotRepo
	teamsSnap           *TeamsSnapshotRepo
	orgUnitsSnap        *OrgUnitsSnapshotRepo
	teamMembershipsSnap *TeamMembershipsSnapshotRepo
}

// AuditStoreRef is a placeholder for the audit store to avoid circular dependency
type AuditStoreRef struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, dsn string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	s := &Postgres{pool: pool}
	s.incidents = &IncidentRepo{pool: pool}
	s.workOrders = &WorkOrderRepo{pool: pool}
	s.attachments = &AttachmentRepo{pool: pool}
	s.schools = &SchoolRepo{pool: pool}
	s.shops = &ServiceShopRepo{pool: pool}
	s.staff = &ServiceStaffRepo{pool: pool}
	s.parts = &PartRepo{pool: pool}
	s.inventory = &InventoryRepo{pool: pool}
	s.workOrderParts = &WorkOrderPartRepo{pool: pool}
	s.schoolsSnap = &SchoolsSnapshotRepo{pool: pool}
	s.devicesSnap = &DevicesSnapshotRepo{pool: pool}
	s.partsSnap = &PartsSnapshotRepo{pool: pool}
	s.ssotState = &SSOTStateRepo{pool: pool}
	s.projectsRepo = &ProjectsRepo{pool: pool}
	s.phasesRepo = &PhasesRepo{pool: pool}
	s.surveysRepo = &SurveysRepo{pool: pool}
	s.surveyRoomsRepo = &SurveyRoomsRepo{pool: pool}
	s.surveyPhotosRepo = &SurveyPhotosRepo{pool: pool}
	s.boqRepo = &BOQRepo{pool: pool}
	s.contactsRepo = &SchoolContactsRepo{pool: pool}
	s.scheduleRepo = &WorkOrderScheduleRepo{pool: pool}
	s.deliverablesRepo = &WorkOrderDeliverablesRepo{pool: pool}
	s.approvalsRepo = &WorkOrderApprovalsRepo{pool: pool}
	s.phaseChecklistsRepo = &PhaseChecklistsRepo{pool: pool}
	s.auditStore = &AuditStoreRef{pool: pool}
	s.messagingRepo = &MessagingRepo{pool: pool}
	s.chatSessionsRepo = &ChatSessionsRepo{pool: pool}
	s.projectTeamRepo = &ProjectTeamRepo{pool: pool}
	s.projectActivitiesRepo = &ProjectActivitiesRepo{pool: pool}
	s.userNotificationsRepo = &UserNotificationsRepo{pool: pool}
	s.workOrderReworkRepo = &WorkOrderReworkRepo{pool: pool}
	s.bulkOperationRepo = &BulkOperationRepo{pool: pool}
	s.featureConfigRepo = &FeatureConfigRepo{pool: pool}
	s.notificationPrefsRepo = &NotificationPrefsRepo{pool: pool}
	s.edtechProfilesRepo = &EdTechProfilesRepo{pool: pool}
	s.edtechProfileHistoryRepo = &EdTechProfileHistoryRepo{pool: pool}
	s.demoLeadsRepo = &DemoLeadsRepo{pool: pool}
	s.demoLeadActivitiesRepo = &DemoLeadActivitiesRepo{pool: pool}
	s.demoSchedulesRepo = &DemoSchedulesRepo{pool: pool}
	s.presentationsRepo = &PresentationsRepo{pool: pool}
	s.presentationViewsRepo = &PresentationViewsRepo{pool: pool}
	s.salesMetricsDailyRepo = &SalesMetricsDailyRepo{pool: pool}
	s.ssotLocationsRepo = &SSOTLocationsRepo{pool: pool}
	s.kbArticlesRepo = &KBArticleRepo{pool: pool}
	s.marketingKBRepo = &MarketingKBRepo{pool: pool}

	// Device inventory
	s.locationsRepo = &LocationsRepo{pool: pool}
	s.assignmentsRepo = &AssignmentsRepo{pool: pool}
	s.groupsRepo = &GroupsRepo{pool: pool}
	s.networkSnapRepo = &NetworkSnapshotRepo{pool: pool}

	// HR SSOT snapshots
	s.peopleSnap = &PeopleSnapshotRepo{pool: pool}
	s.teamsSnap = &TeamsSnapshotRepo{pool: pool}
	s.orgUnitsSnap = &OrgUnitsSnapshotRepo{pool: pool}
	s.teamMembershipsSnap = &TeamMembershipsSnapshotRepo{pool: pool}
	return s, nil
}

func (p *Postgres) Close()                         { p.pool.Close() }
func (p *Postgres) Ping(ctx context.Context) error { return p.pool.Ping(ctx) }
func (p *Postgres) RawPool() *pgxpool.Pool         { return p.pool }

func (p *Postgres) Incidents() *IncidentRepo                          { return p.incidents }
func (p *Postgres) WorkOrders() *WorkOrderRepo                        { return p.workOrders }
func (p *Postgres) Attachments() *AttachmentRepo                      { return p.attachments }
func (p *Postgres) Schools() *SchoolRepo                              { return p.schools }
func (p *Postgres) ServiceShops() *ServiceShopRepo                    { return p.shops }
func (p *Postgres) ServiceStaff() *ServiceStaffRepo                   { return p.staff }
func (p *Postgres) Parts() *PartRepo                                  { return p.parts }
func (p *Postgres) Inventory() *InventoryRepo                         { return p.inventory }
func (p *Postgres) WorkOrderParts() *WorkOrderPartRepo                { return p.workOrderParts }
func (p *Postgres) SchoolsSnapshot() *SchoolsSnapshotRepo             { return p.schoolsSnap }
func (p *Postgres) DevicesSnapshot() *DevicesSnapshotRepo             { return p.devicesSnap }
func (p *Postgres) PartsSnapshot() *PartsSnapshotRepo                 { return p.partsSnap }
func (p *Postgres) SSOTState() *SSOTStateRepo                         { return p.ssotState }
func (p *Postgres) Projects() *ProjectsRepo                           { return p.projectsRepo }
func (p *Postgres) Phases() *PhasesRepo                               { return p.phasesRepo }
func (p *Postgres) Surveys() *SurveysRepo                             { return p.surveysRepo }
func (p *Postgres) SurveyRooms() *SurveyRoomsRepo                     { return p.surveyRoomsRepo }
func (p *Postgres) SurveyPhotos() *SurveyPhotosRepo                   { return p.surveyPhotosRepo }
func (p *Postgres) BOQ() *BOQRepo                                     { return p.boqRepo }
func (p *Postgres) SchoolContacts() *SchoolContactsRepo               { return p.contactsRepo }
func (p *Postgres) WorkOrderSchedules() *WorkOrderScheduleRepo        { return p.scheduleRepo }
func (p *Postgres) WorkOrderDeliverables() *WorkOrderDeliverablesRepo { return p.deliverablesRepo }
func (p *Postgres) WorkOrderApprovals() *WorkOrderApprovalsRepo       { return p.approvalsRepo }
func (p *Postgres) PhaseChecklists() *PhaseChecklistsRepo             { return p.phaseChecklistsRepo }
func (p *Postgres) AuditStorePool() *pgxpool.Pool                     { return p.auditStore.pool }
func (p *Postgres) Messaging() *MessagingRepo                         { return p.messagingRepo }
func (p *Postgres) ChatSessions() *ChatSessionsRepo                   { return p.chatSessionsRepo }
func (p *Postgres) ProjectTeam() *ProjectTeamRepo                     { return p.projectTeamRepo }
func (p *Postgres) ProjectActivities() *ProjectActivitiesRepo         { return p.projectActivitiesRepo }
func (p *Postgres) UserNotifications() *UserNotificationsRepo         { return p.userNotificationsRepo }
func (p *Postgres) WorkOrderRework() *WorkOrderReworkRepo             { return p.workOrderReworkRepo }
func (p *Postgres) BulkOperations() *BulkOperationRepo                { return p.bulkOperationRepo }
func (p *Postgres) FeatureConfig() *FeatureConfigRepo                 { return p.featureConfigRepo }
func (p *Postgres) NotificationPrefs() *NotificationPrefsRepo         { return p.notificationPrefsRepo }
func (p *Postgres) EdTechProfiles() *EdTechProfilesRepo               { return p.edtechProfilesRepo }
func (p *Postgres) EdTechProfileHistory() *EdTechProfileHistoryRepo {
	return p.edtechProfileHistoryRepo
}
func (p *Postgres) DemoLeads() *DemoLeadsRepo                   { return p.demoLeadsRepo }
func (p *Postgres) DemoLeadActivities() *DemoLeadActivitiesRepo { return p.demoLeadActivitiesRepo }
func (p *Postgres) DemoSchedules() *DemoSchedulesRepo           { return p.demoSchedulesRepo }
func (p *Postgres) Presentations() *PresentationsRepo           { return p.presentationsRepo }
func (p *Postgres) PresentationViews() *PresentationViewsRepo   { return p.presentationViewsRepo }
func (p *Postgres) SalesMetricsDaily() *SalesMetricsDailyRepo   { return p.salesMetricsDailyRepo }
func (p *Postgres) SSOTLocations() *SSOTLocationsRepo           { return p.ssotLocationsRepo }
func (p *Postgres) KBArticles() *KBArticleRepo                  { return p.kbArticlesRepo }
func (p *Postgres) MarketingKB() *MarketingKBRepo               { return p.marketingKBRepo }

// Device inventory
func (p *Postgres) Locations() *LocationsRepo             { return p.locationsRepo }
func (p *Postgres) Assignments() *AssignmentsRepo         { return p.assignmentsRepo }
func (p *Postgres) Groups() *GroupsRepo                   { return p.groupsRepo }
func (p *Postgres) NetworkSnapshot() *NetworkSnapshotRepo { return p.networkSnapRepo }

// HR SSOT snapshots
func (p *Postgres) PeopleSnapshot() *PeopleSnapshotRepo     { return p.peopleSnap }
func (p *Postgres) TeamsSnapshot() *TeamsSnapshotRepo       { return p.teamsSnap }
func (p *Postgres) OrgUnitsSnapshot() *OrgUnitsSnapshotRepo { return p.orgUnitsSnap }
func (p *Postgres) TeamMembershipsSnapshot() *TeamMembershipsSnapshotRepo {
	return p.teamMembershipsSnap
}
