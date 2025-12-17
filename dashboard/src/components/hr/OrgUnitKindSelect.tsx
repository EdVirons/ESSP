import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { ORG_UNIT_KINDS } from '@/lib/hr-constants';

interface OrgUnitKindSelectProps {
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
  id?: string;
}

export function OrgUnitKindSelect({
  value,
  onChange,
  disabled = false,
  id,
}: OrgUnitKindSelectProps) {
  return (
    <Select value={value} onValueChange={onChange} disabled={disabled}>
      <SelectTrigger id={id}>
        <SelectValue placeholder="Select type" />
      </SelectTrigger>
      <SelectContent>
        {ORG_UNIT_KINDS.map((kind) => (
          <SelectItem key={kind.value} value={kind.value}>
            {kind.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
