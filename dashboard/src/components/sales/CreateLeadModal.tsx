import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import type { CreateDemoLeadRequest, DemoLeadSource } from '@/types/sales';
import { sourceLabels } from '@/types/sales';
import { useCounties, useSubCounties } from '@/api/ssot';

interface CreateLeadModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateDemoLeadRequest) => void;
  isLoading: boolean;
}

const leadSources: DemoLeadSource[] = ['website', 'referral', 'event', 'cold_outreach', 'inbound'];

export function CreateLeadModal({
  open,
  onClose,
  onSubmit,
  isLoading,
}: CreateLeadModalProps) {
  const [formData, setFormData] = React.useState<CreateDemoLeadRequest>({
    schoolName: '',
    contactName: '',
    contactEmail: '',
    contactPhone: '',
    contactRole: '',
    countyCode: '',
    countyName: '',
    subCountyCode: '',
    subCountyName: '',
    estimatedValue: undefined,
    estimatedDevices: undefined,
    leadSource: 'website',
    notes: '',
    tags: [],
  });

  // Fetch counties and sub-counties
  const { data: countiesData } = useCounties();
  const { data: subCountiesData } = useSubCounties(formData.countyCode);

  const counties = React.useMemo(() => {
    if (!countiesData?.items) return [];
    return countiesData.items;
  }, [countiesData?.items]);

  const subCounties = React.useMemo(() => {
    const items = subCountiesData?.items;
    if (!items) return [];
    return items;
  }, [subCountiesData?.items]);

  const handleCountyChange = (code: string) => {
    const county = counties?.find(c => c.code === code);
    setFormData({
      ...formData,
      countyCode: code,
      countyName: county?.name || '',
      subCountyCode: '',
      subCountyName: '',
    });
  };

  const handleSubCountyChange = (code: string) => {
    const subCounty = subCounties?.find(s => s.code === code);
    setFormData({
      ...formData,
      subCountyCode: code,
      subCountyName: subCounty?.name || '',
    });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  const handleClose = () => {
    setFormData({
      schoolName: '',
      contactName: '',
      contactEmail: '',
      contactPhone: '',
      contactRole: '',
      countyCode: '',
      countyName: '',
      subCountyCode: '',
      subCountyName: '',
      estimatedValue: undefined,
      estimatedDevices: undefined,
      leadSource: 'website',
      notes: '',
      tags: [],
    });
    onClose();
  };

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>Add New Lead</ModalHeader>
      <form onSubmit={handleSubmit}>
        <ModalBody>
          <div className="space-y-4">
            {/* School Info */}
            <div>
              <Label htmlFor="schoolName">School Name *</Label>
              <Input
                id="schoolName"
                value={formData.schoolName}
                onChange={(e) => setFormData({ ...formData, schoolName: e.target.value })}
                placeholder="Enter school name"
                required
              />
            </div>

            {/* Contact Info */}
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="contactName">Contact Name</Label>
                <Input
                  id="contactName"
                  value={formData.contactName}
                  onChange={(e) => setFormData({ ...formData, contactName: e.target.value })}
                  placeholder="Contact person name"
                />
              </div>
              <div>
                <Label htmlFor="contactRole">Role</Label>
                <Input
                  id="contactRole"
                  value={formData.contactRole}
                  onChange={(e) => setFormData({ ...formData, contactRole: e.target.value })}
                  placeholder="e.g., Principal, IT Admin"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="contactEmail">Email</Label>
                <Input
                  id="contactEmail"
                  type="email"
                  value={formData.contactEmail}
                  onChange={(e) => setFormData({ ...formData, contactEmail: e.target.value })}
                  placeholder="email@school.edu"
                />
              </div>
              <div>
                <Label htmlFor="contactPhone">Phone</Label>
                <Input
                  id="contactPhone"
                  value={formData.contactPhone}
                  onChange={(e) => setFormData({ ...formData, contactPhone: e.target.value })}
                  placeholder="+254 7XX XXX XXX"
                />
              </div>
            </div>

            {/* Location Info */}
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="county">County</Label>
                <Select
                  value={formData.countyCode || ''}
                  onValueChange={handleCountyChange}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select county" />
                  </SelectTrigger>
                  <SelectContent>
                    {counties.map((county) => (
                      <SelectItem key={county.code} value={county.code}>
                        {county.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label htmlFor="subCounty">Sub-County</Label>
                <Select
                  value={formData.subCountyCode || ''}
                  onValueChange={handleSubCountyChange}
                  disabled={!formData.countyCode}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select sub-county" />
                  </SelectTrigger>
                  <SelectContent>
                    {subCounties.map((subCounty) => (
                      <SelectItem key={subCounty.code} value={subCounty.code}>
                        {subCounty.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Deal Info */}
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="estimatedValue">Estimated Value (KES)</Label>
                <Input
                  id="estimatedValue"
                  type="number"
                  min="0"
                  value={formData.estimatedValue || ''}
                  onChange={(e) => setFormData({
                    ...formData,
                    estimatedValue: e.target.value ? Number(e.target.value) : undefined
                  })}
                  placeholder="0"
                />
              </div>
              <div>
                <Label htmlFor="estimatedDevices">Estimated Devices</Label>
                <Input
                  id="estimatedDevices"
                  type="number"
                  min="0"
                  value={formData.estimatedDevices || ''}
                  onChange={(e) => setFormData({
                    ...formData,
                    estimatedDevices: e.target.value ? Number(e.target.value) : undefined
                  })}
                  placeholder="0"
                />
              </div>
            </div>

            {/* Lead Source */}
            <div>
              <Label htmlFor="leadSource">Lead Source</Label>
              <Select
                value={formData.leadSource}
                onValueChange={(value) =>
                  setFormData({ ...formData, leadSource: value as DemoLeadSource })
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select source" />
                </SelectTrigger>
                <SelectContent>
                  {leadSources.map((source) => (
                    <SelectItem key={source} value={source}>
                      {sourceLabels[source]}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Notes */}
            <div>
              <Label htmlFor="notes">Notes</Label>
              <Textarea
                id="notes"
                value={formData.notes}
                onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                placeholder="Initial inquiry details, interests, etc."
                rows={3}
              />
            </div>
          </div>
        </ModalBody>
        <ModalFooter>
          <Button type="button" variant="outline" onClick={handleClose}>
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={!formData.schoolName.trim() || isLoading}
          >
            {isLoading ? 'Creating...' : 'Create Lead'}
          </Button>
        </ModalFooter>
      </form>
    </Modal>
  );
}
