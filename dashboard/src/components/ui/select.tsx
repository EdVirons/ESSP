import * as React from 'react';
import { ChevronDown, Check, X } from 'lucide-react';
import { cn } from '@/lib/utils';

// =============================================================================
// Select Option type
// =============================================================================

interface SelectOption {
  value: string;
  label: string;
  disabled?: boolean;
}

// =============================================================================
// Compound Select Context (for Radix-style usage)
// =============================================================================

interface SelectContextValue {
  value?: string;
  onValueChange: (value: string) => void;
  open: boolean;
  setOpen: (open: boolean) => void;
}

const SelectContext = React.createContext<SelectContextValue | undefined>(undefined);

function useSelectContext() {
  const context = React.useContext(SelectContext);
  if (!context) {
    throw new Error('Select components must be used within a Select');
  }
  return context;
}

// =============================================================================
// Select Component (supports both legacy and compound patterns)
// =============================================================================

export interface SelectProps {
  value?: string;
  defaultValue?: string;
  // Legacy props (options pattern)
  onChange?: (value: string) => void;
  options?: SelectOption[];
  placeholder?: string;
  disabled?: boolean;
  className?: string;
  error?: boolean;
  // Compound props
  onValueChange?: (value: string) => void;
  children?: React.ReactNode;
}

