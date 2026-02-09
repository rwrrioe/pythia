import * as React from 'react';
import { ReactNode, useState } from 'react';
import { motion, AnimatePresence } from 'motion/react';
import { BookOpen, LayoutDashboard, User, TrendingUp, Library, Menu, X } from 'lucide-react';
import { GreekPattern } from './greek-pattern';

interface AppLayoutProps {
  children: ReactNode;
  currentScreen: string;
  onNavigate: (screen: string) => void;
}

export function AppLayout({ children, currentScreen, onNavigate }: AppLayoutProps) {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const navItems = [
    { id: 'dashboard', label: 'Dashboard', icon: LayoutDashboard },
    { id: 'library', label: 'Library', icon: Library },
    { id: 'sessions', label: 'Sessions', icon: BookOpen },
    { id: 'stats', label: 'Stats', icon: TrendingUp },
    { id: 'profile', label: 'Profile', icon: User },
  ];

  const handleNavigation = (screen: string) => {
    onNavigate(screen);
    setIsMobileMenuOpen(false);
  };

  return (
    <div className="flex h-screen bg-background overflow-hidden">
      {/* Desktop Sidebar */}
      <aside className="hidden md:flex w-64 bg-card border-r border-border flex-col">
        {/* Logo/Brand */}
        <div className="p-6 border-b border-border">
          <div className="mb-3">
            <GreekPattern className="w-full h-4 text-primary opacity-40" />
          </div>
          <h1 className="text-xl text-foreground" style={{ fontFamily: 'var(--font-body)' }}>
            Pythia
          </h1>
          <p className="text-sm text-muted-foreground mt-1">Language Oracle</p>
        </div>

        {/* Navigation */}
        <nav className="flex-1 p-4">
          <ul className="space-y-1">
            {navItems.map((item) => {
              const isActive = currentScreen === item.id;
              return (
                <li key={item.id}>
                  <button
                    onClick={() => onNavigate(item.id)}
                    className={`
                      w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all duration-200
                      ${isActive 
                        ? 'bg-primary/10 text-primary' 
                        : 'text-foreground hover:bg-muted'
                      }
                    `}
                  >
                    <item.icon className="w-5 h-5" />
                    <span>{item.label}</span>
                  </button>
                </li>
              );
            })}
          </ul>
        </nav>

        {/* Footer Pattern */}
        <div className="p-4 border-t border-border">
          <div className="transform rotate-180">
            <GreekPattern className="w-full h-4 text-primary opacity-40" />
          </div>
        </div>
      </aside>

      {/* Mobile Header */}
      <div className="md:hidden fixed top-0 left-0 right-0 z-50 bg-card border-b border-border">
        <div className="flex items-center justify-between p-4">
          <div className="flex items-center gap-3">
            <h1 className="text-lg text-foreground" style={{ fontFamily: 'var(--font-body)' }}>
              Pythia
            </h1>
          </div>
          <button
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            className="p-2 rounded-lg hover:bg-muted transition-colors"
          >
            {isMobileMenuOpen ? (
              <X className="w-6 h-6 text-foreground" />
            ) : (
              <Menu className="w-6 h-6 text-foreground" />
            )}
          </button>
        </div>
      </div>

      {/* Mobile Menu */}
      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div
            initial={{ x: '100%' }}
            animate={{ x: 0 }}
            exit={{ x: '100%' }}
            transition={{ type: 'tween', duration: 0.3 }}
            className="md:hidden fixed inset-y-0 right-0 z-40 w-64 bg-card border-l border-border shadow-xl"
          >
            <nav className="p-4 pt-20">
              <ul className="space-y-1">
                {navItems.map((item) => {
                  const isActive = currentScreen === item.id;
                  return (
                    <li key={item.id}>
                      <button
                        onClick={() => handleNavigation(item.id)}
                        className={`
                          w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-all duration-200
                          ${isActive 
                            ? 'bg-primary/10 text-primary' 
                            : 'text-foreground hover:bg-muted'
                          }
                        `}
                      >
                        <item.icon className="w-5 h-5" />
                        <span>{item.label}</span>
                      </button>
                    </li>
                  );
                })}
              </ul>
            </nav>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Overlay for mobile menu */}
      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={() => setIsMobileMenuOpen(false)}
            className="md:hidden fixed inset-0 bg-black/50 z-30"
          />
        )}
      </AnimatePresence>

      {/* Main Content Area */}
      <main className="flex-1 overflow-auto pt-16 md:pt-0">
        {children}
      </main>
    </div>
  );
}