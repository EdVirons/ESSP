import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { AlertCircle, Lock, User, Loader2, Sparkles, GraduationCap, Shield, Wrench, School, Headphones, Package, Users, TrendingUp } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';

interface DemoCredential {
  username: string;
  password: string;
  role: string;
  label: string;
  icon: LucideIcon;
  color: string;
}

const demoCredentials: DemoCredential[] = [
  { username: 'admin', password: 'admin123', role: 'ssp_admin', label: 'Admin', icon: Shield, color: 'text-cyan-600 bg-cyan-100' },
  { username: 'school_contact', password: 'school123', role: 'ssp_school_contact', label: 'School Contact', icon: School, color: 'text-green-600 bg-green-100' },
  { username: 'support_agent', password: 'support123', role: 'ssp_support_agent', label: 'Support Agent', icon: Headphones, color: 'text-blue-600 bg-blue-100' },
  { username: 'lead_tech', password: 'lead123', role: 'ssp_lead_tech', label: 'Lead Tech', icon: Users, color: 'text-purple-600 bg-purple-100' },
  { username: 'field_tech', password: 'tech123', role: 'ssp_field_tech', label: 'Field Tech', icon: Wrench, color: 'text-amber-600 bg-amber-100' },
  { username: 'warehouse', password: 'warehouse123', role: 'ssp_warehouse_manager', label: 'Warehouse', icon: Package, color: 'text-rose-600 bg-rose-100' },
  { username: 'sales_marketing', password: 'sales123', role: 'ssp_sales_marketing', label: 'Sales', icon: TrendingUp, color: 'text-orange-600 bg-orange-100' },
];

interface LocationState {
  from?: {
    pathname: string;
  };
}

