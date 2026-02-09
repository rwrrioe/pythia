import * as React from 'react';
import { motion } from 'motion/react';
import { TrendingUp, Calendar, Target, Award, Flame } from 'lucide-react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, BarChart, Bar } from 'recharts';

interface StatsProps {
  onBack: () => void;
}

export function Stats({ onBack }: StatsProps) {
  // Mock data for charts
  const weeklyProgress = [
    { day: 'Mon', words: 12 },
    { day: 'Tue', words: 15 },
    { day: 'Wed', words: 10 },
    { day: 'Thu', words: 14 },
    { day: 'Fri', words: 13 },
    { day: 'Sat', words: 11 },
    { day: 'Sun', words: 12 },
  ];

  const accuracyTrend = [
    { date: 'Jan 8', accuracy: 82 },
    { date: 'Jan 9', accuracy: 85 },
    { date: 'Jan 10', accuracy: 88 },
    { date: 'Jan 11', accuracy: 78 },
    { date: 'Jan 12', accuracy: 90 },
    { date: 'Jan 13', accuracy: 85 },
    { date: 'Jan 14', accuracy: 92 },
  ];

  return (
    <div className="min-h-screen bg-background p-8">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-8"
        >
          <h1 className="text-3xl text-foreground mb-2">Learning Statistics</h1>
          <p className="text-muted-foreground">
            Track your progress and identify patterns in your learning journey.
          </p>
        </motion.div>

        {/* Key Metrics */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8"
        >
          <div className="bg-card border border-border rounded-lg p-6 shadow">
            <div className="flex items-center gap-3 mb-3">
              <div className="p-2 rounded-lg bg-primary/10">
                <Flame className="w-5 h-5 text-orange-500" />
              </div>
              <div>
                <p className="text-2xl font-semibold text-foreground">7</p>
                <p className="text-sm text-muted-foreground">Day Streak</p>
              </div>
            </div>
          </div>

          <div className="bg-card border border-border rounded-lg p-6 shadow">
            <div className="flex items-center gap-3 mb-3">
              <div className="p-2 rounded-lg bg-secondary/10">
                <Target className="w-5 h-5 text-secondary" />
              </div>
              <div>
                <p className="text-2xl font-semibold text-foreground">87%</p>
                <p className="text-sm text-muted-foreground">Avg Accuracy</p>
              </div>
            </div>
          </div>

          <div className="bg-card border border-border rounded-lg p-6 shadow">
            <div className="flex items-center gap-3 mb-3">
              <div className="p-2 rounded-lg bg-accent/10">
                <Calendar className="w-5 h-5 text-accent" />
              </div>
              <div>
                <p className="text-2xl font-semibold text-foreground">23</p>
                <p className="text-sm text-muted-foreground">Sessions</p>
              </div>
            </div>
          </div>

          <div className="bg-card border border-border rounded-lg p-6 shadow">
            <div className="flex items-center gap-3 mb-3">
              <div className="p-2 rounded-lg bg-primary/10">
                <Award className="w-5 h-5 text-primary" />
              </div>
              <div>
                <p className="text-2xl font-semibold text-foreground">142</p>
                <p className="text-sm text-muted-foreground">Words Learned</p>
              </div>
            </div>
          </div>
        </motion.div>

        {/* Charts */}
        <div className="grid md:grid-cols-2 gap-6 mb-8">
          {/* Weekly Progress */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="bg-card border border-border rounded-lg p-6 shadow"
          >
            <h2 className="text-lg font-semibold text-foreground mb-4">Weekly Progress</h2>
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={weeklyProgress}>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(196, 150, 92, 0.1)" />
                <XAxis dataKey="day" stroke="#6B6960" />
                <YAxis stroke="#6B6960" />
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: '#FFFFFF', 
                    border: '1px solid rgba(196, 150, 92, 0.2)',
                    borderRadius: '8px'
                  }} 
                />
                <Bar dataKey="words" fill="#C4965C" radius={[8, 8, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </motion.div>

          {/* Accuracy Trend */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.3 }}
            className="bg-card border border-border rounded-lg p-6 shadow"
          >
            <h2 className="text-lg font-semibold text-foreground mb-4">Accuracy Trend</h2>
            <ResponsiveContainer width="100%" height={250}>
              <AreaChart data={accuracyTrend}>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(196, 150, 92, 0.1)" />
                <XAxis dataKey="date" stroke="#6B6960" />
                <YAxis stroke="#6B6960" domain={[0, 100]} />
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: '#FFFFFF', 
                    border: '1px solid rgba(196, 150, 92, 0.2)',
                    borderRadius: '8px'
                  }} 
                />
                <Area 
                  type="monotone" 
                  dataKey="accuracy" 
                  stroke="#8B9C7D" 
                  fill="#8B9C7D" 
                  fillOpacity={0.3} 
                />
              </AreaChart>
            </ResponsiveContainer>
          </motion.div>
        </div>

        {/* Study Insights */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          className="bg-card border border-border rounded-lg p-6 shadow"
        >
          <h2 className="text-lg font-semibold text-foreground mb-4 flex items-center gap-2">
            <TrendingUp className="w-5 h-5 text-primary" />
            Study Insights
          </h2>
          <div className="grid md:grid-cols-3 gap-6">
            <div className="border-l-4 border-primary pl-4">
              <p className="text-sm text-muted-foreground mb-1">Best Performance</p>
              <p className="text-xl font-semibold text-foreground">Saturday</p>
              <p className="text-xs text-muted-foreground mt-1">92% average accuracy</p>
            </div>
            <div className="border-l-4 border-secondary pl-4">
              <p className="text-sm text-muted-foreground mb-1">Most Active Time</p>
              <p className="text-xl font-semibold text-foreground">Evening</p>
              <p className="text-xs text-muted-foreground mt-1">6 PM - 9 PM</p>
            </div>
            <div className="border-l-4 border-accent pl-4">
              <p className="text-sm text-muted-foreground mb-1">Words per Day</p>
              <p className="text-xl font-semibold text-foreground">12.4</p>
              <p className="text-xs text-muted-foreground mt-1">Last 7 days average</p>
            </div>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