export function Select({
  value,
  defaultValue,
  onChange,
  onValueChange,
  options,
  placeholder = 'Select...',
  disabled = false,
  className,
  error = false,
  children,
}: SelectProps) {
  const [internalValue, setInternalValue] = React.useState(defaultValue);
  const [open, setOpen] = React.useState(false);
  const containerRef = React.useRef<HTMLDivElement>(null);

  const isControlled = value !== undefined;
  const currentValue = isControlled ? value : internalValue;

  // Combined change handler for both legacy and compound patterns
  const handleValueChange = (newValue: string) => {
    if (!isControlled) {
      setInternalValue(newValue);
    }
    onChange?.(newValue);
    onValueChange?.(newValue);
  };

  // Close on click outside (for legacy pattern)
  React.useEffect(() => {
    if (!options) return; // Only for legacy pattern
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [options]);

  // If options are provided, use legacy rendering
  if (options) {
    const selectedOption = options.find((opt) => opt.value === currentValue);

    return (
      <div ref={containerRef} className={cn('relative', className)}>
        <button
          type="button"
          onClick={() => !disabled && setOpen(!open)}
          className={cn(
            'flex h-10 w-full items-center justify-between rounded-md border bg-white px-3 py-2 text-sm',
            'focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2',
            error ? 'border-red-300' : 'border-gray-300',
            disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'
          )}
          disabled={disabled}
        >
          <span className={selectedOption ? 'text-gray-900' : 'text-gray-500'}>
            {selectedOption?.label || placeholder}
          </span>
          <ChevronDown className="h-4 w-4 text-gray-400" />
        </button>

        {open && (
          <div className="absolute z-50 mt-1 w-full rounded-md border border-gray-200 bg-white shadow-lg">
            <div className="max-h-60 overflow-auto py-1">
              {options.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  onClick={() => {
                    if (!option.disabled) {
                      handleValueChange(option.value);
                      setOpen(false);
                    }
                  }}
                  className={cn(
                    'flex w-full items-center justify-between px-3 py-2 text-sm',
                    option.disabled
                      ? 'cursor-not-allowed text-gray-400'
                      : 'cursor-pointer text-gray-900 hover:bg-gray-50',
                    option.value === currentValue && 'bg-blue-50'
                  )}
                  disabled={option.disabled}
                >
                  <span>{option.label}</span>
                  {option.value === currentValue && <Check className="h-4 w-4 text-blue-600" />}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>
    );
  }

  // Otherwise, use compound component pattern
  return (
    <SelectContext.Provider value={{ value: currentValue, onValueChange: handleValueChange, open, setOpen }}>
      <div className={cn('relative', className)}>{children}</div>
    </SelectContext.Provider>
  );
}

// =============================================================================
// Compound Components (for Radix-style usage)
// =============================================================================

export interface SelectTriggerProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children?: React.ReactNode;
}

export const SelectTrigger = React.forwardRef<HTMLButtonElement, SelectTriggerProps>(
  ({ className, children, ...props }, ref) => {
    const { open, setOpen } = useSelectContext();

    return (
      <button
        ref={ref}
        type="button"
        role="combobox"
        aria-expanded={open}
        onClick={() => setOpen(!open)}
        className={cn(
          'flex h-10 w-full items-center justify-between rounded-md border border-gray-300 bg-white px-3 py-2 text-sm',
          'focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2',
          'disabled:cursor-not-allowed disabled:opacity-50',
          className
        )}
        {...props}
      >
        {children}
        <ChevronDown className="h-4 w-4 opacity-50" />
      </button>
    );
  }
);
SelectTrigger.displayName = 'SelectTrigger';

export function SelectValue({ placeholder }: { placeholder?: string }) {
  const { value } = useSelectContext();
  return <span className={value ? 'text-gray-900' : 'text-gray-500'}>{value || placeholder}</span>;
}

export interface SelectContentProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
}

export const SelectContent = React.forwardRef<HTMLDivElement, SelectContentProps>(
  ({ className, children, ...props }, _ref) => {
    const { open, setOpen } = useSelectContext();
    const contentRef = React.useRef<HTMLDivElement>(null);

    React.useEffect(() => {
      const handleClickOutside = (e: MouseEvent) => {
        if (contentRef.current && !contentRef.current.contains(e.target as Node)) {
          setOpen(false);
        }
      };

      if (open) {
        document.addEventListener('mousedown', handleClickOutside);
      }

      return () => {
        document.removeEventListener('mousedown', handleClickOutside);
      };
    }, [open, setOpen]);

    if (!open) return null;

    return (
      <div
        ref={contentRef}
        className={cn(
          'absolute z-50 mt-1 w-full min-w-[8rem] overflow-hidden rounded-md border border-gray-200 bg-white shadow-md',
          'animate-in fade-in-0 zoom-in-95',
          className
        )}
        {...props}
      >
        <div className="max-h-60 overflow-auto py-1">{children}</div>
      </div>
    );
  }
);
SelectContent.displayName = 'SelectContent';

export interface SelectItemProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  value: string;
  children: React.ReactNode;
}

export const SelectItem = React.forwardRef<HTMLButtonElement, SelectItemProps>(
  ({ className, value, children, ...props }, ref) => {
    const { value: selectedValue, onValueChange, setOpen } = useSelectContext();
    const isSelected = selectedValue === value;

    return (
      <button
        ref={ref}
        type="button"
        role="option"
        aria-selected={isSelected}
        onClick={() => {
          onValueChange(value);
          setOpen(false);
        }}
        className={cn(
          'relative flex w-full cursor-pointer select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none',
          'hover:bg-gray-100 focus:bg-gray-100',
          isSelected && 'bg-gray-100',
          className
        )}
        {...props}
      >
        <span className="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
          {isSelected && <Check className="h-4 w-4" />}
        </span>
        {children}
      </button>
    );
  }
);
SelectItem.displayName = 'SelectItem';

export function SelectGroup({ children }: { children: React.ReactNode }) {
  return <div className="p-1">{children}</div>;
}

export function SelectLabel({ children, className }: { children: React.ReactNode; className?: string }) {
  return <div className={cn('py-1.5 pl-8 pr-2 text-sm font-semibold', className)}>{children}</div>;
}

export function SelectSeparator({ className }: { className?: string }) {
  return <div className={cn('-mx-1 my-1 h-px bg-gray-100', className)} />;
}

// =============================================================================
// MultiSelect Component
// =============================================================================

interface MultiSelectProps {
  value: string[];
  onChange: (value: string[]) => void;
  options: SelectOption[];
  placeholder?: string;
  disabled?: boolean;
  className?: string;
  maxDisplay?: number;
}

export function MultiSelect({
  value = [],
  onChange,
  options,
  placeholder = 'Select...',
  disabled = false,
  className,
  maxDisplay = 3,
}: MultiSelectProps) {
  const [open, setOpen] = React.useState(false);
  const containerRef = React.useRef<HTMLDivElement>(null);

  React.useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const selectedLabels = value
    .map((v) => options.find((opt) => opt.value === v)?.label)
    .filter(Boolean);

  const displayText =
    selectedLabels.length === 0
      ? placeholder
      : selectedLabels.length <= maxDisplay
      ? selectedLabels.join(', ')
      : `${selectedLabels.slice(0, maxDisplay).join(', ')} +${selectedLabels.length - maxDisplay}`;

  const handleToggle = (optionValue: string) => {
    if (value.includes(optionValue)) {
      onChange(value.filter((v) => v !== optionValue));
    } else {
      onChange([...value, optionValue]);
    }
  };

  const handleRemove = (optionValue: string, e: React.MouseEvent) => {
    e.stopPropagation();
    onChange(value.filter((v) => v !== optionValue));
  };

  return (
    <div ref={containerRef} className={cn('relative', className)}>
      <button
        type="button"
        onClick={() => !disabled && setOpen(!open)}
        className={cn(
          'flex h-10 w-full items-center justify-between rounded-md border border-gray-300 bg-white px-3 py-2 text-sm',
          'focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2',
          disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'
        )}
        disabled={disabled}
      >
        <span className={value.length > 0 ? 'text-gray-900' : 'text-gray-500'}>
          {displayText}
        </span>
        <ChevronDown className="h-4 w-4 text-gray-400" />
      </button>

      {open && (
        <div className="absolute z-50 mt-1 w-full rounded-md border border-gray-200 bg-white shadow-lg">
          <div className="max-h-60 overflow-auto py-1">
            {options.map((option) => (
              <button
                key={option.value}
                type="button"
                onClick={() => !option.disabled && handleToggle(option.value)}
                className={cn(
                  'flex w-full items-center justify-between px-3 py-2 text-sm',
                  option.disabled
                    ? 'cursor-not-allowed text-gray-400'
                    : 'cursor-pointer text-gray-900 hover:bg-gray-50',
                  value.includes(option.value) && 'bg-blue-50'
                )}
                disabled={option.disabled}
              >
                <span>{option.label}</span>
                {value.includes(option.value) && <Check className="h-4 w-4 text-blue-600" />}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Selected tags */}
      {value.length > 0 && (
        <div className="flex flex-wrap gap-1 mt-2">
          {value.map((v) => {
            const option = options.find((opt) => opt.value === v);
            if (!option) return null;
            return (
              <span
                key={v}
                className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs bg-blue-100 text-blue-800"
              >
                {option.label}
                <button
                  type="button"
                  onClick={(e) => handleRemove(v, e)}
                  className="hover:bg-blue-200 rounded-full p-0.5"
                >
                  <X className="h-3 w-3" />
                </button>
              </span>
            );
          })}
        </div>
      )}
    </div>
  );
}
