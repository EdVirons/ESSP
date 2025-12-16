import * as React from 'react';
import { X } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';

type SheetSide = 'left' | 'right' | 'top' | 'bottom';

interface SheetProps {
  open: boolean;
  onClose: () => void;
  children: React.ReactNode;
  side?: SheetSide;
  className?: string;
}

const sideStyles: Record<SheetSide, string> = {
  left: 'inset-y-0 left-0 h-full w-full max-w-md data-[state=open]:animate-slide-in-from-left',
  right: 'inset-y-0 right-0 h-full w-full max-w-md data-[state=open]:animate-slide-in-from-right',
  top: 'inset-x-0 top-0 w-full max-h-[80vh] data-[state=open]:animate-slide-in-from-top',
  bottom: 'inset-x-0 bottom-0 w-full max-h-[80vh] data-[state=open]:animate-slide-in-from-bottom',
};

export function Sheet({ open, onClose, children, side = 'right', className }: SheetProps) {
  // Handle escape key
  React.useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    if (open) {
      document.addEventListener('keydown', handleEscape);
      document.body.style.overflow = 'hidden';
    }
    return () => {
      document.removeEventListener('keydown', handleEscape);
      document.body.style.overflow = '';
    };
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50">
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/50 transition-opacity"
        onClick={onClose}
        aria-hidden="true"
      />
      {/* Sheet content */}
      <div
        className={cn(
          'fixed bg-white shadow-xl overflow-auto',
          sideStyles[side],
          className
        )}
        data-state={open ? 'open' : 'closed'}
        role="dialog"
        aria-modal="true"
      >
        {children}
      </div>
    </div>
  );
}

interface SheetHeaderProps {
  children: React.ReactNode;
  onClose?: () => void;
  className?: string;
}

export function SheetHeader({ children, onClose, className }: SheetHeaderProps) {
  return (
    <div className={cn('flex items-center justify-between p-6 border-b border-gray-200', className)}>
      <div className="text-lg font-semibold text-gray-900">{children}</div>
      {onClose && (
        <Button variant="ghost" size="sm" onClick={onClose} className="h-8 w-8 p-0">
          <X className="h-4 w-4" />
          <span className="sr-only">Close</span>
        </Button>
      )}
    </div>
  );
}

interface SheetBodyProps {
  children: React.ReactNode;
  className?: string;
}

export function SheetBody({ children, className }: SheetBodyProps) {
  return <div className={cn('p-6 flex-1 overflow-auto', className)}>{children}</div>;
}

interface SheetFooterProps {
  children: React.ReactNode;
  className?: string;
}

export function SheetFooter({ children, className }: SheetFooterProps) {
  return (
    <div className={cn('flex items-center justify-end gap-3 p-6 border-t border-gray-200', className)}>
      {children}
    </div>
  );
}
