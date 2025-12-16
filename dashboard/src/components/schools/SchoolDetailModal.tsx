import {
  School,
  Sparkles,
  Laptop,
  Wifi,
  Target,
  AlertCircle,
  CheckCircle2,
  Loader2,
  X,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import type { SchoolSnapshot } from '@/api/ssot';
import { formatDate } from '@/lib/utils';
import { useEdTechProfile } from '@/hooks/useEdTechProfile';

interface SchoolDetailModalProps {
  school: SchoolSnapshot;
  onClose: () => void;
}

export function SchoolDetailModal({ school, onClose }: SchoolDetailModalProps) {
  const { data: edtechData, isLoading: edtechLoading } = useEdTechProfile(school.schoolId);
  const edtechProfile = edtechData?.profile;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/50" onClick={onClose} />

      {/* Modal */}
      <div className="relative bg-white rounded-lg shadow-lg w-full max-w-lg mx-4 max-h-[90vh] overflow-y-auto">
        <div className="p-6">
          {/* Header */}
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-full bg-blue-50">
                <School className="h-5 w-5 text-blue-600" />
              </div>
              <div>
                <h2 className="text-lg font-semibold text-gray-900">{school.name}</h2>
                <p className="text-sm text-gray-500">School Details</p>
              </div>
            </div>
            <Button variant="ghost" size="icon" onClick={onClose}>
              <X className="h-4 w-4" />
            </Button>
          </div>

          {/* Content */}
          <div className="space-y-4">
            {/* Basic Info */}
            <div className="grid grid-cols-2 gap-4">
              <InfoItem label="School ID" value={school.schoolId} />
              <InfoItem label="KNEC Code" value={school.knecCode} />
              <InfoItem label="UIC" value={school.uic} />
              <InfoItem label="Level" value={school.level} />
              <InfoItem label="Type" value={school.type} capitalize />
              <InfoItem label="Sex" value={school.sex} capitalize />
              <InfoItem label="Accommodation" value={school.accommodation} capitalize />
              <InfoItem label="Cluster" value={school.cluster} />
            </div>

            {/* Location */}
            <div className="border-t pt-4">
              <h3 className="text-sm font-medium text-gray-900 mb-2">Location</h3>
              <div className="grid grid-cols-2 gap-4">
                <InfoItem label="County" value={school.countyName} />
                <InfoItem label="County Code" value={school.countyCode} />
                <InfoItem label="Sub-County" value={school.subCountyName} />
                <InfoItem label="Sub-County Code" value={school.subCountyCode} />
              </div>
              {(school.latitude || school.longitude) && (
                <div className="mt-2 text-sm text-gray-500">
                  Coordinates: {school.latitude?.toFixed(6)}, {school.longitude?.toFixed(6)}
                </div>
              )}
            </div>

            {/* Metadata */}
            <div className="border-t pt-4">
              <h3 className="text-sm font-medium text-gray-900 mb-2">Metadata</h3>
              <InfoItem label="Last Updated" value={formatDate(school.updatedAt)} />
            </div>

            {/* EdTech Assessment */}
            <div className="border-t pt-4">
              <div className="flex items-center gap-2 mb-3">
                <div className="flex h-6 w-6 items-center justify-center rounded-md bg-gradient-to-br from-indigo-500 to-purple-600 text-white">
                  <Sparkles className="h-3 w-3" />
                </div>
                <h3 className="text-sm font-medium text-gray-900">EdTech Assessment</h3>
              </div>

              {edtechLoading ? (
                <div className="flex items-center justify-center py-4">
                  <Loader2 className="h-5 w-5 text-indigo-600 animate-spin" />
                  <span className="ml-2 text-sm text-gray-500">Loading assessment...</span>
                </div>
              ) : edtechProfile ? (
                <div className="space-y-3">
                  {/* Status Badge */}
                  <div className="flex items-center gap-2">
                    {edtechProfile.status === 'completed' ? (
                      <Badge className="bg-green-100 text-green-800 border-green-200">
                        <CheckCircle2 className="h-3 w-3 mr-1" />
                        Complete
                      </Badge>
                    ) : (
                      <Badge className="bg-amber-100 text-amber-800 border-amber-200">
                        Draft
                      </Badge>
                    )}
                    {edtechProfile.completedAt && (
                      <span className="text-xs text-gray-500">
                        Completed {formatDate(edtechProfile.completedAt)}
                      </span>
                    )}
                  </div>

                  {/* Quick Stats */}
                  <div className="grid grid-cols-3 gap-2">
                    <div className="bg-indigo-50 rounded-lg p-2 text-center">
                      <div className="flex items-center justify-center gap-1">
                        <Laptop className="h-3 w-3 text-indigo-600" />
                        <span className="text-lg font-bold text-indigo-600">
                          {edtechProfile.totalDevices || 0}
                        </span>
                      </div>
                      <div className="text-[10px] text-gray-500">Devices</div>
                    </div>
                    <div className="bg-teal-50 rounded-lg p-2 text-center">
                      <div className="flex items-center justify-center gap-1">
                        <Wifi className="h-3 w-3 text-teal-600" />
                        <span className="text-xs font-semibold text-teal-600 truncate">
                          {edtechProfile.networkQuality || '-'}
                        </span>
                      </div>
                      <div className="text-[10px] text-gray-500">Network</div>
                    </div>
                    <div className="bg-amber-50 rounded-lg p-2 text-center">
                      <div className="flex items-center justify-center gap-1">
                        <AlertCircle className="h-3 w-3 text-amber-600" />
                        <span className="text-lg font-bold text-amber-600">
                          {edtechProfile.biggestChallenges?.length || 0}
                        </span>
                      </div>
                      <div className="text-[10px] text-gray-500">Challenges</div>
                    </div>
                  </div>

                  {/* Device Breakdown */}
                  {edtechProfile.deviceTypes && (
                    <div className="bg-gray-50 rounded-lg p-2">
                      <div className="text-xs font-medium text-gray-700 mb-1">Device Breakdown</div>
                      <div className="flex flex-wrap gap-1.5 text-xs">
                        {edtechProfile.deviceTypes.laptops > 0 && (
                          <span className="bg-white px-1.5 py-0.5 rounded border text-gray-600">
                            Laptops: {edtechProfile.deviceTypes.laptops}
                          </span>
                        )}
                        {edtechProfile.deviceTypes.chromebooks > 0 && (
                          <span className="bg-white px-1.5 py-0.5 rounded border text-gray-600">
                            Chromebooks: {edtechProfile.deviceTypes.chromebooks}
                          </span>
                        )}
                        {edtechProfile.deviceTypes.tablets > 0 && (
                          <span className="bg-white px-1.5 py-0.5 rounded border text-gray-600">
                            Tablets: {edtechProfile.deviceTypes.tablets}
                          </span>
                        )}
                        {edtechProfile.deviceTypes.desktops > 0 && (
                          <span className="bg-white px-1.5 py-0.5 rounded border text-gray-600">
                            Desktops: {edtechProfile.deviceTypes.desktops}
                          </span>
                        )}
                      </div>
                    </div>
                  )}

                  {/* Top Priority */}
                  {edtechProfile.priorityRanking?.[0] && (
                    <div className="flex items-center gap-2 text-xs bg-emerald-50 rounded-lg p-2">
                      <Target className="h-3 w-3 text-emerald-600" />
                      <span className="text-gray-600">Top Priority:</span>
                      <span className="font-medium text-emerald-700">{edtechProfile.priorityRanking[0]}</span>
                    </div>
                  )}

                  {/* AI Summary */}
                  {edtechProfile.aiSummary && (
                    <div className="bg-purple-50 rounded-lg p-2">
                      <div className="flex items-center gap-1 text-xs font-medium text-purple-700 mb-1">
                        <Sparkles className="h-3 w-3" />
                        AI Summary
                      </div>
                      <p className="text-xs text-gray-700 line-clamp-3">{edtechProfile.aiSummary}</p>
                    </div>
                  )}

                  {/* AI Recommendations Count */}
                  {edtechProfile.aiRecommendations && edtechProfile.aiRecommendations.length > 0 && (
                    <div className="text-xs text-gray-500">
                      {edtechProfile.aiRecommendations.length} AI recommendation{edtechProfile.aiRecommendations.length !== 1 ? 's' : ''} available
                    </div>
                  )}
                </div>
              ) : (
                <div className="bg-gray-50 rounded-lg p-4 text-center">
                  <Sparkles className="h-6 w-6 text-gray-300 mx-auto mb-2" />
                  <p className="text-sm text-gray-500">No EdTech assessment completed</p>
                  <p className="text-xs text-gray-400 mt-1">
                    The school contact can complete this from their dashboard
                  </p>
                </div>
              )}
            </div>
          </div>

          {/* Footer */}
          <div className="mt-6 flex justify-end">
            <Button variant="outline" onClick={onClose}>
              Close
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}

// Helper component for displaying info items
interface InfoItemProps {
  label: string;
  value?: string | null;
  capitalize?: boolean;
}

function InfoItem({ label, value, capitalize }: InfoItemProps) {
  return (
    <div>
      <dt className="text-xs text-gray-500">{label}</dt>
      <dd className={`text-sm font-medium text-gray-900 ${capitalize ? 'capitalize' : ''}`}>
        {value || '-'}
      </dd>
    </div>
  );
}
