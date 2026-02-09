import * as React from 'react';
import { motion } from 'motion/react';
import { User, Mail, Calendar, Award, TrendingUp, BookOpen, Flame, Target } from 'lucide-react';

interface ProfileProps {
  onBack: () => void;
}

export function Profile({ onBack }: ProfileProps) {
  const userStats = [
    { label: 'Total Words Learned', value: 142, icon: BookOpen, color: 'text-primary' },
    { label: 'Current Streak', value: '7 days', icon: Flame, color: 'text-orange-500' },
    { label: 'Total Sessions', value: 23, icon: Target, color: 'text-secondary' },
    { label: 'Average Accuracy', value: '87%', icon: TrendingUp, color: 'text-green-600' },
  ];

  const achievements = [
    { title: 'First Session', desc: 'Completed your first learning session', date: 'Dec 20, 2025', earned: true },
    { title: 'Week Warrior', desc: 'Maintained a 7-day streak', date: 'Jan 14, 2026', earned: true },
    { title: 'Century Club', desc: 'Learned 100+ words', date: 'Jan 10, 2026', earned: true },
    { title: 'Month Master', desc: 'Maintain a 30-day streak', date: 'Not yet earned', earned: false },
  ];

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="bg-gradient-to-br from-primary/10 to-accent/10 border-b border-border p-8">
        <div className="max-w-4xl mx-auto">
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
          >
            <div className="flex items-start gap-6">
              <div className="w-24 h-24 rounded-full bg-gradient-to-br from-primary to-accent flex items-center justify-center shadow-lg">
                <User className="w-12 h-12 text-white" />
              </div>
              <div className="flex-1">
                <h1 className="text-3xl text-foreground mb-2">Alex Seeker</h1>
                <div className="flex flex-col gap-2 text-muted-foreground">
                  <div className="flex items-center gap-2">
                    <Mail className="w-4 h-4" />
                    <span>alex.seeker@example.com</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Calendar className="w-4 h-4" />
                    <span>Member since December 2025</span>
                  </div>
                </div>
              </div>
            </div>
          </motion.div>
        </div>
      </div>

      <div className="max-w-4xl mx-auto p-8">
        {/* Stats Grid */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
        >
          <h2 className="text-xl text-foreground mb-4">Your Statistics</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
            {userStats.map((stat, index) => (
              <motion.div
                key={stat.label}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 + index * 0.05 }}
                className="bg-card border border-border rounded-lg p-6 shadow"
              >
                <div className="flex items-center gap-2 mb-3">
                  <stat.icon className={`w-5 h-5 ${stat.color}`} />
                </div>
                <p className="text-2xl font-semibold text-foreground mb-1">{stat.value}</p>
                <p className="text-sm text-muted-foreground">{stat.label}</p>
              </motion.div>
            ))}
          </div>
        </motion.div>

        {/* Achievements */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <h2 className="text-xl text-foreground mb-4 flex items-center gap-2">
            <Award className="w-6 h-6 text-primary" />
            Achievements
          </h2>
          <div className="space-y-3">
            {achievements.map((achievement, index) => (
              <motion.div
                key={achievement.title}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 0.4 + index * 0.05 }}
                className={`
                  p-4 rounded-lg border transition-all duration-200
                  ${achievement.earned
                    ? 'bg-primary/5 border-primary/30'
                    : 'bg-muted/30 border-border opacity-60'
                  }
                `}
              >
                <div className="flex items-start gap-4">
                  <div className={`
                    w-12 h-12 rounded-full flex items-center justify-center
                    ${achievement.earned 
                      ? 'bg-gradient-to-br from-primary to-accent' 
                      : 'bg-muted'
                    }
                  `}>
                    <Award className={`w-6 h-6 ${achievement.earned ? 'text-white' : 'text-muted-foreground'}`} />
                  </div>
                  <div className="flex-1">
                    <h3 className="font-semibold text-foreground">{achievement.title}</h3>
                    <p className="text-sm text-muted-foreground mt-1">{achievement.desc}</p>
                    <p className="text-xs text-muted-foreground mt-2">{achievement.date}</p>
                  </div>
                </div>
              </motion.div>
            ))}
          </div>
        </motion.div>

        {/* Learning Preferences */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5 }}
          className="mt-8 bg-card border border-border rounded-lg p-6 shadow"
        >
          <h2 className="text-xl text-foreground mb-4">Learning Preferences</h2>
          <div className="space-y-4">
            <div className="flex justify-between items-center py-3 border-b border-border">
              <div>
                <p className="font-medium text-foreground">Primary Language</p>
                <p className="text-sm text-muted-foreground">German</p>
              </div>
            </div>
            <div className="flex justify-between items-center py-3 border-b border-border">
              <div>
                <p className="font-medium text-foreground">Current Level</p>
                <p className="text-sm text-muted-foreground">B1 - Intermediate</p>
              </div>
            </div>
            <div className="flex justify-between items-center py-3">
              <div>
                <p className="font-medium text-foreground">Daily Goal</p>
                <p className="text-sm text-muted-foreground">15 words per session</p>
              </div>
            </div>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