export function Login() {
  const navigate = useNavigate();
  const location = useLocation();
  const { login, isLoading, error } = useAuth();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [localError, setLocalError] = useState<string | null>(null);

  const state = location.state as LocationState;
  const from = state?.from?.pathname || '/overview';

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLocalError(null);

    if (!username.trim()) {
      setLocalError('Username is required');
      return;
    }

    if (!password.trim()) {
      setLocalError('Password is required');
      return;
    }

    const success = await login({ username, password });
    if (success) {
      navigate(from, { replace: true });
    }
  };

  const displayError = localError || error;

  return (
    <div className="flex min-h-screen">
      {/* Left side - Branding */}
      <div className="hidden lg:flex lg:w-1/2 header-gradient relative overflow-hidden">
        {/* Decorative elements */}
        <div className="absolute top-0 left-0 w-full h-full opacity-10">
          <div className="absolute top-20 left-20 w-64 h-64 rounded-full bg-cyan-300 blur-3xl" />
          <div className="absolute bottom-20 right-20 w-96 h-96 rounded-full bg-teal-300 blur-3xl" />
          <div className="absolute top-1/2 left-1/2 w-48 h-48 rounded-full bg-white blur-2xl" />
        </div>

        <div className="relative z-10 flex flex-col justify-center px-12 text-white">
          <div className="flex items-center gap-4 mb-8">
            <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-cyan-400 to-cyan-600 shadow-xl">
              <Sparkles className="h-8 w-8" />
            </div>
            <div>
              <h1 className="text-4xl font-bold">ESSP</h1>
              <p className="text-cyan-200 text-lg">Education Sector Support Platform</p>
            </div>
          </div>

          <div className="space-y-6 max-w-md">
            <h2 className="text-2xl font-semibold">
              Empowering Education Through Technology
            </h2>
            <p className="text-cyan-100 text-lg leading-relaxed">
              A comprehensive platform for managing school infrastructure, devices, and maintenance operations across the education sector.
            </p>

            <div className="space-y-4 pt-4">
              <div className="flex items-center gap-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/10 backdrop-blur">
                  <GraduationCap className="h-6 w-6 text-cyan-300" />
                </div>
                <div>
                  <p className="font-semibold">School Management</p>
                  <p className="text-cyan-200 text-sm">Track schools, devices, and programs</p>
                </div>
              </div>

              <div className="flex items-center gap-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-white/10 backdrop-blur">
                  <Shield className="h-6 w-6 text-cyan-300" />
                </div>
                <div>
                  <p className="font-semibold">Secure & Reliable</p>
                  <p className="text-cyan-200 text-sm">Enterprise-grade security</p>
                </div>
              </div>
            </div>
          </div>

          <div className="mt-12 pt-8 border-t border-white/20">
            <p className="text-cyan-200 text-sm">
              Powered by EdVirons Technology
            </p>
          </div>
        </div>
      </div>

      {/* Right side - Login form */}
      <div className="flex w-full lg:w-1/2 items-center justify-center bg-gradient-to-br from-slate-50 to-cyan-50 p-8">
        <div className="w-full max-w-md">
          {/* Mobile logo */}
          <div className="lg:hidden flex items-center justify-center gap-3 mb-8">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-cyan-500 to-teal-600 shadow-lg">
              <Sparkles className="h-6 w-6 text-white" />
            </div>
            <span className="text-2xl font-bold text-gray-900">ESSP</span>
          </div>

          <Card className="shadow-xl border-0 bg-white/80 backdrop-blur">
            <CardHeader className="space-y-1 text-center pb-2">
              <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-gradient-to-br from-cyan-100 to-teal-100">
                <Lock className="h-7 w-7 text-cyan-600" />
              </div>
              <CardTitle className="text-2xl font-bold text-gray-900">Welcome Back</CardTitle>
              <CardDescription className="text-gray-500">
                Sign in to access the ESSP Dashboard
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                {displayError && (
                  <div className="flex items-center gap-2 rounded-lg bg-red-50 border border-red-100 p-3 text-sm text-red-600">
                    <AlertCircle className="h-4 w-4 flex-shrink-0" />
                    <span>{displayError}</span>
                  </div>
                )}

                <div className="space-y-2">
                  <label htmlFor="username" className="text-sm font-medium text-gray-700">
                    Username
                  </label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                    <Input
                      id="username"
                      type="text"
                      placeholder="Enter your username"
                      value={username}
                      onChange={(e) => setUsername(e.target.value)}
                      className="pl-10 border-gray-200 focus:border-cyan-500 focus:ring-cyan-500"
                      autoComplete="username"
                      autoFocus
                      disabled={isLoading}
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <label htmlFor="password" className="text-sm font-medium text-gray-700">
                    Password
                  </label>
                  <div className="relative">
                    <Lock className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                    <Input
                      id="password"
                      type="password"
                      placeholder="Enter your password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      className="pl-10 border-gray-200 focus:border-cyan-500 focus:ring-cyan-500"
                      autoComplete="current-password"
                      disabled={isLoading}
                    />
                  </div>
                </div>

                <Button
                  type="submit"
                  className="w-full bg-gradient-to-r from-cyan-600 to-teal-600 hover:from-cyan-700 hover:to-teal-700 text-white shadow-lg shadow-cyan-500/25 transition-all hover:shadow-xl hover:shadow-cyan-500/30"
                  disabled={isLoading}
                >
                  {isLoading ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Signing in...
                    </>
                  ) : (
                    'Sign In'
                  )}
                </Button>
              </form>

              <div className="mt-6 rounded-lg bg-gradient-to-r from-slate-50 to-gray-50 border border-gray-200 p-4">
                <p className="text-xs text-gray-600 mb-3 font-medium">Demo Accounts (click to fill credentials)</p>
                <div className="grid grid-cols-2 gap-2">
                  {demoCredentials.map((cred) => {
                    const Icon = cred.icon;
                    return (
                      <button
                        key={cred.username}
                        type="button"
                        onClick={() => {
                          setUsername(cred.username);
                          setPassword(cred.password);
                        }}
                        className="flex items-center gap-2 text-left text-sm bg-white rounded-lg px-3 py-2 border border-gray-200 hover:border-gray-300 hover:shadow-sm transition-all"
                      >
                        <div className={`flex h-7 w-7 items-center justify-center rounded-md ${cred.color.split(' ')[1]}`}>
                          <Icon className={`h-3.5 w-3.5 ${cred.color.split(' ')[0]}`} />
                        </div>
                        <div className="min-w-0 flex-1">
                          <div className="font-medium text-gray-700 text-xs truncate">{cred.label}</div>
                          <div className="text-[10px] text-gray-400 truncate">{cred.role.replace('ssp_', '')}</div>
                        </div>
                      </button>
                    );
                  })}
                </div>
              </div>

              <div className="mt-6 text-center">
                <p className="text-xs text-gray-400">
                  ESSP Admin Dashboard v1.0
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
