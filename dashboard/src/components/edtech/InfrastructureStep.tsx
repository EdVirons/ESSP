import { Laptop, Smartphone, Monitor, Tablet, HardDrive } from 'lucide-react';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Checkbox } from '@/components/ui/checkbox';
import { QuizQuestion } from './QuizQuestion';
import type { InfrastructureStepData, EdTechFormOptions } from '@/types/edtech';

interface InfrastructureStepProps {
  data: InfrastructureStepData;
  options: EdTechFormOptions | undefined;
  onChange: (data: InfrastructureStepData) => void;
  currentQuestion: number;
  totalQuestions: number;
}

export function InfrastructureStep({
  data,
  options,
  onChange,
  currentQuestion,
  totalQuestions,
}: InfrastructureStepProps) {
  const updateField = <K extends keyof InfrastructureStepData>(
    field: K,
    value: InfrastructureStepData[K]
  ) => {
    onChange({ ...data, [field]: value });
  };

  const updateDeviceType = (type: keyof typeof data.deviceTypes, value: number) => {
    const newDeviceTypes = { ...data.deviceTypes, [type]: value };
    const totalDevices = Object.values(newDeviceTypes).reduce((sum, v) => sum + (v || 0), 0);
    onChange({
      ...data,
      deviceTypes: newDeviceTypes,
      totalDevices,
    });
  };

  const toggleSoftware = (software: string) => {
    const current = data.existingSoftware || [];
    if (current.includes(software)) {
      updateField('existingSoftware', current.filter((s) => s !== software));
    } else {
      updateField('existingSoftware', [...current, software]);
    }
  };

  const renderQuestion = () => {
    switch (currentQuestion) {
      case 1:
        return (
          <QuizQuestion
            questionNumber={1}
            totalQuestions={totalQuestions}
            question="How many digital devices do you use in school currently?"
            description="Enter the count for each device type your school owns or uses."
          >
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="flex items-center gap-3 p-3 bg-white rounded-lg border border-gray-200">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-100">
                  <Laptop className="h-5 w-5 text-blue-600" />
                </div>
                <div className="flex-1">
                  <Label htmlFor="laptops" className="text-sm font-medium">Laptops</Label>
                  <Input
                    id="laptops"
                    type="number"
                    min={0}
                    value={data.deviceTypes.laptops || ''}
                    onChange={(e) => updateDeviceType('laptops', parseInt(e.target.value) || 0)}
                    placeholder="0"
                    className="mt-1"
                  />
                </div>
              </div>

              <div className="flex items-center gap-3 p-3 bg-white rounded-lg border border-gray-200">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-green-100">
                  <Laptop className="h-5 w-5 text-green-600" />
                </div>
                <div className="flex-1">
                  <Label htmlFor="chromebooks" className="text-sm font-medium">Chromebooks</Label>
                  <Input
                    id="chromebooks"
                    type="number"
                    min={0}
                    value={data.deviceTypes.chromebooks || ''}
                    onChange={(e) => updateDeviceType('chromebooks', parseInt(e.target.value) || 0)}
                    placeholder="0"
                    className="mt-1"
                  />
                </div>
              </div>

              <div className="flex items-center gap-3 p-3 bg-white rounded-lg border border-gray-200">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-purple-100">
                  <Tablet className="h-5 w-5 text-purple-600" />
                </div>
                <div className="flex-1">
                  <Label htmlFor="tablets" className="text-sm font-medium">Tablets</Label>
                  <Input
                    id="tablets"
                    type="number"
                    min={0}
                    value={data.deviceTypes.tablets || ''}
                    onChange={(e) => updateDeviceType('tablets', parseInt(e.target.value) || 0)}
                    placeholder="0"
                    className="mt-1"
                  />
                </div>
              </div>

              <div className="flex items-center gap-3 p-3 bg-white rounded-lg border border-gray-200">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-amber-100">
                  <Monitor className="h-5 w-5 text-amber-600" />
                </div>
                <div className="flex-1">
                  <Label htmlFor="desktops" className="text-sm font-medium">Desktops</Label>
                  <Input
                    id="desktops"
                    type="number"
                    min={0}
                    value={data.deviceTypes.desktops || ''}
                    onChange={(e) => updateDeviceType('desktops', parseInt(e.target.value) || 0)}
                    placeholder="0"
                    className="mt-1"
                  />
                </div>
              </div>

              <div className="flex items-center gap-3 p-3 bg-white rounded-lg border border-gray-200">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-gray-100">
                  <HardDrive className="h-5 w-5 text-gray-600" />
                </div>
                <div className="flex-1">
                  <Label htmlFor="other" className="text-sm font-medium">Other Devices</Label>
                  <Input
                    id="other"
                    type="number"
                    min={0}
                    value={data.deviceTypes.other || ''}
                    onChange={(e) => updateDeviceType('other', parseInt(e.target.value) || 0)}
                    placeholder="0"
                    className="mt-1"
                  />
                </div>
              </div>

              <div className="flex items-center gap-3 p-3 bg-indigo-50 rounded-lg border border-indigo-200">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-indigo-100">
                  <Smartphone className="h-5 w-5 text-indigo-600" />
                </div>
                <div className="flex-1">
                  <Label className="text-sm font-medium text-indigo-700">Total Devices</Label>
                  <div className="mt-1 h-10 flex items-center text-2xl font-bold text-indigo-600">
                    {data.totalDevices || 0}
                  </div>
                </div>
              </div>
            </div>
          </QuizQuestion>
        );

      case 2:
        return (
          <QuizQuestion
            questionNumber={2}
            totalQuestions={totalQuestions}
            question="What is the average age of your devices?"
            description="Select the option that best describes the age of most devices in your school."
          >
            <Select value={data.deviceAge} onValueChange={(v) => updateField('deviceAge', v)}>
              <SelectTrigger className="w-full max-w-md">
                <SelectValue placeholder="Select device age range" />
              </SelectTrigger>
              <SelectContent>
                {(options?.deviceAge || []).map((age) => (
                  <SelectItem key={age} value={age}>
                    {age.replace(/_/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase())}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </QuizQuestion>
        );

      case 3:
        return (
          <QuizQuestion
            questionNumber={3}
            totalQuestions={totalQuestions}
            question="How would you rate your network quality?"
            description="Consider factors like reliability, coverage, and speed when rating."
          >
            <Select value={data.networkQuality} onValueChange={(v) => updateField('networkQuality', v)}>
              <SelectTrigger className="w-full max-w-md">
                <SelectValue placeholder="Select network quality" />
              </SelectTrigger>
              <SelectContent>
                {(options?.networkQuality || []).map((quality) => (
                  <SelectItem key={quality} value={quality}>
                    {quality.charAt(0).toUpperCase() + quality.slice(1)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </QuizQuestion>
        );

      case 4:
        return (
          <QuizQuestion
            questionNumber={4}
            totalQuestions={totalQuestions}
            question="What is your internet connection speed?"
            description="Select the speed tier that best matches your school's internet connection."
          >
            <Select value={data.internetSpeed} onValueChange={(v) => updateField('internetSpeed', v)}>
              <SelectTrigger className="w-full max-w-md">
                <SelectValue placeholder="Select internet speed" />
              </SelectTrigger>
              <SelectContent>
                {(options?.internetSpeed || []).map((speed) => (
                  <SelectItem key={speed} value={speed}>
                    {speed.charAt(0).toUpperCase() + speed.slice(1)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </QuizQuestion>
        );

      case 5:
        return (
          <QuizQuestion
            questionNumber={5}
            totalQuestions={totalQuestions}
            question="What Learning Management System does your school use?"
            description="Select your primary LMS platform, or 'None' if you don't use one."
          >
            <Select value={data.lmsPlatform} onValueChange={(v) => updateField('lmsPlatform', v)}>
              <SelectTrigger className="w-full max-w-md">
                <SelectValue placeholder="Select LMS platform" />
              </SelectTrigger>
              <SelectContent>
                {(options?.lmsPlatforms || []).map((lms) => (
                  <SelectItem key={lms} value={lms}>
                    {lms}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </QuizQuestion>
        );

      case 6:
        return (
          <QuizQuestion
            questionNumber={6}
            totalQuestions={totalQuestions}
            question="What educational software does your school currently use?"
            description="Select all the software tools and platforms your school uses."
          >
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {(options?.existingSoftware || []).map((software) => (
                <label
                  key={software}
                  className="flex items-center gap-3 p-3 bg-white rounded-lg border border-gray-200 cursor-pointer hover:border-indigo-300 hover:bg-indigo-50/50 transition-colors"
                >
                  <Checkbox
                    id={`software-${software}`}
                    checked={data.existingSoftware?.includes(software)}
                    onCheckedChange={() => toggleSoftware(software)}
                  />
                  <span className="text-sm font-medium">{software}</span>
                </label>
              ))}
            </div>
          </QuizQuestion>
        );

      case 7:
        return (
          <QuizQuestion
            questionNumber={7}
            totalQuestions={totalQuestions}
            question="How many IT staff support your school?"
            description="Include full-time IT staff, part-time support, and teachers with IT responsibilities."
          >
            <div className="max-w-md">
              <Input
                id="itStaffCount"
                type="number"
                min={0}
                value={data.itStaffCount || ''}
                onChange={(e) => updateField('itStaffCount', parseInt(e.target.value) || 0)}
                placeholder="Enter number of IT staff"
                className="text-lg"
              />
            </div>
          </QuizQuestion>
        );

      default:
        return null;
    }
  };

  return renderQuestion();
}
