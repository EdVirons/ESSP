import * as React from 'react';
import { Search, Users, Building2, UsersRound, Mail, Phone, BadgeCheck, Loader2, Plus, RefreshCw, ChevronLeft, ChevronRight, FolderTree } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { usePeople, useTeams, useOrgUnits, useCreatePerson, useCreateTeam, useCreateOrgUnit, useUpdatePerson, useUpdateTeam, useUpdateOrgUnit, useDeletePerson, useDeleteTeam, useDeleteOrgUnit, useSyncAllHR } from '@/api/hr';
import { CreatePersonModal, CreateTeamModal, CreateOrgUnitModal, PersonDetailModal, TeamDetailModal, OrgUnitDetailModal, OrgTree } from '@/components/hr';
import { toast } from '@/lib/toast';
import type { PersonSnapshot, TeamSnapshot, OrgUnitSnapshot, CreatePersonInput, CreateTeamInput, CreateOrgUnitInput } from '@/types/hr';

const PAGE_SIZE = 20;

export function HR() {
  const [searchQuery, setSearchQuery] = React.useState('');
  const [activeTab, setActiveTab] = React.useState('people');
  const [showPersonModal, setShowPersonModal] = React.useState(false);
  const [showTeamModal, setShowTeamModal] = React.useState(false);
  const [showOrgUnitModal, setShowOrgUnitModal] = React.useState(false);
  const [selectedPerson, setSelectedPerson] = React.useState<PersonSnapshot | null>(null);
  const [selectedTeam, setSelectedTeam] = React.useState<TeamSnapshot | null>(null);
  const [selectedOrgUnit, setSelectedOrgUnit] = React.useState<OrgUnitSnapshot | null>(null);

  // Pagination state for each tab
  const [peoplePage, setPeoplePage] = React.useState(0);
  const [teamsPage, setTeamsPage] = React.useState(0);
  const [orgUnitsPage, setOrgUnitsPage] = React.useState(0);

  // Fetch data with pagination
  const { data: peopleData, isLoading: peopleLoading } = usePeople({ limit: PAGE_SIZE, offset: peoplePage * PAGE_SIZE });
  const { data: teamsData, isLoading: teamsLoading } = useTeams({ limit: PAGE_SIZE, offset: teamsPage * PAGE_SIZE });
  const { data: orgUnitsData, isLoading: orgUnitsLoading } = useOrgUnits({ limit: PAGE_SIZE, offset: orgUnitsPage * PAGE_SIZE });

  // Mutations
  const createPerson = useCreatePerson();
  const createTeam = useCreateTeam();
  const createOrgUnit = useCreateOrgUnit();
  const updatePerson = useUpdatePerson();
  const updateTeam = useUpdateTeam();
  const updateOrgUnit = useUpdateOrgUnit();
  const deletePerson = useDeletePerson();
  const deleteTeam = useDeleteTeam();
  const deleteOrgUnit = useDeleteOrgUnit();
  const syncAll = useSyncAllHR();

  // Sync handler
  const handleSync = async () => {
    try {
      await syncAll.mutateAsync();
      toast.success('HR data synced successfully');
    } catch (error) {
      toast.error('Failed to sync HR data', String(error));
    }
  };

  // Filter people by search
  const filteredPeople = React.useMemo(() => {
    const people = peopleData?.items || [];
    if (!searchQuery) return people;
    const query = searchQuery.toLowerCase();
    return people.filter(
      (p) =>
        p.fullName?.toLowerCase().includes(query) ||
        p.email?.toLowerCase().includes(query) ||
        p.title?.toLowerCase().includes(query)
    );
  }, [peopleData?.items, searchQuery]);

  // Filter teams by search
  const filteredTeams = React.useMemo(() => {
    const teams = teamsData?.items || [];
    if (!searchQuery) return teams;
    const query = searchQuery.toLowerCase();
    return teams.filter(
      (t) =>
        t.name?.toLowerCase().includes(query) ||
        t.key?.toLowerCase().includes(query) ||
        t.description?.toLowerCase().includes(query)
    );
  }, [teamsData?.items, searchQuery]);

  // Filter org units by search
  const filteredOrgUnits = React.useMemo(() => {
    const units = orgUnitsData?.items || [];
    if (!searchQuery) return units;
    const query = searchQuery.toLowerCase();
    return units.filter(
      (o) =>
        o.name?.toLowerCase().includes(query) ||
        o.code?.toLowerCase().includes(query)
    );
  }, [orgUnitsData?.items, searchQuery]);

  const stats = {
    people: peopleData?.total ?? peopleData?.items?.length ?? 0,
    teams: teamsData?.total ?? teamsData?.items?.length ?? 0,
    orgUnits: orgUnitsData?.total ?? orgUnitsData?.items?.length ?? 0,
  };

  // Pagination calculations
  const peopleTotalPages = Math.ceil(stats.people / PAGE_SIZE);
  const teamsTotalPages = Math.ceil(stats.teams / PAGE_SIZE);
  const orgUnitsTotalPages = Math.ceil(stats.orgUnits / PAGE_SIZE);

  const orgUnitsList = orgUnitsData?.items || [];

  // Create handlers
  const handleCreatePerson = async (data: Parameters<typeof createPerson.mutateAsync>[0]) => {
    try {
      await createPerson.mutateAsync(data);
      setShowPersonModal(false);
      toast.success('Person created successfully');
    } catch (error) {
      toast.error('Failed to create person', String(error));
    }
  };

  const handleCreateTeam = async (data: Parameters<typeof createTeam.mutateAsync>[0]) => {
    try {
      await createTeam.mutateAsync(data);
      setShowTeamModal(false);
      toast.success('Team created successfully');
    } catch (error) {
      toast.error('Failed to create team', String(error));
    }
  };

  const handleCreateOrgUnit = async (data: Parameters<typeof createOrgUnit.mutateAsync>[0]) => {
    try {
      await createOrgUnit.mutateAsync(data);
      setShowOrgUnitModal(false);
      toast.success('Org unit created successfully');
    } catch (error) {
      toast.error('Failed to create org unit', String(error));
    }
  };

  // Update handlers
  const handleUpdatePerson = async (id: string, data: Partial<CreatePersonInput>) => {
    try {
      await updatePerson.mutateAsync({ id, ...data });
      setSelectedPerson(null);
      toast.success('Person updated successfully');
    } catch (error) {
      toast.error('Failed to update person', String(error));
    }
  };

  const handleUpdateTeam = async (id: string, data: Partial<CreateTeamInput>) => {
    try {
      await updateTeam.mutateAsync({ id, ...data });
      setSelectedTeam(null);
      toast.success('Team updated successfully');
    } catch (error) {
      toast.error('Failed to update team', String(error));
    }
  };

  const handleUpdateOrgUnit = async (id: string, data: Partial<CreateOrgUnitInput>) => {
    try {
      await updateOrgUnit.mutateAsync({ id, ...data });
      setSelectedOrgUnit(null);
      toast.success('Org unit updated successfully');
    } catch (error) {
      toast.error('Failed to update org unit', String(error));
    }
  };

  // Delete handlers
  const handleDeletePerson = async (id: string) => {
    try {
      await deletePerson.mutateAsync(id);
      setSelectedPerson(null);
      toast.success('Person deleted successfully');
    } catch (error) {
      toast.error('Failed to delete person', String(error));
    }
  };

  const handleDeleteTeam = async (id: string) => {
    try {
      await deleteTeam.mutateAsync(id);
      setSelectedTeam(null);
      toast.success('Team deleted successfully');
    } catch (error) {
      toast.error('Failed to delete team', String(error));
    }
  };

  const handleDeleteOrgUnit = async (id: string) => {
    try {
      await deleteOrgUnit.mutateAsync(id);
      setSelectedOrgUnit(null);
      toast.success('Org unit deleted successfully');
    } catch (error) {
      toast.error('Failed to delete org unit', String(error));
    }
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex justify-between items-start">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">HR Directory</h1>
          <p className="text-sm text-gray-500">
            View people, teams, and organizational structure
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            onClick={handleSync}
            variant="outline"
            className="gap-2"
            disabled={syncAll.isPending}
          >
            <RefreshCw className={`h-4 w-4 ${syncAll.isPending ? 'animate-spin' : ''}`} />
            {syncAll.isPending ? 'Syncing...' : 'Sync'}
          </Button>
          <Button onClick={() => setShowPersonModal(true)} className="gap-2">
            <Plus className="h-4 w-4" />
            Add Person
          </Button>
          <Button onClick={() => setShowTeamModal(true)} variant="outline" className="gap-2">
            <Plus className="h-4 w-4" />
            Create Team
          </Button>
          <Button onClick={() => setShowOrgUnitModal(true)} variant="outline" className="gap-2">
            <Plus className="h-4 w-4" />
            Add Org Unit
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <div className="p-3 bg-violet-100 rounded-lg">
              <Users className="h-6 w-6 text-violet-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">People</p>
              <p className="text-2xl font-bold">
                {peopleLoading ? <Loader2 className="h-6 w-6 animate-spin" /> : stats.people}
              </p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <div className="p-3 bg-blue-100 rounded-lg">
              <UsersRound className="h-6 w-6 text-blue-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Teams</p>
              <p className="text-2xl font-bold">
                {teamsLoading ? <Loader2 className="h-6 w-6 animate-spin" /> : stats.teams}
              </p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 flex items-center gap-4">
            <div className="p-3 bg-green-100 rounded-lg">
              <Building2 className="h-6 w-6 text-green-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Org Units</p>
              <p className="text-2xl font-bold">
                {orgUnitsLoading ? <Loader2 className="h-6 w-6 animate-spin" /> : stats.orgUnits}
              </p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Search */}
      <div className="relative max-w-md">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <Input
          placeholder="Search people, teams, or org units..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-9"
        />
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="people" className="gap-2">
            <Users className="h-4 w-4" />
            People ({filteredPeople.length})
          </TabsTrigger>
          <TabsTrigger value="teams" className="gap-2">
            <UsersRound className="h-4 w-4" />
            Teams ({filteredTeams.length})
          </TabsTrigger>
          <TabsTrigger value="org" className="gap-2">
            <Building2 className="h-4 w-4" />
            Org Units ({filteredOrgUnits.length})
          </TabsTrigger>
          <TabsTrigger value="tree" className="gap-2">
            <FolderTree className="h-4 w-4" />
            Org Tree
          </TabsTrigger>
        </TabsList>

        {/* People Tab */}
        <TabsContent value="people" className="mt-4">
          {peopleLoading ? (
            <div className="flex items-center justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-violet-600" />
              <span className="ml-2 text-gray-500">Loading people...</span>
            </div>
          ) : filteredPeople.length === 0 ? (
            <Card>
              <CardContent className="p-8 text-center text-gray-500">
                No people found
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {filteredPeople.map((person) => (
                  <Card
                    key={person.personId}
                    className="hover:shadow-md transition-shadow cursor-pointer"
                    onClick={() => setSelectedPerson(person)}
                  >
                    <CardContent className="p-4">
                      <div className="flex items-start gap-3">
                        <div className="h-12 w-12 rounded-full bg-violet-100 flex items-center justify-center text-violet-600 font-semibold text-lg">
                          {person.givenName?.[0]}{person.familyName?.[0]}
                        </div>
                        <div className="flex-1 min-w-0">
                          <h3 className="font-semibold text-gray-900 truncate">
                            {person.fullName || `${person.givenName} ${person.familyName}`}
                          </h3>
                          {person.title && (
                            <p className="text-sm text-gray-500 truncate">{person.title}</p>
                          )}
                          <div className="mt-2 space-y-1">
                            {person.email && (
                              <div className="flex items-center gap-1 text-sm text-gray-600">
                                <Mail className="h-3 w-3" />
                                <span className="truncate">{person.email}</span>
                              </div>
                            )}
                            {person.phone && (
                              <div className="flex items-center gap-1 text-sm text-gray-600">
                                <Phone className="h-3 w-3" />
                                <span>{person.phone}</span>
                              </div>
                            )}
                          </div>
                          <div className="mt-2">
                            <Badge
                              variant={person.status === 'active' ? 'default' : 'secondary'}
                              className="text-xs"
                            >
                              <BadgeCheck className="h-3 w-3 mr-1" />
                              {person.status}
                            </Badge>
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
              {/* Pagination */}
              {peopleTotalPages > 1 && (
                <div className="flex items-center justify-between border-t pt-4">
                  <p className="text-sm text-gray-500">
                    Showing {peoplePage * PAGE_SIZE + 1} to {Math.min((peoplePage + 1) * PAGE_SIZE, stats.people)} of {stats.people}
                  </p>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setPeoplePage(p => Math.max(0, p - 1))}
                      disabled={peoplePage === 0}
                    >
                      <ChevronLeft className="h-4 w-4" />
                      Previous
                    </Button>
                    <span className="text-sm text-gray-600 min-w-[100px] text-center">
                      Page {peoplePage + 1} of {peopleTotalPages}
                    </span>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setPeoplePage(p => Math.min(peopleTotalPages - 1, p + 1))}
                      disabled={peoplePage >= peopleTotalPages - 1}
                    >
                      Next
                      <ChevronRight className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              )}
            </div>
          )}
        </TabsContent>

        {/* Teams Tab */}
        <TabsContent value="teams" className="mt-4">
          {teamsLoading ? (
            <div className="flex items-center justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
              <span className="ml-2 text-gray-500">Loading teams...</span>
            </div>
          ) : filteredTeams.length === 0 ? (
            <Card>
              <CardContent className="p-8 text-center text-gray-500">
                No teams found
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {filteredTeams.map((team) => (
                  <Card
                    key={team.teamId}
                    className="hover:shadow-md transition-shadow cursor-pointer"
                    onClick={() => setSelectedTeam(team)}
                  >
                    <CardHeader className="pb-2">
                      <div className="flex items-center gap-3">
                        <div className="p-2 bg-blue-100 rounded-lg">
                          <UsersRound className="h-5 w-5 text-blue-600" />
                        </div>
                        <div>
                          <CardTitle className="text-lg">{team.name}</CardTitle>
                          <Badge variant="outline" className="text-xs mt-1">
                            {team.key}
                          </Badge>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      {team.description && (
                        <p className="text-sm text-gray-600">{team.description}</p>
                      )}
                    </CardContent>
                  </Card>
                ))}
              </div>
              {/* Pagination */}
              {teamsTotalPages > 1 && (
                <div className="flex items-center justify-between border-t pt-4">
                  <p className="text-sm text-gray-500">
                    Showing {teamsPage * PAGE_SIZE + 1} to {Math.min((teamsPage + 1) * PAGE_SIZE, stats.teams)} of {stats.teams}
                  </p>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setTeamsPage(p => Math.max(0, p - 1))}
                      disabled={teamsPage === 0}
                    >
                      <ChevronLeft className="h-4 w-4" />
                      Previous
                    </Button>
                    <span className="text-sm text-gray-600 min-w-[100px] text-center">
                      Page {teamsPage + 1} of {teamsTotalPages}
                    </span>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setTeamsPage(p => Math.min(teamsTotalPages - 1, p + 1))}
                      disabled={teamsPage >= teamsTotalPages - 1}
                    >
                      Next
                      <ChevronRight className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              )}
            </div>
          )}
        </TabsContent>

        {/* Org Units Tab */}
        <TabsContent value="org" className="mt-4">
          {orgUnitsLoading ? (
            <div className="flex items-center justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-green-600" />
              <span className="ml-2 text-gray-500">Loading org units...</span>
            </div>
          ) : filteredOrgUnits.length === 0 ? (
            <Card>
              <CardContent className="p-8 text-center text-gray-500">
                No org units found
              </CardContent>
            </Card>
          ) : (
            <div className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {filteredOrgUnits.map((unit) => (
                  <Card
                    key={unit.orgUnitId}
                    className="hover:shadow-md transition-shadow cursor-pointer"
                    onClick={() => setSelectedOrgUnit(unit)}
                  >
                    <CardContent className="p-4">
                      <div className="flex items-start gap-3">
                        <div className="p-2 bg-green-100 rounded-lg">
                          <Building2 className="h-5 w-5 text-green-600" />
                        </div>
                        <div>
                          <h3 className="font-semibold text-gray-900">{unit.name}</h3>
                          <div className="flex items-center gap-2 mt-1">
                            <Badge variant="outline" className="text-xs">
                              {unit.code}
                            </Badge>
                            {unit.kind && (
                              <Badge variant="secondary" className="text-xs">
                                {unit.kind}
                              </Badge>
                            )}
                          </div>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
              {/* Pagination */}
              {orgUnitsTotalPages > 1 && (
                <div className="flex items-center justify-between border-t pt-4">
                  <p className="text-sm text-gray-500">
                    Showing {orgUnitsPage * PAGE_SIZE + 1} to {Math.min((orgUnitsPage + 1) * PAGE_SIZE, stats.orgUnits)} of {stats.orgUnits}
                  </p>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setOrgUnitsPage(p => Math.max(0, p - 1))}
                      disabled={orgUnitsPage === 0}
                    >
                      <ChevronLeft className="h-4 w-4" />
                      Previous
                    </Button>
                    <span className="text-sm text-gray-600 min-w-[100px] text-center">
                      Page {orgUnitsPage + 1} of {orgUnitsTotalPages}
                    </span>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setOrgUnitsPage(p => Math.min(orgUnitsTotalPages - 1, p + 1))}
                      disabled={orgUnitsPage >= orgUnitsTotalPages - 1}
                    >
                      Next
                      <ChevronRight className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              )}
            </div>
          )}
        </TabsContent>

        {/* Org Tree Tab */}
        <TabsContent value="tree" className="mt-4">
          {orgUnitsLoading ? (
            <div className="flex items-center justify-center p-12">
              <Loader2 className="h-8 w-8 animate-spin text-green-600" />
              <span className="ml-2 text-gray-500">Loading org tree...</span>
            </div>
          ) : (
            <Card>
              <CardContent className="p-4">
                <OrgTree
                  orgUnits={orgUnitsData?.items || []}
                  onSelect={setSelectedOrgUnit}
                  selectedId={selectedOrgUnit?.orgUnitId}
                />
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>

      {/* Modals */}
      <CreatePersonModal
        open={showPersonModal}
        onClose={() => setShowPersonModal(false)}
        onSubmit={handleCreatePerson}
        isLoading={createPerson.isPending}
        orgUnits={orgUnitsList}
      />
      <CreateTeamModal
        open={showTeamModal}
        onClose={() => setShowTeamModal(false)}
        onSubmit={handleCreateTeam}
        isLoading={createTeam.isPending}
        orgUnits={orgUnitsList}
      />
      <CreateOrgUnitModal
        open={showOrgUnitModal}
        onClose={() => setShowOrgUnitModal(false)}
        onSubmit={handleCreateOrgUnit}
        isLoading={createOrgUnit.isPending}
        orgUnits={orgUnitsList}
      />

      {/* Detail Modals */}
      <PersonDetailModal
        open={!!selectedPerson}
        onClose={() => setSelectedPerson(null)}
        person={selectedPerson}
        orgUnits={orgUnitsList}
        onUpdate={handleUpdatePerson}
        onDelete={handleDeletePerson}
        isUpdating={updatePerson.isPending}
        isDeleting={deletePerson.isPending}
      />
      <TeamDetailModal
        open={!!selectedTeam}
        onClose={() => setSelectedTeam(null)}
        team={selectedTeam}
        orgUnits={orgUnitsList}
        onUpdate={handleUpdateTeam}
        onDelete={handleDeleteTeam}
        isUpdating={updateTeam.isPending}
        isDeleting={deleteTeam.isPending}
      />
      <OrgUnitDetailModal
        open={!!selectedOrgUnit}
        onClose={() => setSelectedOrgUnit(null)}
        orgUnit={selectedOrgUnit}
        orgUnits={orgUnitsList}
        onUpdate={handleUpdateOrgUnit}
        onDelete={handleDeleteOrgUnit}
        isUpdating={updateOrgUnit.isPending}
        isDeleting={deleteOrgUnit.isPending}
      />
    </div>
  );
}
