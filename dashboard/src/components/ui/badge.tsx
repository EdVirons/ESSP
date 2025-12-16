import * as React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';

const badgeVariants = cva(
  'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:ring-offset-2',
  {
    variants: {
      variant: {
        default:
          'border-transparent bg-cyan-100 text-cyan-800',
        secondary:
          'border-transparent bg-gray-100 text-gray-700',
        destructive:
          'border-transparent bg-red-100 text-red-700',
        success:
          'border-transparent bg-emerald-100 text-emerald-700',
        warning:
          'border-transparent bg-amber-100 text-amber-700',
        info:
          'border-transparent bg-blue-100 text-blue-700',
        purple:
          'border-transparent bg-purple-100 text-purple-700',
        orange:
          'border-transparent bg-orange-100 text-orange-700',
        outline: 'text-gray-600 border-gray-300 bg-white',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  }
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  );
}

export { Badge, badgeVariants };
