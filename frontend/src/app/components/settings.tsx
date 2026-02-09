import * as React from 'react';
import { motion } from 'motion/react';
import { ArrowLeft, User, Bell, Globe, Moon, Sun } from 'lucide-react';
import { GreekPattern } from './greek-pattern';
import { useState } from 'react';

interface SettingsProps {
  onBack: () => void;
}

export function Settings({ onBack }: SettingsProps) {
  const [notifications, setNotifications] = useState(true);
  const [darkMode, setDarkMode] = useState(false);
  const [language, setLanguage] = useState('en');

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted to-background p-4 md:p-8 relative">
      <div className="absolute top-0 left-0 right-0">
        <GreekPattern className="w-full h-6 text-primary opacity-40" />
      </div>

      <div className="max-w-3xl mx-auto relative z-10 pt-8">
        {/* Header */}
        <div className="flex items-center gap-4 mb-8">
          <button onClick={onBack} className="p-2 rounded-lg hover:bg-muted transition-colors">
            <ArrowLeft className="w-6 h-6 text-foreground" />
          </button>
          <div>
            <h1 className="text-3xl md:text-4xl text-primary" style={{ fontFamily: 'var(--font-heading)' }}>
              Settings
            </h1>
            <p className="text-muted-foreground">Customize your Oracle experience</p>
          </div>
        </div>

        <div className="space-y-6">
          {/* Profile Section */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="bg-card border border-border rounded-xl p-6 shadow-xl"
          >
            <div className="flex items-center gap-4 mb-6">
              <div className="p-3 bg-primary/10 rounded-full">
                <User className="w-6 h-6 text-primary" />
              </div>
              <h2 className="text-2xl font-bold text-foreground" style={{ fontFamily: 'var(--font-heading)' }}>
                Profile
              </h2>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-semibold text-foreground mb-2">Username</label>
                <input
                  type="text"
                  defaultValue="Seeker of Knowledge"
                  className="w-full p-3 bg-input-background border border-border rounded-lg 
                             focus:outline-none focus:ring-2 focus:ring-primary/50"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-foreground mb-2">Email</label>
                <input
                  type="email"
                  defaultValue="seeker@pythia.app"
                  className="w-full p-3 bg-input-background border border-border rounded-lg 
                             focus:outline-none focus:ring-2 focus:ring-primary/50"
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-foreground mb-2">Learning Level</label>
                <select
                  className="w-full p-3 bg-input-background border border-border rounded-lg 
                             focus:outline-none focus:ring-2 focus:ring-primary/50"
                >
                  <option value="a2">A2 - Elementary</option>
                  <option value="b1" selected>B1 - Intermediate</option>
                  <option value="b2">B2 - Upper Intermediate</option>
                </select>
              </div>
            </div>
          </motion.div>

          {/* Preferences */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="bg-card border border-border rounded-xl p-6 shadow-xl"
          >
            <div className="flex items-center gap-4 mb-6">
              <div className="p-3 bg-secondary/10 rounded-full">
                <Globe className="w-6 h-6 text-secondary" />
              </div>
              <h2 className="text-2xl font-bold text-foreground" style={{ fontFamily: 'var(--font-heading)' }}>
                Preferences
              </h2>
            </div>
            <div className="space-y-4">
              <div className="flex items-center justify-between p-4 bg-muted/30 rounded-lg">
                <div className="flex items-center gap-3">
                  <Bell className="w-5 h-5 text-foreground" />
                  <div>
                    <p className="font-semibold text-foreground">Notifications</p>
                    <p className="text-sm text-muted-foreground">Reminders for review sessions</p>
                  </div>
                </div>
                <button
                  onClick={() => setNotifications(!notifications)}
                  className={`relative w-14 h-8 rounded-full transition-colors ${
                    notifications ? 'bg-secondary' : 'bg-muted'
                  }`}
                >
                  <div
                    className={`absolute top-1 left-1 w-6 h-6 bg-white rounded-full transition-transform ${
                      notifications ? 'translate-x-6' : 'translate-x-0'
                    }`}
                  />
                </button>
              </div>

              <div className="flex items-center justify-between p-4 bg-muted/30 rounded-lg">
                <div className="flex items-center gap-3">
                  {darkMode ? <Moon className="w-5 h-5 text-foreground" /> : <Sun className="w-5 h-5 text-foreground" />}
                  <div>
                    <p className="font-semibold text-foreground">Dark Mode</p>
                    <p className="text-sm text-muted-foreground">Toggle dark theme</p>
                  </div>
                </div>
                <button
                  onClick={() => setDarkMode(!darkMode)}
                  className={`relative w-14 h-8 rounded-full transition-colors ${
                    darkMode ? 'bg-secondary' : 'bg-muted'
                  }`}
                >
                  <div
                    className={`absolute top-1 left-1 w-6 h-6 bg-white rounded-full transition-transform ${
                      darkMode ? 'translate-x-6' : 'translate-x-0'
                    }`}
                  />
                </button>
              </div>

              <div>
                <label className="block text-sm font-semibold text-foreground mb-2">Target Language</label>
                <select
                  value={language}
                  onChange={(e) => setLanguage(e.target.value)}
                  className="w-full p-3 bg-input-background border border-border rounded-lg 
                             focus:outline-none focus:ring-2 focus:ring-primary/50"
                >
                  <option value="en">English → Spanish</option>
                  <option value="fr">English → French</option>
                  <option value="de">English → German</option>
                  <option value="it">English → Italian</option>
                  <option value="pt">English → Portuguese</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-semibold text-foreground mb-2">Daily Goal</label>
                <select
                  className="w-full p-3 bg-input-background border border-border rounded-lg 
                             focus:outline-none focus:ring-2 focus:ring-primary/50"
                >
                  <option value="5">5 words per day</option>
                  <option value="10" selected>10 words per day</option>
                  <option value="15">15 words per day</option>
                  <option value="20">20 words per day</option>
                </select>
              </div>
            </div>
          </motion.div>

          {/* Study Settings */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="bg-card border border-border rounded-xl p-6 shadow-xl"
          >
            <h2 className="text-2xl font-bold text-foreground mb-4" style={{ fontFamily: 'var(--font-heading)' }}>
              Study Settings
            </h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-semibold text-foreground mb-2">Default Study Mode</label>
                <select
                  className="w-full p-3 bg-input-background border border-border rounded-lg 
                             focus:outline-none focus:ring-2 focus:ring-primary/50"
                >
                  <option value="flashcard" selected>Flashcards</option>
                  <option value="quiz">Multiple Choice Quiz</option>
                  <option value="typing">Typing Practice</option>
                  <option value="mixed">Mixed Mode</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-semibold text-foreground mb-2">Auto-Advance After Answer</label>
                <select
                  className="w-full p-3 bg-input-background border border-border rounded-lg 
                             focus:outline-none focus:ring-2 focus:ring-primary/50"
                >
                  <option value="1">1 second</option>
                  <option value="2" selected>2 seconds</option>
                  <option value="3">3 seconds</option>
                  <option value="manual">Manual</option>
                </select>
              </div>
            </div>
          </motion.div>

          {/* Save Button */}
          <motion.button
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="w-full py-4 bg-gradient-to-r from-primary to-accent text-primary-foreground 
                       rounded-lg font-semibold text-lg shadow-lg hover:shadow-xl 
                       transform hover:scale-105 transition-all duration-300"
          >
            Save Changes
          </motion.button>
        </div>
      </div>
    </div>
  );
}
