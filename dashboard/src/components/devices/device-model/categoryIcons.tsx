import * as React from 'react';
import { Laptop, Tablet, Monitor } from 'lucide-react';
import type { DeviceCategory } from '@/types/device';

export const categoryIcons: Record<DeviceCategory, React.ReactNode> = {
  laptop: <Laptop className="h-4 w-4" />,
  tablet: <Tablet className="h-4 w-4" />,
  desktop: <Monitor className="h-4 w-4" />,
  chromebook: <Laptop className="h-4 w-4" />,
  other: <Monitor className="h-4 w-4" />,
};
