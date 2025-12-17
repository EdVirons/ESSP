import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import type { OrgUnitSnapshot } from '@/types/hr';

interface OrgUnitSelectProps {
  value: string;
  onChange: (value: string) => void;
  orgUnits: OrgUnitSnapshot[];
  excludeId?: string; // Exclude this ID to prevent circular references
  placeholder?: string;
  allowNone?: boolean;
  noneLabel?: string;
  disabled?: boolean;
  id?: string;
}

export function OrgUnitSelect({
  value,
  onChange,
  orgUnits,
  excludeId,
  placeholder = 'Select org unit',
  allowNone = true,
  noneLabel = 'None',
  disabled = false,
  id,
}: OrgUnitSelectProps) {
  const filteredUnits = excludeId
    ? orgUnits.filter((u) => u.orgUnitId !== excludeId)
    : orgUnits;

  return (
    <Select value={value} onValueChange={onChange} disabled={disabled}>
      <SelectTrigger id={id}>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {allowNone && <SelectItem value="">{noneLabel}</SelectItem>}
        {filteredUnits.map((unit) => (
          <SelectItem key={unit.orgUnitId} value={unit.orgUnitId}>
            {unit.name} {unit.code && `(${unit.code})`}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
