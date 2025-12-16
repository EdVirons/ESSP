import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Calendar, Clock, MapPin, Video, UserPlus, X, Loader2 } from 'lucide-react';
import type { CreateDemoScheduleRequest, DemoAttendee } from '@/types/sales';

interface ScheduleDemoModalProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateDemoScheduleRequest) => void;
  isLoading: boolean;
  leadName?: string;
}

const durationOptions = [
  { value: '30', label: '30 minutes' },
  { value: '45', label: '45 minutes' },
  { value: '60', label: '1 hour' },
  { value: '90', label: '1.5 hours' },
  { value: '120', label: '2 hours' },
];

export function ScheduleDemoModal({
  open,
  onClose,
  onSubmit,
  isLoading,
  leadName,
}: ScheduleDemoModalProps) {
  const [formData, setFormData] = React.useState<CreateDemoScheduleRequest>({
    scheduledDate: '',
    scheduledTime: '10:00',
    durationMinutes: 60,
    location: '',
    meetingLink: '',
    attendees: [],
  });

  const [newAttendee, setNewAttendee] = React.useState<DemoAttendee>({
    name: '',
    email: '',
    role: '',
  });

  // Get tomorrow's date as minimum date
  const minDate = React.useMemo(() => {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    return tomorrow.toISOString().split('T')[0];
  }, []);

  const handleAddAttendee = () => {
    if (newAttendee.name && newAttendee.email) {
      setFormData({
        ...formData,
        attendees: [...(formData.attendees || []), newAttendee],
      });
      setNewAttendee({ name: '', email: '', role: '' });
    }
  };

  const handleRemoveAttendee = (index: number) => {
    setFormData({
      ...formData,
      attendees: formData.attendees?.filter((_, i) => i !== index),
    });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  const handleClose = () => {
    setFormData({
      scheduledDate: '',
      scheduledTime: '10:00',
      durationMinutes: 60,
      location: '',
      meetingLink: '',
      attendees: [],
    });
    setNewAttendee({ name: '', email: '', role: '' });
    onClose();
  };

  const isValid = formData.scheduledDate && formData.scheduledTime;

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>
        <div className="flex items-center gap-2">
          <Calendar className="h-5 w-5 text-purple-600" />
          Schedule Demo
        </div>
        {leadName && (
          <p className="text-sm font-normal text-gray-500 mt-1">
            for {leadName}
          </p>
        )}
      </ModalHeader>
      <form onSubmit={handleSubmit}>
        <ModalBody>
          <div className="space-y-4">
            {/* Date and Time */}
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label htmlFor="scheduledDate" className="flex items-center gap-2">
                  <Calendar className="h-4 w-4 text-gray-500" />
                  Date *
                </Label>
                <Input
                  id="scheduledDate"
                  type="date"
                  min={minDate}
                  value={formData.scheduledDate}
                  onChange={(e) => setFormData({ ...formData, scheduledDate: e.target.value })}
                  required
                />
              </div>
              <div>
                <Label htmlFor="scheduledTime" className="flex items-center gap-2">
                  <Clock className="h-4 w-4 text-gray-500" />
                  Time *
                </Label>
                <Input
                  id="scheduledTime"
                  type="time"
                  value={formData.scheduledTime}
                  onChange={(e) => setFormData({ ...formData, scheduledTime: e.target.value })}
                  required
                />
              </div>
            </div>

            {/* Duration */}
            <div>
              <Label htmlFor="duration">Duration</Label>
              <Select
                value={String(formData.durationMinutes)}
                onValueChange={(value) =>
                  setFormData({ ...formData, durationMinutes: Number(value) })
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select duration" />
                </SelectTrigger>
                <SelectContent>
                  {durationOptions.map((option) => (
                    <SelectItem key={option.value} value={option.value}>
                      {option.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Location */}
            <div>
              <Label htmlFor="location" className="flex items-center gap-2">
                <MapPin className="h-4 w-4 text-gray-500" />
                Location
              </Label>
              <Input
                id="location"
                value={formData.location}
                onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                placeholder="e.g., School premises, Conference room"
              />
            </div>

            {/* Meeting Link */}
            <div>
              <Label htmlFor="meetingLink" className="flex items-center gap-2">
                <Video className="h-4 w-4 text-gray-500" />
                Meeting Link (for virtual demos)
              </Label>
              <Input
                id="meetingLink"
                type="url"
                value={formData.meetingLink}
                onChange={(e) => setFormData({ ...formData, meetingLink: e.target.value })}
                placeholder="https://meet.google.com/..."
              />
            </div>

            {/* Attendees Section */}
            <div className="space-y-3">
              <Label className="flex items-center gap-2">
                <UserPlus className="h-4 w-4 text-gray-500" />
                Attendees
              </Label>

              {/* Current Attendees */}
              {formData.attendees && formData.attendees.length > 0 && (
                <div className="space-y-2">
                  {formData.attendees.map((attendee, index) => (
                    <div
                      key={index}
                      className="flex items-center justify-between bg-gray-50 rounded-lg px-3 py-2 text-sm"
                    >
                      <div>
                        <span className="font-medium">{attendee.name}</span>
                        <span className="text-gray-500 mx-2">•</span>
                        <span className="text-gray-600">{attendee.email}</span>
                        {attendee.role && (
                          <>
                            <span className="text-gray-500 mx-2">•</span>
                            <span className="text-gray-500">{attendee.role}</span>
                          </>
                        )}
                      </div>
                      <button
                        type="button"
                        onClick={() => handleRemoveAttendee(index)}
                        className="text-gray-400 hover:text-red-500"
                        aria-label={`Remove ${attendee.name}`}
                      >
                        <X className="h-4 w-4" />
                      </button>
                    </div>
                  ))}
                </div>
              )}

              {/* Add New Attendee */}
              <div className="bg-gray-50 rounded-lg p-3 space-y-3">
                <div className="grid grid-cols-2 gap-2">
                  <Input
                    placeholder="Name"
                    value={newAttendee.name}
                    onChange={(e) => setNewAttendee({ ...newAttendee, name: e.target.value })}
                  />
                  <Input
                    placeholder="Email"
                    type="email"
                    value={newAttendee.email}
                    onChange={(e) => setNewAttendee({ ...newAttendee, email: e.target.value })}
                  />
                </div>
                <div className="flex gap-2">
                  <Input
                    placeholder="Role (optional)"
                    value={newAttendee.role}
                    onChange={(e) => setNewAttendee({ ...newAttendee, role: e.target.value })}
                    className="flex-1"
                  />
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={handleAddAttendee}
                    disabled={!newAttendee.name || !newAttendee.email}
                  >
                    Add
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </ModalBody>
        <ModalFooter>
          <Button type="button" variant="outline" onClick={handleClose}>
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={!isValid || isLoading}
            className="bg-purple-600 hover:bg-purple-700"
          >
            {isLoading ? (
              <>
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                Scheduling...
              </>
            ) : (
              'Schedule Demo'
            )}
          </Button>
        </ModalFooter>
      </form>
    </Modal>
  );
}
