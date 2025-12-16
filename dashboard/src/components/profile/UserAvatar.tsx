import { cn } from '@/lib/utils';
import { User } from 'lucide-react';

interface UserAvatarProps {
  src?: string | null;
  alt?: string;
  fallback?: string;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  className?: string;
}

const sizeClasses = {
  sm: 'h-8 w-8 text-xs',
  md: 'h-10 w-10 text-sm',
  lg: 'h-12 w-12 text-base',
  xl: 'h-16 w-16 text-lg',
};

const iconSizes = {
  sm: 'h-4 w-4',
  md: 'h-5 w-5',
  lg: 'h-6 w-6',
  xl: 'h-8 w-8',
};

// Generate initials from a name
function getInitials(name?: string): string {
  if (!name) return '';
  const parts = name.trim().split(/\s+/);
  if (parts.length === 1) {
    return parts[0].substring(0, 2).toUpperCase();
  }
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}

// Generate a consistent color based on a string
function stringToColor(str?: string): string {
  if (!str) return 'bg-gradient-to-br from-gray-400 to-gray-500';

  const colors = [
    'bg-gradient-to-br from-cyan-500 to-teal-600',
    'bg-gradient-to-br from-blue-500 to-indigo-600',
    'bg-gradient-to-br from-purple-500 to-pink-600',
    'bg-gradient-to-br from-emerald-500 to-green-600',
    'bg-gradient-to-br from-amber-500 to-orange-600',
    'bg-gradient-to-br from-rose-500 to-red-600',
    'bg-gradient-to-br from-indigo-500 to-purple-600',
    'bg-gradient-to-br from-teal-500 to-cyan-600',
  ];

  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }

  return colors[Math.abs(hash) % colors.length];
}

export function UserAvatar({ src, alt, fallback, size = 'md', className }: UserAvatarProps) {
  const initials = getInitials(fallback || alt);
  const bgColor = stringToColor(fallback || alt);

  if (src) {
    return (
      <img
        src={src}
        alt={alt || 'User avatar'}
        className={cn(
          'rounded-full object-cover ring-2 ring-white',
          sizeClasses[size],
          className
        )}
        onError={(e) => {
          // Hide the image and show fallback
          e.currentTarget.style.display = 'none';
          e.currentTarget.nextElementSibling?.classList.remove('hidden');
        }}
      />
    );
  }

  if (initials) {
    return (
      <div
        className={cn(
          'flex items-center justify-center rounded-full font-medium text-white ring-2 ring-white',
          sizeClasses[size],
          bgColor,
          className
        )}
        title={alt}
      >
        {initials}
      </div>
    );
  }

  return (
    <div
      className={cn(
        'flex items-center justify-center rounded-full bg-gradient-to-br from-cyan-100 to-teal-100 text-cyan-600 ring-2 ring-white shadow-sm',
        sizeClasses[size],
        className
      )}
      title={alt}
    >
      <User className={iconSizes[size]} />
    </div>
  );
}
